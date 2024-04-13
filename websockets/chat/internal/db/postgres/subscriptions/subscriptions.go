package subscriptions

import (
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscriptions struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
}

type SubscriptionsArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
}

func NewSubscriptions(ua SubscriptionsArgs) *Subscriptions {
	return &Subscriptions{db: ua.DB, utils: ua.Utils, cfg: ua.Cfg, shared: ua.Shared}
}

// Get a subscription of a user to a topic.
func (s *Subscriptions) Get(topic string, user types.Uid, keepDeleted bool) (*types.Subscription, error) {
	ctx, cancel := s.utils.GetContext(time.Duration(s.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	var sub types.Subscription
	var userId int64
	var modeWant, modeGiven []byte

	stmt := `
		SELECT 
			created_at,
			updated_at,
			deleted_at,
			user_id AS user,
			topic,
			del_id,
			received_seq_id,
			read_seq_id,
			mode_want,
			mode_given,
			private 
		FROM subscriptions WHERE topic = $1 AND user_id = $2
	`
	err := s.db.QueryRow(ctx, stmt, topic, store.DecodeUid(user)).Scan(
		&sub.CreatedAt,
		&sub.UpdatedAt,
		&sub.DeletedAt,
		&userId,
		&sub.Topic,
		&sub.DelId,
		&sub.RecvSeqId,
		&sub.ReadSeqId,
		&modeWant,
		&modeGiven,
		&sub.Private,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Nothing found - clear the error
			err = nil
		}
		return nil, err
	}

	if !keepDeleted && sub.DeletedAt != nil {
		return nil, nil
	}

	sub.User = store.EncodeUid(userId).String()
	sub.ModeWant.Scan(modeWant)
	sub.ModeGiven.Scan(modeGiven)

	return &sub, nil
}

// SubsForUser loads all user's subscriptions. Does NOT load Public or Private
//
// values and does not load deleted subscriptions.
func (s *Subscriptions) SubsForUser(forUser types.Uid) ([]types.Subscription, error) {
	q := `
		SELECT 
			created_at,
			updated_at,
			deleted_at,
			user_id AS user,
			topic,
			del_id,
			received_seq_id,
			read_seq_id,
			mode_want,
			mode_given 
		FROM subscriptions 
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	args := []any{store.DecodeUid(forUser)}
	ctx, cancel := s.utils.GetContext(time.Duration(s.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}
	rows, err := s.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []types.Subscription
	var sub types.Subscription
	var userId int64
	var modeWant, modeGiven []byte

	for rows.Next() {
		if err = rows.Scan(
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&sub.DeletedAt,
			&userId,
			&sub.Topic,
			&sub.DelId,
			&sub.RecvSeqId,
			&sub.ReadSeqId,
			&modeWant,
			&modeGiven,
		); err != nil {
			break
		}

		sub.User = store.EncodeUid(userId).String()
		sub.ModeWant.Scan(modeWant)
		sub.ModeGiven.Scan(modeGiven)
		subs = append(subs, sub)
	}

	if err == nil {
		err = rows.Err()
	}

	return subs, err
}

// SubsForTopic fetches all subsciptions for a topic. Does NOT load Public value.
//
// # The difference between UsersForTopic vs SubsForTopic is that the former
//
// loads user.public+trusted, the latter does not.
func (s *Subscriptions) SubsForTopic(topic string, keepDeleted bool, opts *types.QueryOpt) ([]types.Subscription, error) {
	q := `
		SELECT 
			created_at,
			updated_at,
			deleted_at,
			user_id AS user,
			topic,
			del_id,
			received_seq_id,
			read_seq_id,
			mode_want,
			mode_given,
			private 
		FROM subscriptions WHERE topic = ?
	`

	args := []any{topic}

	if !keepDeleted {
		// Filter out deleted rows.
		q += " AND deleted_at IS NULL "
	}

	limit := s.cfg.MaxResults
	if opts != nil {
		// Ignore IfModifiedSince - we must return all entries
		// Those unmodified will be stripped of Public & Private.

		if !opts.User.IsZero() {
			q += " AND userid=?"
			args = append(args, store.DecodeUid(opts.User))
		}

		if opts.Limit > 0 && opts.Limit < limit {
			limit = opts.Limit
		}
	}

	q += " LIMIT ? "
	args = append(args, limit)
	q, args = s.shared.ExpandQuery(q, args...)

	ctx, cancel := s.utils.GetContext(time.Duration(s.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	rows, err := s.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []types.Subscription
	var sub types.Subscription
	var userId int64
	var modeWant, modeGiven []byte

	for rows.Next() {
		if err = rows.Scan(
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&sub.DeletedAt,
			&userId,
			&sub.Topic,
			&sub.DelId,
			&sub.RecvSeqId,
			&sub.ReadSeqId,
			&modeWant,
			&modeGiven,
			&sub.Private,
		); err != nil {
			break
		}

		sub.User = store.EncodeUid(userId).String()
		sub.ModeWant.Scan(modeWant)
		sub.ModeGiven.Scan(modeGiven)
		subs = append(subs, sub)
	}

	if err == nil {
		err = rows.Err()
	}

	return subs, err
}

// Update updates one or multiple subscriptions to a topic.
func (s *Subscriptions) Update(topic string, user types.Uid, update map[string]any) error {
	ctx, cancel := s.utils.GetContext(time.Duration(s.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	cols, args := s.shared.UpdateByMap(update)
	args = append(args, topic)
	q := "UPDATE subscriptions SET " + strings.Join(cols, ",") + " WHERE topic = ?"

	if !user.IsZero() {
		// Update just one topic subscription
		args = append(args, store.DecodeUid(user))
		q += " AND user_id = ? "
	}

	q, args = s.shared.ExpandQuery(q, args...)
	if _, err = tx.Exec(ctx, q, args...); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Delete marks subscription as deleted.
func (s *Subscriptions) Delete(topic string, user types.Uid) error {
	ctx, cancel := s.utils.GetContext(time.Duration(s.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	decUID := store.DecodeUid(user)
	now := types.TimeNow()

	stmt := `
	  UPDATE subscriptions SET 
				updated_at = $1,
		    deleted_at = $2 
		WHERE 
				topic = $3 
		AND user_id = $4 
		AND deleted_at IS NULL
	`
	res, err := tx.Exec(ctx, stmt, now, now, topic, decUID)
	if err != nil {
		return err
	}

	affected := res.RowsAffected()
	if affected == 0 {
		// ensure tx.Rollback() above is ran
		err = types.ErrNotFound
		return err
	}

	// Remove records of messages soft-deleted by this user.
	_, err = tx.Exec(ctx, "DELETE FROM dellog WHERE topic = $1 AND deleted_for = $2", topic, decUID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
