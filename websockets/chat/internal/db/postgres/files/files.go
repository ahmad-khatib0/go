package files

import (
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Files struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
	uGen   *types.UidGenerator
}

type FilesArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
	UGen   *types.UidGenerator
}

func NewFiles(fa FilesArgs) *Files {
	return &Files{db: fa.DB, utils: fa.Utils, cfg: fa.Cfg, shared: fa.Shared, uGen: fa.UGen}
}

// StartUpload initializes a file upload
func (f *Files) StartUpload(fd *types.FileDef) error {
	ctx, cancel := f.utils.GetContext(time.Duration(f.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	var user any
	if fd.User != "" {
		user = f.uGen.DecodeUid(types.ParseUid(fd.User))
	}

	stmt := `
	  INSERT INTO file_uploads(
				id, created_at, updated_at, user_id, status, mime_type, size, location
		) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := f.db.Exec(
		ctx,
		stmt,
		f.uGen.DecodeUid(fd.Uid()),
		fd.CreatedAt,
		fd.UpdatedAt,
		user,
		fd.Status,
		fd.MimeType,
		fd.Size,
		fd.Location,
	)

	return err
}

// FinishUpload marks file upload as completed, successfully or otherwise
func (f *Files) FinishUpload(fd *types.FileDef, success bool, size int64) (*types.FileDef, error) {
	ctx, cancel := f.utils.GetContext(time.Duration(f.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := f.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	now := types.TimeNow()
	if success {
		stmt := ` UPDATE file_uploads SET updated_at = $1, status = $2, size = $3 WHERE id = $4 `
		_, err = tx.Exec(ctx, stmt, now, types.UploadCompleted, size, f.uGen.DecodeUid(fd.Uid()))
		if err != nil {
			return nil, err
		}

		fd.Status = types.UploadCompleted
		fd.Size = size
	} else {
		// Deleting the record: there is no value in keeping it in the DB.
		_, err = tx.Exec(ctx, "DELETE FROM file_uploads WHERE id = $1", f.uGen.DecodeUid(fd.Uid()))
		if err != nil {
			return nil, err
		}

		fd.Status = types.UploadFailed
		fd.Size = 0
	}

	fd.UpdatedAt = now
	return fd, tx.Commit(ctx)
}

// Get fetches a record of a specific file
func (f *Files) Get(fid string) (*types.FileDef, error) {
	id := types.ParseUid(fid)
	if id.IsZero() {
		return nil, types.ErrMalformed
	}

	ctx, cancel := f.utils.GetContext(time.Duration(f.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	var fd types.FileDef
	var ID int64
	var userId int64

	stmt := `
	  SELECT 
				id,
				created_at,
				updated_at,
				user_id AS user,
				status,
				mime_type,
				size,
				location 
	  FROM file_uploads WHERE id = $1
	`

	err := f.db.QueryRow(ctx,
		stmt,
		f.uGen.DecodeUid(id)).Scan(&ID,
		&fd.CreatedAt,
		&fd.UpdatedAt,
		&userId,
		&fd.Status,
		&fd.MimeType,
		&fd.Size,
		&fd.Location,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	fd.SetUid(f.uGen.EncodeUid(ID))
	fd.User = f.uGen.EncodeUid(userId).String()
	return &fd, nil
}

// DeleteUnused deletes file upload records.
func (f *Files) DeleteUnused(olderThan time.Time, limit int) ([]string, error) {
	ctx, cancel := f.utils.GetContext(time.Duration(f.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := f.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Garbage collecting entries which as either marked as deleted,
	// or lack message references, or have no user assigned.
	query := `
	  SELECT fu.id, fu.location FROM file_uploads AS fu 
		LEFT JOIN file_message_links AS fml ON fml.file_id = fu.id 
		WHERE fml.id IS NULL
	`

	var args []any
	if !olderThan.IsZero() {
		query += " AND fu.updated_at < ? "
		args = append(args, olderThan)
	}

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	query, _ = f.shared.ExpandQuery(query, args...)
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var locations []string
	var ids []any

	for rows.Next() {
		var id int
		var loc string
		if err = rows.Scan(&id, &loc); err != nil {
			break
		}

		if loc != "" {
			locations = append(locations, loc)
		}

		ids = append(ids, id)
	}

	if err == nil {
		err = rows.Err()
	}
	rows.Close()

	if err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		query, ids = f.shared.ExpandQuery("DELETE FROM file_uploads WHERE id IN (?)", ids)
		_, err = tx.Exec(ctx, query, ids...)
		if err != nil {
			return nil, err
		}
	}

	return locations, tx.Commit(ctx)
}

// LinkAttachments connects given topic or message to the file record IDs from the list.
func (f *Files) LinkAttachments(topic string, userId, msgId types.Uid, fids []string) error {
	if len(fids) == 0 || (topic == "" && msgId.IsZero() && userId.IsZero()) {
		return types.ErrMalformed
	}

	now := types.TimeNow()
	var args []any
	var linkId any
	var linkBy string

	if !msgId.IsZero() {
		linkBy = "message_id"
		linkId = int64(msgId)

	} else if topic != "" {
		linkBy = "topic"
		linkId = topic
		// Only one attachment per topic is permitted at this time.
		fids = fids[0:1]

	} else {
		linkBy = "user_id"
		linkId = f.uGen.DecodeUid(userId)
		// Only one attachment per user is permitted at this time.
		fids = fids[0:1]
	}

	// Decoded ids
	var dids []any
	for _, fid := range fids {
		id := types.ParseUid(fid)
		if id.IsZero() {
			return types.ErrMalformed
		}

		dids = append(dids, f.uGen.DecodeUid(id))
	}

	for _, id := range dids {
		// createdat,fileid,[msgid|topic|userid]
		args = append(args, now, id, linkId)
	}

	ctx, cancel := f.utils.GetContext(time.Duration(f.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := f.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Unlink earlier uploads on the same topic or user allowing them to be garbage-collected.
	if msgId.IsZero() {
		sql := "DELETE FROM file_message_links WHERE " + linkBy + " = $1 "
		_, err = tx.Exec(ctx, sql, linkId)
		if err != nil {
			return err
		}
	}

	stmt := "INSERT INTO file_message_links(created_at, file_id, " + linkBy + ") VALUES (?,?,?)"

	query, args := f.shared.ExpandQuery(stmt+strings.Repeat(",(?,?,?)", len(dids)-1), args...)
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
