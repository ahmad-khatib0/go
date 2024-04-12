package credentials

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Credentials struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
}

type CredentialsArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
}

func NewCredentials(ua CredentialsArgs) *Credentials {
	return &Credentials{db: ua.DB, utils: ua.Utils, cfg: ua.Cfg, shared: ua.Shared}
}

// Upsert adds or updates a validation record. Returns true if inserted, false if updated.
//
// 1. if credential is validated:
//
// 1.1 Hard-delete unconfirmed equivalent record, if exists.
//
// 1.2 Insert new. Report error if duplicate.
//
// 2. if credential is not validated:
//
// 2.1 Check if validated equivalent exist. If so, report an error.
//
// 2.2 Soft-delete all unvalidated records of the same method.
//
// 2.3 Undelete existing credential. Return if successful.
//
// 2.4 Insert new credential record.
func (c *Credentials) Upsert(cred *types.Credential) (bool, error) {
	var err error

	ctx, cancel := c.utils.GetContext(time.Duration(c.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	now := types.TimeNow()
	userId := c.shared.DecodeUidString(cred.User)

	// Enforce uniqueness:
	// if credential is confirmed,         "method:value" must be unique.
	// if credential is not yet confirmed, "userid:method:value" is unique.
	synth := cred.Method + ":" + cred.Value

	if !cred.Done {
		// Check if this credential is already validated.
		var done bool
		err = tx.QueryRow(ctx, "SELECT done FROM credentials WHERE synthetic = $1", synth).Scan(&done)
		if err == nil {
			// Assign err to ensure closing of a transaction.
			err = types.ErrDuplicate
			return false, err
		}
		if err != pgx.ErrNoRows {
			return false, err
		}

		// We are going to insert new record.
		// NOTE: note how this method allow a user to register with different methods
		synth = cred.User + ":" + synth

		// Adding new unvalidated credential. Deactivate all unvalidated records of this user and method.
		stmt := "UPDATE credentials SET deleted_at = $1 WHERE user_id = $2 AND method = $3 AND done = FALSE"
		_, err = tx.Exec(ctx, stmt, now, userId, cred.Method)
		if err != nil {
			return false, err
		}

		// Assume that the record exists and try to update it: undelete, update timestamp and response value.
		stmt = "UPDATE credentials SET updated_at = $1, deleted_at = NULL, response = $2, done = FALSE WHERE synthetic = $3"
		res, err := tx.Exec(ctx, stmt, cred.UpdatedAt, cred.Response, synth)
		if err != nil {
			return false, err
		}

		// If record was updated, then all is fine.
		if numrows := res.RowsAffected(); numrows > 0 {
			return false, tx.Commit(ctx)
		}
	} else {
		// Hard-deleting unconformed record if it exists.
		_, err = tx.Exec(ctx, "DELETE FROM credentials WHERE synthetic = $1", cred.User+":"+synth)
		if err != nil {
			return false, err
		}
	}

	stmt := `
	  INSERT INTO credentials(created_at, updated_at, method, value, synthetic, user_id, response, done) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_,
		err = tx.Exec(ctx,
		stmt,
		cred.CreatedAt,
		cred.UpdatedAt,
		cred.Method,
		cred.Value,
		synth,
		userId,
		cred.Response,
		cred.Done,
	)

	if err != nil {
		if c.shared.IsDupe(err) {
			return true, types.ErrDuplicate
		}
		return true, err
	}

	return true, tx.Commit(ctx)
}

// Del deletes either credentials of the given user. If method is blank all credentials
//
// are removed. If value is blank all credentials of the given method are removed.
func (c *Credentials) Del(uid types.Uid, method, value string) error {
	ctx, cancel := c.utils.GetContext(time.Duration(c.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = c.shared.CredDel(ctx, tx, uid, method, value)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetActive returns currently active unvalidated credential of the given user and method.
func (c *Credentials) GetActive(uid types.Uid, method string) (*types.Credential, error) {

	ctx, cancel := c.utils.GetContext(time.Duration(c.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	var cred types.Credential

	stmt := `
	  SELECT created_at, updated_at, method, value, response, done, retries 
	  FROM credentials WHERE user_id = $1 AND deleted_at IS NULL AND method = $2 AND done = FALSE
	`
	err := c.db.QueryRow(ctx, stmt, store.DecodeUid(uid), method).Scan(
		&cred.CreatedAt,
		&cred.UpdatedAt,
		&cred.Method,
		&cred.Value,
		&cred.Response,
		&cred.Done,
		&cred.Retries,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	cred.User = uid.String()

	return &cred, nil
}

// Confirm marks given credential method as confirmed.
func (c *Credentials) Confirm(uid types.Uid, method string) error {
	ctx, cancel := c.utils.GetContext(time.Duration(c.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := `
	  UPDATE credentials SET updated_at = $1, done = TRUE, synthetic = CONCAT(method,':',value) 
	  WHERE user_id = $2 AND method = $3 AND deleted_at IS NULL AND done = FALSE	
	`

	res, err := c.db.Exec(ctx, stmt, types.TimeNow(), store.DecodeUid(uid), method)
	if err != nil {
		if c.shared.IsDupe(err) {
			return types.ErrDuplicate
		}
		return err
	}

	if numrows := res.RowsAffected(); numrows < 1 {
		return types.ErrNotFound
	}

	return nil
}

// Fail increments failure count of the given validation method.
func (c *Credentials) Fail(uid types.Uid, method string) error {
	ctx, cancel := c.utils.GetContext(time.Duration(c.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := "UPDATE credentials SET updated_at = $1, retries = retries+1 WHERE user_id = $2 AND method = $3 AND done = FALSE"
	_, err := c.db.Exec(ctx, stmt, types.TimeNow(), store.DecodeUid(uid), method)
	return err
}
