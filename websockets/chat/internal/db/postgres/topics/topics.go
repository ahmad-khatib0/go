package topics

import (
	"context"
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/common"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Topics struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
	logger *logger.Logger
}

type TopicsArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
	Logger *logger.Logger
}

func NewTopics(ua TopicsArgs) *Topics {
	return &Topics{db: ua.DB, utils: ua.Utils, cfg: ua.Cfg, shared: ua.Shared}
}

func (t *Topics) Create(topic *types.Topic) error {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = t.createTopic(ctx, tx, topic)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// CreateP2P given two users creates a p2p topic
func (t *Topics) CreateP2P(initiator, invited *types.Subscription) error {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}
	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = t.createSubscription(ctx, tx, initiator, false)
	if err != nil {
		return err
	}

	err = t.createSubscription(ctx, tx, invited, true)
	if err != nil {
		return err
	}

	topic := &types.Topic{ObjHeader: types.ObjHeader{ID: initiator.Topic}}
	topic.ObjHeader.MergeTimes(&initiator.ObjHeader)
	topic.TouchedAt = initiator.GetTouchedAt()
	err = t.createTopic(ctx, tx, topic)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get loads a single topic by name, if it exists. If the topic does not exist the call returns (nil, nil)
func (t *Topics) Get(topic string) (*types.Topic, error) {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	// Fetch topic by name
	var tt = new(types.Topic)
	var owner int64

	stmt := `
	  SELECT 
			created_at,
			updated_at,
			state,
			state_at,
			touched_at,
			name AS id,
			use_bt,
			access,
			owner,
			seq_id,
			del_id,
			public,
			trusted,
			tags
	  FROM topics WHERE name = $1
	`

	err := t.db.QueryRow(ctx, stmt, topic).Scan(
		&tt.CreatedAt,
		&tt.UpdatedAt,
		&tt.State,
		&tt.StateAt,
		&tt.TouchedAt,
		&tt.ID,
		&tt.UseBt,
		&tt.Access,
		&owner,
		&tt.SeqId,
		&tt.DelId,
		&tt.Public,
		&tt.Trusted,
		&tt.Tags,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil // Nothing found - clear the error
		}
		return nil, err
	}

	tt.Owner = store.EncodeUid(owner).String()
	return tt, nil
}

