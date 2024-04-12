package auth

import (
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
}

type AuthArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
}

func NewAuth(ua AuthArgs) *Auth {
	return &Auth{db: ua.DB}
}

// GetUniqueRecord returns user_id, auth level, secret, expire for a given unique value i.e. login.
func (a *Auth) GetUniqueRecord(unique string) (types.Uid, auth.Level, []byte, time.Time, error) {
	var expires time.Time

	var record struct {
		Userid  int64
		Authlvl auth.Level
		Secret  []byte
		Expires *time.Time
	}

	ctx, cancel := a.utils.GetContext(time.Duration(a.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := "SELECT user_id, secret, expires, level FROM auth WHERE user_name = $1"
	if err := a.db.QueryRow(ctx, stmt, unique).Scan(
		&record.Userid,
		&record.Secret,
		&record.Expires,
		&record.Authlvl,
	); err != nil {
		if err == pgx.ErrNoRows {
			err = nil // Nothing found - clear the error
		}
		return types.ZeroUid, 0, nil, expires, err
	}

	if record.Expires != nil {
		expires = *record.Expires
	}

	return store.EncodeUid(record.Userid), record.Authlvl, record.Secret, expires, nil
}

// GetRecord returns authentication record given user ID and method.
func (a *Auth) GetRecord(uid types.Uid, scheme string) (string, auth.Level, []byte, time.Time, error) {
	var expires time.Time

	var record struct {
		Uname   string
		Authlvl auth.Level
		Secret  []byte
		Expires *time.Time
	}

	ctx, cancel := a.utils.GetContext(time.Duration(a.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := "SELECT user_name, secret, expires, level FROM auth WHERE user_id = $1 AND scheme = $2"
	if err := a.db.QueryRow(ctx, stmt, store.DecodeUid(uid), scheme).Scan(
		&record.Uname,
		&record.Secret,
		&record.Expires,
		&record.Authlvl,
	); err != nil {
		if err == pgx.ErrNoRows {
			err = types.ErrNotFound // Nothing found - use standard error.
		}
		return "", 0, nil, expires, err
	}

	if record.Expires != nil {
		expires = *record.Expires
	}

	return record.Uname, record.Authlvl, record.Secret, expires, nil
}

// AddRecord creates new authentication record
func (a *Auth) AddRecord(uid types.Uid, scheme, unique string, authLvl auth.Level, secret []byte, expires time.Time) error {

	var exp *time.Time
	if !expires.IsZero() {
		exp = &expires
	}
	ctx, cancel := a.utils.GetContext(time.Duration(a.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := "INSERT INTO auth(user_name, user_id, scheme, level, secret, expires) VALUES($1, $2, $3, $4, $5, $6)"
	if _, err := a.db.Exec(ctx, stmt, unique, store.DecodeUid(uid), scheme, authLvl, secret, exp); err != nil {
		if a.shared.IsDupe(err) {
			return types.ErrDuplicate
		}
		return err
	}
	return nil
}

// DelScheme deletes an existing authentication scheme for the user.
func (a *Auth) DelScheme(user types.Uid, scheme string) error {
	ctx, cancel := a.utils.GetContext(time.Duration(a.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	_, err := a.db.Exec(ctx, "DELETE FROM auth WHERE user_id = $1 AND scheme = $2", store.DecodeUid(user), scheme)
	return err
}

// DelAllRecords deletes all authentication records for the user.
func (a *Auth) DelAllRecords(user types.Uid) (int, error) {
	ctx, cancel := a.utils.GetContext(time.Duration(a.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	res, err := a.db.Exec(ctx, "DELETE FROM auth WHERE user_id = $1", store.DecodeUid(user))
	if err != nil {
		return 0, err
	}

	count := res.RowsAffected()
	return int(count), nil
}

// UpdRecord modifies an authentication record. Only non-default/non-zero values are updated.
func (a *Auth) UpdRecord(uid types.Uid, scheme, unique string, authLvl auth.Level, secret []byte, expires time.Time) error {

	parapg := []string{" level = ? "}
	args := []any{authLvl}

	if unique != "" {
		parapg = append(parapg, " user_name = ? ")
		args = append(args, unique)
	}

	if len(secret) > 0 {
		parapg = append(parapg, " secret = ? ")
		args = append(args, secret)
	}

	if !expires.IsZero() {
		parapg = append(parapg, " expires = ? ")
		args = append(args, expires)
	}

	args = append(args, store.DecodeUid(uid), scheme)

	ctx, cancel := a.utils.GetContext(time.Duration(a.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := "UPDATE auth SET " + strings.Join(parapg, ",") + " WHERE user_id = ? AND scheme = ?"
	sql, args := a.shared.ExpandQuery(stmt, args...)

	resp, err := a.db.Exec(ctx, sql, args...)
	if a.shared.IsDupe(err) {
		return types.ErrDuplicate
	}

	if count := resp.RowsAffected(); count <= 0 {
		return types.ErrNotFound
	}

	return err
}
