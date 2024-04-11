package users

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	t "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
}

type UsersArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
}

func NewUsers(ua UsersArgs) *Users {
	return &Users{db: ua.DB, utils: ua.Utils, cfg: ua.Cfg, shared: ua.Shared}
}

// UserCreate creates a new user. Returns error and true if error
//
// is due to duplicate user name, false for any other error
func (u *Users) Create(user *t.User) error {
	ctx, cancel := u.utils.GetContext(time.Duration(u.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := u.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	decUid := store.DecodeUid(user.Uid())
	stmt := `
	  INSERT INTO users(id, created_at, updated_at, state, access, public, trusted, tags) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err = tx.Exec(ctx, stmt,
		decUid,
		user.CreatedAt,
		user.UpdatedAt,
		user.State,
		user.Access,
		u.utils.ToJSON(user.Public),
		u.utils.ToJSON(user.Trusted),
		user.Tags,
	)

	if err != nil {
		return err
	}

	// Save user's tags to a separate table to make user findable.
	if err := u.shared.AddTags(ctx, tx, "user_tags", "user_id", decUid, user.Tags, false); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Get fetches a single user by user id. If user is not found it returns (nil, nil)
func (u *Users) Get(uid t.Uid) (*t.User, error) {
	ctx, cancel := u.utils.GetContext(time.Duration(u.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	var user t.User
	var id int64

	stmt := `SELECT * FROM users WHERE id = $1 AND state != $2`
	row, err := u.db.Query(ctx, stmt, store.DecodeUid(uid), t.StateDeleted)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		// Nothing found: user does not exist or marked as soft-deleted
		return nil, nil
	}

	err = row.Scan(&id, &user.CreatedAt, &user.UpdatedAt, &user.State, &user.StateAt, &user.Access, &user.LastSeen, &user.UserAgent, &user.Public, &user.Trusted, &user.Tags)
	if err == nil {
		user.SetUid(uid)
		return &user, nil
	}

	return nil, err
}

// UserGetAll returns user records for a given list of user IDs
func (u *Users) GetAll(ids ...t.Uid) ([]t.User, error) {
	uids := make([]any, len(ids))
	for i, id := range ids {
		uids[i] = store.DecodeUid(id)
	}

	users := []t.User{}
	ctx, cancel := u.utils.GetContext(time.Duration(u.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := "SELECT * FROM users WHERE id = ANY ($1) AND state != $2"
	rows, err := u.db.Query(ctx, stmt, uids, t.StateDeleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user t.User
		var id int64
		if err = rows.Scan(
			&id,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.State,
			&user.StateAt,
			&user.Access,
			&user.LastSeen,
			&user.UserAgent,
			&user.Public,
			&user.Trusted,
			&user.Tags,
		); err != nil {
			users = nil
			break
		}

		if user.State == t.StateDeleted {
			continue
		}

		user.SetUid(store.EncodeUid(id))
		users = append(users, user)
	}
	if err == nil {
		err = rows.Err()
	}

	return users, err

}