// TopicsForUser loads user's contact list:
//
//	p2p and grp topics, except for 'me' & 'fnd' subscriptions.
//
// Reads and denormalizes Public value.
func (t *Topics) TopicsForUser(uid types.Uid, keepDeleted bool, opts *types.QueryOpt) ([]types.Subscription, error) {
	// Fetch ALL user's subscriptions, even those which has not been modified recently.
	//
	// We are going to use these subscriptions to fetch topics and users which may have been modified recently.

	q := `
		SELECT 
			created_at,
			updated_at,
			deleted_at,
			topic,
			del_id,
			received_seq_id,
			read_seq_id,
			mode_want,
			mode_given,
			private 
		FROM subscriptions WHERE user_id=? 
	`

	args := []any{store.DecodeUid(uid)}
	if !keepDeleted {
		q += " AND deleted_at IS NULL " // Filter out deleted rows.
	}

	limit := 0
	ipg := time.Time{}

	if opts != nil {
		if opts.Topic != "" {
			q += " AND topic=? "
			args = append(args, opts.Topic)
		}

		// Apply the limit only when the client does not manage the cache (or cold start).
		//
		// Otherwise have to get all subscriptions and do a manual join with users/topics.
		if opts.IfModifiedSince == nil {
			if opts.Limit > 0 && opts.Limit < t.cfg.MaxResults {
				limit = opts.Limit
			} else {
				limit = t.cfg.MaxResults
			}
		} else {
			ipg = *opts.IfModifiedSince
		}
	} else {
		limit = t.cfg.MaxResults
	}

	if limit > 0 {
		q += " LIMIT ?"
		args = append(args, limit)
	}

	q, args = t.shared.ExpandQuery(q, args...)

	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	rows, err := t.db.Query(ctx, q, args...)
	if err != nil {
		rows.Close()
		return nil, err
	}

	// Fetch subscriptions. Two queries are needed: users table (p2p) and topics table (grp).
	// Prepare a list of separate subscriptions to users vs topics
	join := make(map[string]types.Subscription) // Keeping these to make a join with table for .private and .access
	topq := make([]any, 0, 16)
	usrq := make([]any, 0, 16)

	for rows.Next() {
		var sub types.Subscription
		var modeWant, modeGiven []byte

		if err = rows.Scan(
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&sub.DeletedAt,
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

		sub.ModeWant.Scan(modeWant)
		sub.ModeGiven.Scan(modeGiven)
		sub.User = uid.String()
		tname := sub.Topic
		tcat := types.GetTopicCat(tname) // i.e p2p, grp etc

		if tcat == types.TopicCatMe || tcat == types.TopicCatFnd {
			// One of 'me', 'fnd' subscriptions, skip. Don't skip 'sys' subscription.
			continue

		} else if tcat == types.TopicCatP2P {
			// P2P subscription, find the other user to get user.Public and user.Trusted.
			uid1, uid2, _ := types.TopicsParseP2P(tname) // user 1, and the second user ids

			if uid1 == uid {
				usrq = append(usrq, store.DecodeUid(uid2))
				sub.SetWith(uid2.UserId())
			} else {
				usrq = append(usrq, store.DecodeUid(uid1))
				sub.SetWith(uid1.UserId())
			}

			topq = append(topq, tname) // append the topic name
		} else {
			// Group or 'sys' subscription.
			if tcat == types.TopicCatGrp {
				// Maybe convert channel name to topic name (it won't be changed if its already a group).
				tname = types.TopicsChnToGrp(tname)
			}
			topq = append(topq, tname)
		}

		sub.Private = t.utils.FromJSON(sub.Private)
		join[tname] = sub
	}

	if err == nil {
		err = rows.Err()
	}
	rows.Close()

	if err != nil {
		return nil, err
	}

	var subs []types.Subscription
	if len(join) == 0 {
		return subs, nil
	}

	// Fetch grp topics and join to subscriptions.
	if len(topq) > 0 {
		q = `
		   SELECT 
					created_at,
					updated_at,
					state,
					state_at,
					touched_at,
					name AS id,
					use_bt,
					access,
					seq_id,
					del_id,
					public,
					trusted,
					tags
		   FROM topics WHERE name IN (?)
		`

		newArgs := []any{topq}
		if !keepDeleted {
			// Optionally skip deleted topics.
			q += " AND state != ? "
			newArgs = append(newArgs, types.StateDeleted)
		}

		if !ipg.IsZero() {
			// Use cache timestamp if provided: get newer entries only.
			q += " AND touched_at > ? "
			newArgs = append(newArgs, ipg)

			if limit > 0 && limit < len(topq) {
				// No point in fetching more than the requested limit.
				q += " ORDER BY touched_at LIMIT ? "
				newArgs = append(newArgs, limit)
			}
		}

		q, newArgs = t.shared.ExpandQuery(q, newArgs...)
		ctx2, cancel2 := t.utils.GetContext(time.Duration(t.cfg.SqlTimeout))
		if cancel2 != nil {
			defer cancel2()
		}

		rows, err = t.db.Query(ctx2, q, newArgs...)
		if err != nil {
			rows.Close()
			return nil, err
		}

		var top types.Topic
		for rows.Next() {
			if err = rows.Scan(
				&top.CreatedAt,
				&top.UpdatedAt,
				&top.State,
				&top.StateAt,
				&top.TouchedAt,
				&top.ID,
				&top.UseBt,
				&top.Access,
				&top.SeqId,
				&top.DelId,
				&top.Public,
				&top.Trusted,
				&top.Tags,
			); err != nil {
				break
			}

			sub := join[top.ID]
			// Check if sub.UpdatedAt needs to be adjusted to earlier or later time.
			sub.UpdatedAt = common.SelectLatestTime(sub.UpdatedAt, top.UpdatedAt)
			sub.SetState(top.State)
			sub.SetTouchedAt(top.TouchedAt)
			sub.SetSeqId(top.SeqId)

			if types.GetTopicCat(sub.Topic) == types.TopicCatGrp {
				sub.SetPublic(top.Public)
				sub.SetTrusted(top.Trusted)
			}

			// Put back the updated value of a subsription, will process further below
			join[top.ID] = sub
		}

		if err == nil {
			err = rows.Err()
		}
		rows.Close()

		if err != nil {
			return nil, err
		}
	}

	// Fetch p2p users and join to p2p subscriptions.
	if len(usrq) > 0 {
		q = `
		  SELECT 
			  id,
				created_at,
				updated_at,
				state,
				state_at,
				access,
				last_seen,
				user_agent,
				public,
				trusted,
				tags 
		  FROM users WHERE id IN (?)
		`

		newArgs := []any{usrq}
		if !keepDeleted {
			// Optionally skip deleted users.
			q += " AND state != ? "
			newArgs = append(newArgs, types.StateDeleted)
		}

		// Ignoring ipg: we need all users to get LastSeen and UserAgent.
		q, newArgs = t.shared.ExpandQuery(q, newArgs...)

		ctx3, cancel3 := t.utils.GetContext(time.Duration(t.cfg.SqlTimeout))
		if cancel3 != nil {
			defer cancel3()
		}

		rows, err = t.db.Query(ctx3, q, newArgs...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var usr2 types.User
			var id int64
			if err = rows.Scan(
				&id,
				&usr2.CreatedAt,
				&usr2.UpdatedAt,
				&usr2.State,
				&usr2.StateAt,
				&usr2.Access,
				&usr2.LastSeen,
				&usr2.UserAgent,
				&usr2.Public,
				&usr2.Trusted,
				&usr2.Tags,
			); err != nil {
				break
			}

			usr2.ID = store.EncodeUid(id).String()
			joinOn := uid.P2PName(types.ParseUid(usr2.ID))

			if sub, ok := join[joinOn]; ok {
				sub.UpdatedAt = common.SelectLatestTime(sub.UpdatedAt, usr2.UpdatedAt)
				sub.SetState(usr2.State)
				sub.SetPublic(usr2.Public)
				sub.SetTrusted(usr2.Trusted)
				sub.SetDefaultAccess(usr2.Access.Auth, usr2.Access.Anon)
				sub.SetLastSeenAndUA(usr2.LastSeen, usr2.UserAgent)
				join[joinOn] = sub
			}
		}

		if err == nil {
			err = rows.Err()
		}

		if err != nil {
			return nil, err
		}
	}

	subs = make([]types.Subscription, 0, len(join))
	for _, sub := range join {
		subs = append(subs, sub)
	}

	return common.SelectEarliestUpdatedSubs(subs, opts, t.cfg.MaxResults), nil
}

// UsersForTopic loads users subscribed to the given topic.
//
// The difference between UsersForTopic vs SubsForTopic is that the former loads user.Public,
//
// the latter does not.
func (t *Topics) UsersForTopic(topic string, keepDeleted bool, opts *types.QueryOpt) ([]types.Subscription, error) {
	tcat := types.GetTopicCat(topic)

	// Fetch all subscribed users. The number of users is not large
	q := `
		SELECT 
			s.created_at,
			s.updated_at,
			s.deleted_at,
			s.user_id,
			s.topic,
			s.del_id,
			s.received_seq_id,
			s.read_seq_id,
			s.mode_want,
			s.mode_given,
			u.public,
			u.trusted,
			u.last_seen,
			u.user_agent,
			s.private
		FROM subscriptions AS s 
		JOIN users AS u ON s.user_id = u.id
		WHERE s.topic = ?
	`

	args := []any{topic}
	if !keepDeleted {
		// Filter out rows with users deleted
		q += " AND u.state != ? "
		args = append(args, types.StateDeleted)

		// For p2p topics we must load all subscriptions including deleted.
		// Otherwise it will be impossible to swipe Public values.
		if tcat != types.TopicCatP2P {
			// Filter out deleted subscriptions.
			q += " AND s.deleted_at IS NULL"
		}
	}

	limit := t.cfg.MaxResults
	var oneUser types.Uid

	if opts != nil {
		// Ignore IfModifiedSince:
		// loading all entries because a topic cannot have too many subscribers.
		// Those unmodified will be stripped of Public & Private.

		if !opts.User.IsZero() {
			// For p2p topics we have to fetch both users otherwise public cannot be swapped.
			if tcat != types.TopicCatP2P {
				q += " AND s.user_id = ?"
				args = append(args, store.DecodeUid(opts.User))
			}
			oneUser = opts.User
		}

		if opts.Limit > 0 && opts.Limit < limit {
			limit = opts.Limit
		}
	}

	q += " LIMIT ? "
	args = append(args, limit)
	q, args = t.shared.ExpandQuery(q, args...)

	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	rows, err := t.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch subscriptions
	var sub types.Subscription
	var subs []types.Subscription
	var userId int64
	var modeWant, modeGiven []byte
	var lastSeen *time.Time = nil
	var userAgent string
	var public, trusted any

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
			&public,
			&trusted,
			&lastSeen,
			&userAgent,
			&sub.Private,
		); err != nil {
			break
		}

		sub.User = store.EncodeUid(userId).String()
		sub.SetPublic(public)
		sub.SetTrusted(trusted)
		sub.SetLastSeenAndUA(lastSeen, userAgent)
		sub.ModeWant.Scan(modeWant)
		sub.ModeGiven.Scan(modeGiven)
		subs = append(subs, sub)
	}

	if err == nil {
		err = rows.Err()
	}

	if err == nil && tcat == types.TopicCatP2P && len(subs) > 0 {
		// Swap public & lastSeen values of P2P topics as expected.
		if len(subs) == 1 {
			// The other user is deleted, nothing we can do.
			subs[0].SetPublic(nil)
			subs[0].SetTrusted(nil)
			subs[0].SetLastSeenAndUA(nil, "")

		} else {
			tmp := subs[0].GetPublic()
			subs[0].SetPublic(subs[1].GetPublic())
			subs[1].SetPublic(tmp)

			tmp = subs[0].GetTrusted()
			subs[0].SetTrusted(subs[1].GetTrusted())
			subs[1].SetTrusted(tmp)

			lastSeen := subs[0].GetLastSeen()
			userAgent = subs[0].GetUserAgent()
			subs[0].SetLastSeenAndUA(subs[1].GetLastSeen(), subs[1].GetUserAgent())
			subs[1].SetLastSeenAndUA(lastSeen, userAgent)
		}

		// Remove deleted and unneeded subscriptions
		if !keepDeleted || !oneUser.IsZero() {
			var xsubs []types.Subscription
			for i := range subs {
				if (subs[i].DeletedAt != nil && !keepDeleted) || (!oneUser.IsZero() && subs[i].Uid() != oneUser) {
					continue
				}
				xsubs = append(xsubs, subs[i])
			}
			subs = xsubs
		}
	}

	return subs, err
}

