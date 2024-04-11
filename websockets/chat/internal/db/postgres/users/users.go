package users

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
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

// UserDelete deletes specified user: wipes completely (hard-delete) or marks as deleted.
// TODO: report when the user is not found.
func (u *Users) Delete(uid t.Uid, hard bool) error {
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

	now := t.TimeNow()
	decUid := store.DecodeUid(uid)

	if hard {
		// Delete user's devices,  if t.ErrNotFound = user  that means the user no devices.
		if err = u.shared.DeviceDelete(ctx, tx, uid, ""); err != nil && err != t.ErrNotFound {
			return err
		}

		// Delete user's subscriptions in all topics.
		if err = u.shared.SubDelForUser(ctx, tx, uid, true); err != nil {
			return err
		}

		// Delete records of messages soft-deleted for the user.
		if _, err = tx.Exec(ctx, "DELETE FROM dellog WHERE deleted_for = $1", decUid); err != nil {
			return err
		}

		// Can't delete user's messages in all topics because we cannot notify topics of such deletion.
		// Just leave the messages there marked as sent by "not found" user.

		// Delete topics where the user is the owner:

		// First delete all messages in those topics.
		stmt := "DELETE FROM delete_logs USING topics WHERE topics.name = delete_logs.topic AND topics.owner = $1"
		_, err = tx.Exec(ctx, stmt, decUid)
		if err != nil {
			return err
		}
		stmt = "DELETE FROM messages USING topics WHERE topics.name = messages.topic AND topics.owner = $1"
		if _, err = tx.Exec(ctx, stmt, decUid); err != nil {
			return err
		}

		// Delete all subscriptions:
		stmt = "DELETE FROM subscriptions USING topics WHERE topics.name = subscriptions.topic AND topics.owner = $1"
		if _, err = tx.Exec(ctx, stmt, decUid); err != nil {
			return err
		}

		// Delete topic tags.
		stmt = "DELETE FROM topic_tags USING topics WHERE topics.name = topic_tags.topic AND topics.owner = $1"
		if _, err = tx.Exec(ctx, stmt, decUid); err != nil {
			return err
		}

		// And finally delete the topics.
		stmt = "DELETE FROM topics WHERE owner = $1"
		if _, err = tx.Exec(ctx, stmt, decUid); err != nil {
			return err
		}

		// Delete user's authentication records.
		if _, err = tx.Exec(ctx, "DELETE FROM auth WHERE user_id = $1", decUid); err != nil {
			return err
		}

		// Delete all credentials.
		if err = u.shared.CredDel(ctx, tx, uid, "", ""); err != nil && err != types.ErrNotFound {
			return err
		}

		if _, err = tx.Exec(ctx, "DELETE FROM user_tags WHERE user_id = $1 ", decUid); err != nil {
			return err
		}

		if _, err = tx.Exec(ctx, "DELETE FROM users WHERE id = $1 ", decUid); err != nil {
			return err
		}
	} else {

		// Disable all user's subscriptions. That includes p2p subscriptions. No need to delete them.
		if err = u.shared.SubDelForUser(ctx, tx, uid, true); err != nil {
			return err
		}

		// Disable all subscriptions to topics where the user is the owner.
		stmt := `
		  UPDATE subscriptions SET updated_at = $1 , deleted_at = $2 
      FROM topics WHERE subscriptions.topic = topics.name AND topics.owner = $3;
		`
		if _, err = tx.Exec(ctx, stmt, now, now, decUid); err != nil {
			return err
		}

		// Disable p2p topics with the user (p2p topic's owner is 0).
		stmt = "UPDATE topics SET updated_at = $1, touched_at = $2, state = $3, state_at = $4 WHERE owner = $5"
		if _, err = tx.Exec(ctx, stmt, now, now, types.StateOK, now, decUid); err != nil {
			return err
		}

		// Disable p2p topics with the user (p2p topic's owner is 0).
		stmt = `
		  UPDATE topics SET updated_at = $1, touched_at = $2, state_at = $3, state = $4
			FROM subscriptions 
			WHERE 
						topics.name = subscriptions.topic
        AND topics.owner = 0 AND subscriptions.user_id = $5
		`
		if _, err = tx.Exec(ctx, stmt, now, now, now, types.StateDeleted, decUid); err != nil {
			return err
		}

		// Disable the other user's subscription to a disabled p2p topic.
		stmt = `
			UPDATE subscriptions AS s_one SET updated_at = $1, deleted_at = $2
			FROM subscriptions AS s_two WHERE s_one.topic = s_two.topic
			AND s_two.user_id = $3 AND s_two.topic LIKE 'p2p%'
		`
		if _, err = tx.Exec(ctx, stmt, now, now, decUid); err != nil {
			return err
		}

		// Disable user.
		stmt = "UPDATE users SET updated_at = $1, state = $2, state_at = $3 WHERE id = $4"
		if _, err = tx.Exec(ctx, stmt, now, types.StateDeleted, decUid); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
