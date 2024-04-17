package messages

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Messages struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
	uGen   *types.UidGenerator
}

type MessagesArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
	UGen   *types.UidGenerator
}

func NewMessages(ma MessagesArgs) *Messages {
	return &Messages{db: ma.DB, utils: ma.Utils, cfg: ma.Cfg, shared: ma.Shared, uGen: ma.UGen}
}

func (m *Messages) Save(msg *types.Message) error {
	ctx, cancel := m.utils.GetContext(time.Duration(m.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	// store assignes message ID, but we don't use it. Message IDs are not used anywhere.
	// Using a sequential ID provided by the database.
	var id int
	stmt := `
	  INSERT INTO messages(
			created_at,
			updated_at,
			seq_id,
			topic,
			"from",
			head,
			content
		) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	err := m.db.QueryRow(
		ctx,
		stmt,
		msg.CreatedAt,
		msg.UpdatedAt,
		msg.SeqId,
		msg.Topic,
		m.uGen.DecodeUid(types.ParseUid(msg.From)),
		msg.Head,
		m.utils.ToJSON(msg.Content),
	).Scan(&id)

	if err == nil {
		// Replacing ID given by store by ID given by the DB.
		msg.SetUid(types.Uid(id))
	}

	return err
}

// GetAll() returns messages matching the query
func (m *Messages) GetAll(topic string, forUser types.Uid, opts *types.QueryOpt) ([]types.Message, error) {
	var limit = m.cfg.MaxMessageResults
	var lower = 0
	var upper = 1<<31 - 1

	if opts != nil {
		if opts.Since > 0 {
			lower = opts.Since
		}

		if opts.Before > 0 {
			// MySQL BETWEEN is inclusive-inclusive, Tinode API requires inclusive-exclusive, thus -1
			upper = opts.Before - 1
		}

		if opts.Limit > 0 && opts.Limit < limit {
			limit = opts.Limit
		}
	}

	unum := m.uGen.DecodeUid(forUser)

	ctx, cancel := m.utils.GetContext(time.Duration(m.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := `
	  SELECT 
			m.created_at,
			m.updated_at,
			m.deleted_at,
			m.del_id,
			m.seq_id,
			m.topic,
			m.from,
			m.head,
			m.content
		FROM messages AS m 
		LEFT JOIN delete_logs AS d ON 
			  	 d.topic = m.topic 
		   AND m.seq_id BETWEEN d.low AND d.hi - 1 
		   AND d.deleted_for = $1
		WHERE 
					 m.del_id = 0 
		   AND m.topic = $2 
		   AND m.seq_id BETWEEN $3 AND $4 
		   AND d.deleted_for IS NULL
		ORDER BY m.seq_id DESC LIMIT $5
	`

	rows, err := m.db.Query(ctx, stmt, unum, topic, lower, upper, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgs := make([]types.Message, 0, limit)
	for rows.Next() {
		var msg types.Message
		var from int64
		if err = rows.Scan(
			&msg.CreatedAt,
			&msg.UpdatedAt,
			&msg.DeletedAt,
			&msg.DelId,
			&msg.SeqId,
			&msg.Topic,
			&from,
			&msg.Head,
			&msg.Content,
		); err != nil {
			break
		}

		msg.From = m.uGen.EncodeUid(from).String()
		msgs = append(msgs, msg)
	}

	if err == nil {
		err = rows.Err()
	}

	return msgs, err
}

// Get ranges of deleted messages
func (m *Messages) GetDeleted(topic string, forUser types.Uid, opts *types.QueryOpt) ([]types.DelMessage, error) {
	var limit = m.cfg.MaxResults
	var lower = 0
	var upper = 1<<31 - 1

	if opts != nil {
		if opts.Since > 0 {
			lower = opts.Since
		}

		if opts.Before > 1 {
			// DelRange is inclusive-exclusive, while BETWEEN is inclusive-inclisive.
			upper = opts.Before - 1
		}

		if opts.Limit > 0 && opts.Limit < limit {
			limit = opts.Limit
		}
	}

	// Fetch log of deletions
	ctx, cancel := m.utils.GetContext(time.Duration(m.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := `
	  SELECT topic, deleted_for, del_id, low, hi FROM delete_logs 
		WHERE 
					topic = $1 AND del_id BETWEEN $2 AND $3
		  AND (deleted_for = 0 OR deleted_for = $4) 
		ORDER BY del_id LIMIT $5
	`
	rows, err := m.db.Query(ctx, stmt, topic, lower, upper, m.uGen.DecodeUid(forUser), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dellog struct {
		Topic      string
		Deletedfor int64
		Delid      int
		Low        int
		Hi         int
	}

	var dmsgs []types.DelMessage
	var dmsg types.DelMessage
	for rows.Next() {
		if err = rows.Scan(
			&dellog.Topic,
			&dellog.Deletedfor,
			&dellog.Delid,
			&dellog.Low,
			&dellog.Hi,
		); err != nil {
			dmsgs = nil
			break
		}

		if dellog.Delid != dmsg.DelId {
			if dmsg.DelId > 0 {
				dmsgs = append(dmsgs, dmsg)
			}

			dmsg.DelId = dellog.Delid
			dmsg.Topic = dellog.Topic
			if dellog.Deletedfor > 0 {
				dmsg.DeletedFor = m.uGen.EncodeUid(dellog.Deletedfor).String()
			} else {
				dmsg.DeletedFor = ""
			}

			dmsg.SeqIdRanges = nil
		}

		if dellog.Hi <= dellog.Low+1 {
			dellog.Hi = 0
		}

		dmsg.SeqIdRanges = append(dmsg.SeqIdRanges, types.Range{Low: dellog.Low, Hi: dellog.Hi})
	}

	if err == nil {
		err = rows.Err()
	}

	if err == nil {
		if dmsg.DelId > 0 {
			dmsgs = append(dmsgs, dmsg)
		}
	}

	return dmsgs, err
}

// DeleteList marks messages as deleted.
//
// Soft or Hard is defined by forUser value: forUser.IsZero == true is hard.
func (m *Messages) DeleteList(topic string, toDel *types.DelMessage) (err error) {
	ctx, cancel := m.utils.GetContext(time.Duration(m.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := m.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	if err = m.shared.MessageDeleteList(ctx, tx, topic, toDel); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