// OwnTopics loads a slice of topic names where the user is the owner.
func (t *Topics) OwnTopics(uid types.Uid) ([]string, error) {
	stmt := "SELECT name FROM topics WHERE owner= $1 "
	return t.getTopicNamesForUser(uid, stmt)
}

// ChannelsForUser loads a slice of topic names where the user is a channel
//
// reader and notifications (P) are enabled.
func (t *Topics) ChannelsForUser(uid types.Uid) ([]string, error) {
	stmt := `
	  SELECT topic FROM subscriptions 
		WHERE 
					user_id = $1 
			AND topic LIKE 'chn%' 
			AND POSITION('P' IN mode_want) > 0 
			AND POSITION('P' IN mode_given) > 0 
			AND deleted_at IS NULL
	`
	return t.getTopicNamesForUser(uid, stmt)
}

// Share() creates topic subscriptions
func (t *Topics) Share(shares []*types.Subscription) error {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	for _, sub := range shares {
		err = t.createSubscription(ctx, tx, sub, true)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// Delete deletes specified topic.
func (t *Topics) Delete(topic string, isChan, hard bool) error {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// If the topic is a channel, must try to delete subscriptions under both grpXXX and chnXXX names.
	args := []any{topic}
	if isChan {
		args = append(args, types.TopicsGrpToChn(topic))
	}

	if hard {
		// Delete subscriptions. If this is a channel, delete both group subscriptions and channel subscriptions.
		q, args := t.shared.ExpandQuery("DELETE FROM subscriptions WHERE topic IN (?) ", args)

		if _, err = tx.Exec(ctx, q, args...); err != nil {
			return err
		}

		if err = t.shared.MessageDeleteList(ctx, tx, topic, nil); err != nil {
			return err
		}

		if _, err = tx.Exec(ctx, "DELETE FROM topic_tags WHERE topic = $1", topic); err != nil {
			return err
		}

		if _, err = tx.Exec(ctx, "DELETE FROM topics WHERE name = $1", topic); err != nil {
			return err
		}

	} else {
		now := types.TimeNow()
		stmt := `UPDATE subscriptions SET updated_at = ?, deleted_at = ? WHERE topic IN (?)`
		q, args := t.shared.ExpandQuery(stmt, now, now, args)

		if _, err = tx.Exec(ctx, q, args); err != nil {
			return err
		}

		stmt = `UPDATE topics SET updated_at = $1, touched_at = $2, state = $3, state_at = $4 WHERE name=$5 `
		if _, err = tx.Exec(ctx, stmt, now, now, types.StateDeleted, now, topic); err != nil {
			return err
		}

	}
	return tx.Commit(ctx)
}

// UpdateOnMessage increments Topic's or User's SeqId value and updates TouchedAt timestamp.
func (t *Topics) UpdateOnMessage(topic string, msg *types.Message) error {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := `UPDATE topics SET seq_id = $1, touched_at = $2 WHERE name = $3`
	_, err := t.db.Exec(ctx, stmt, msg.SeqId, msg.CreatedAt, topic)
	return err
}

// Update() updates a topic record.
func (t *Topics) Update(topic string, update map[string]any) error {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	if t, u := update["touched_at"], update["updated_at"]; t == nil && u != nil {
		update["touched_at"] = u
	}

	cols, args := t.shared.UpdateByMap(update)
	stmt := "UPDATE topics SET " + strings.Join(cols, ",") + " WHERE name = ?"

	q, args := t.shared.ExpandQuery(stmt, args, topic)
	_, err = tx.Exec(ctx, q, args...)

	if err != nil {
		return err
	}

	// Tags are also stored in a separate table
	if tags := t.shared.ExtractTags(update); tags != nil {
		// First delete all user tags
		_, err = tx.Exec(ctx, "DELETE FROM topic_tags WHERE topic = $1", topic)
		if err != nil {
			return err
		}

		// Now insert new tags
		err = t.shared.AddTags(ctx, tx, "topic_tags", "topic", topic, tags, false)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// UpdateTopicOwner updates topic's owner
func (t *Topics) UpdateTopicOwner(topic string, newOwner types.Uid) error {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := `UPDATE topics SET owner = $1 WHERE name = $2`
	_, err := t.db.Exec(ctx, stmt, store.DecodeUid(newOwner), topic)

	return err
}

func (t *Topics) createTopic(ctx context.Context, tx pgx.Tx, topic *types.Topic) error {
	stmt := `
	  INSERT INTO topics(
			created_at,
			updated_at,
			touched_at,
			state,
			name,
			use_bt,
			owner,
			access,
			public,
			trusted,
			tags
		)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`
	_,
		err := tx.Exec(ctx,
		stmt,
		topic.CreatedAt,
		topic.UpdatedAt,
		topic.TouchedAt,
		topic.State,
		topic.ID,
		topic.UseBt,
		store.DecodeUid(types.ParseUid(topic.Owner)),
		topic.Access,
		t.utils.ToJSON(topic.Public),
		t.utils.ToJSON(topic.Trusted),
		topic.Tags,
	)
	if err != nil {
		return err
	}

	// Save topic's tags to a separate table to make topic findable.
	return t.shared.AddTags(ctx, tx, "topic_tags", "topic", topic.ID, topic.Tags, false)
}

// createSubscription() If undelete = true - update subscription on duplicate key,
// otherwise ignore the duplicate.
func (t *Topics) createSubscription(ctx context.Context, tx pgx.Tx, sub *types.Subscription, undelete bool) error {
	isOwner := (sub.ModeGiven & sub.ModeWant).IsOwner()

	jpriv := t.utils.ToJSON(sub.Private)
	decUID := store.DecodeUid(types.ParseUid(sub.User))

	_, err2 := tx.Exec(ctx, "SAVEPOINT createSub")
	if err2 != nil {
		t.logger.Sugar().Infof("failed to create savepoint <subscriptions(createSubscription)>: %s", err2.Error())
	}

	stmt := `
	  INSERT INTO subscriptions(
			created_at,
			updated_at,
			deleted_at,
			user_id,
			topic,
			mode_want,
			mode_given,
			private
	  )
		VALUES($1, $2, NULL, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(
		ctx,
		stmt,
		sub.CreatedAt,
		sub.UpdatedAt,
		decUID,
		sub.Topic,
		sub.ModeWant.String(),
		sub.ModeGiven.String(),
		jpriv,
	)

	if err != nil && t.shared.IsDupe(err) {
		_, err2 = tx.Exec(ctx, "ROLLBACK TO SAVEPOINT createSub")
		if err2 != nil {
			t.logger.Sugar().Infof("failed to rollback savepoint <subscriptions(createSubscription)>: %s", err2.Error())
		}

		if undelete {
			stmt := `
			  UPDATE subscriptions SET 
					createdat = $1, 
				  updatedat = $2,
				  deletedat = NULL,
				  modeWant = $3,
				  modeGiven = $4,
				  del_id = 0,
          received_seq_id = 0,
				  read_seq_id = 0 
				WHERE topic = $5 AND user_id = $6
			`
			_,
				err = tx.Exec(ctx, stmt,
				sub.CreatedAt,
				sub.UpdatedAt,
				sub.ModeWant.String(),
				sub.ModeGiven.String(),
				sub.Topic,
				decUID,
			)

		} else {
			stmt := `
			  UPDATE subscriptions SET 
					created_at = $1,
					updated_at = $2,
					deleted_at = NULL,
					mode_want = $3,
					mode_given = $4,
				  del_id = 0,
				  received_seq_id = 0,
				  read_seq_id = 0,
				  private = $5
				WHERE topic=$6 AND user_id = $7
			`
			_,
				err = tx.Exec(ctx, stmt,
				sub.CreatedAt,
				sub.UpdatedAt,
				sub.ModeWant.String(),
				sub.ModeGiven.String(),
				jpriv,
				sub.Topic,
				decUID,
			)
		}

	} else {
		_, err2 = tx.Exec(ctx, "RELEASE SAVEPOINT createSub")
		if err2 != nil {
			t.logger.Sugar().Infof("failed to release savepoint <subscriptions(createSubscription)>: %s", err2.Error())
		}
	}

	if err == nil && isOwner {
		_, err = tx.Exec(ctx, "UPDATE topics SET owner = $1 WHERE name = $2", decUID, sub.Topic)
	}

	return err
}

// getTopicNamesForUser reads a slice of strings using provided query.
func (t *Topics) getTopicNamesForUser(uid types.Uid, sqlQuery string) ([]string, error) {
	ctx, cancel := t.utils.GetContext(time.Duration(t.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	rows, err := t.db.Query(ctx, sqlQuery, store.DecodeUid(uid))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	var name string
	for rows.Next() {
		if err = rows.Scan(&name); err != nil {
			break
		}
		names = append(names, name)
	}

	if err == nil {
		err = rows.Err()
	}

	return names, err
}
