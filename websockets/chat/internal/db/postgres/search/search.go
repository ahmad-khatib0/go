package search

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Search struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
}

type SearchArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
}

func NewSearch(ua SearchArgs) *Search {
	return &Search{db: ua.DB, utils: ua.Utils, cfg: ua.Cfg, shared: ua.Shared}
}

// FindUsers searches for new contacts given a list of tags.
//
// Returns a list of users who match given tags, such as "email:jdoe@example.com" or "tel:+18003287448".
func (s *Search) FindUsers(user types.Uid, req [][]string, opt []string, activeOnly bool) ([]types.Subscription, error) {
	index := make(map[string]struct{})
	var args []any
	stateConstraint := ""

	if activeOnly {
		args = append(args, types.StateOK)
		stateConstraint = " u.state = ? AND "
	}

	allReq := s.utils.FlattenDoubleSlice(req)
	allTags := append(allReq, opt...)
	for _, tag := range allTags {
		index[tag] = struct{}{}
	}

	args = append(args, allTags)

	query := `
	  SELECT 
			u.id, 
			u.created_at,
			u.updated_at,
			u.access,
			u.public,
			u.trusted,
			u.tags, 
			COUNT(*) AS matches 
	  FROM users AS u 
		LEFT JOIN user_tags AS t ON t.user_id = u.id WHERE 
	` + stateConstraint + " t.tag IN (?) GROUP BY u.id, u.created_at, u.updated_at "

	if len(allReq) > 0 {
		query += " HAVING "
		first := true
		for _, reqDisjunction := range req {
			if len(reqDisjunction) > 0 {
				if !first {
					query += " AND "
				} else {
					first = false
				}
				// At least one of the tags must be present.
				query += " COUNT(t.tag IN (?) OR NULL) > = 1 "
				args = append(args, reqDisjunction)
			}
		}
	}

	query, args = s.shared.ExpandQuery(query+" ORDER BY matches DESC LIMIT ? ", args, s.cfg.MaxResults)
	ctx, cancel := s.utils.GetContext(time.Duration(s.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	// Get users matched by tags, sort by number of matches from high to low.
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userId int64
	var public, trusted any
	var access types.DefaultAccess
	var userTags types.StringSlice
	var ignored int
	var sub types.Subscription
	var subs []types.Subscription
	thisUser := store.DecodeUid(user)

	for rows.Next() {
		if err = rows.Scan(
			&userId,
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&access,
			&public,
			&trusted,
			&userTags,
			&ignored,
		); err != nil {
			subs = nil
			break
		}

		if userId == thisUser {
			// Skip the callee
			continue
		}

		sub.User = store.EncodeUid(userId).String()
		sub.SetPublic(public)
		sub.SetTrusted(trusted)
		sub.SetDefaultAccess(access.Auth, access.Anon)

		foundTags := make([]string, 0, 1)
		for _, tag := range userTags {
			if _, ok := index[tag]; ok {
				foundTags = append(foundTags, tag)
			}
		}

		sub.Private = foundTags
		subs = append(subs, sub)
	}

	if err == nil {
		err = rows.Err()
	}

	return subs, err
}

// Returns a list of topics with matching tags.
//
// Searching the 'topics.Tags' for the given tags using respective index.
func (s *Search) FindTopics(req [][]string, opt []string, activeOnly bool) ([]types.Subscription, error) {
	index := make(map[string]struct{})
	var args []any
	stateConstraint := ""

	if activeOnly {
		args = append(args, types.StateOK)
		stateConstraint = "t.state = ? AND "
	}

	allReq := s.utils.FlattenDoubleSlice(req)
	allTags := append(allReq, opt...)
	for _, tag := range allTags {
		index[tag] = struct{}{}
	}

	args = append(args, allTags)

	query := `
	  SELECT 
			t.id,
			t.name AS topic,
			t.created_at,
			t.updated_at,
			t.use_bt,
			t.access,
			t.public,
			t.trusted,
			t.tags,
			COUNT(*) AS matches 
	  FROM topics AS t 
		LEFT JOIN topic_tags AS tt ON t.name = tt.topic WHERE 
	` + stateConstraint + " tt.tag IN (?) GROUP BY t.id, t.name, t.created_at, t.updated_at, t.use_bt "

	if len(allReq) > 0 {
		query += " HAVING "
		first := true
		for _, reqDisjunction := range req {
			if len(reqDisjunction) > 0 {
				if !first {
					query += " AND "
				} else {
					first = false
				}

				// At least one of the tags must be present.
				query += " COUNT(tt.tag IN (?) OR NULL) > = 1 "
				args = append(args, reqDisjunction)
			}
		}
	}

	query, args = s.shared.ExpandQuery(query+" ORDER BY matches DESC LIMIT ? ", args, s.cfg.MaxResults)
	ctx, cancel := s.utils.GetContext(time.Duration(s.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var access types.DefaultAccess
	var public, trusted any
	var topicTags types.StringSlice
	var id int
	var ignored int
	var isChan bool
	var sub types.Subscription
	var subs []types.Subscription

	for rows.Next() {
		if err = rows.Scan(&id, &sub.Topic, &sub.CreatedAt, &sub.UpdatedAt, &isChan, &access,
			&public, &trusted, &topicTags, &ignored); err != nil {
			subs = nil
			break
		}

		if isChan {
			sub.Topic = types.TopicsGrpToChn(sub.Topic)
		}

		sub.SetPublic(public)
		sub.SetTrusted(trusted)
		sub.SetDefaultAccess(access.Auth, access.Anon)
		foundTags := make([]string, 0, 1)

		for _, tag := range topicTags {
			if _, ok := index[tag]; ok {
				foundTags = append(foundTags, tag)
			}
		}

		sub.Private = foundTags
		subs = append(subs, sub)
	}

	if err == nil {
		err = rows.Err()
	}

	if err != nil {
		return nil, err
	}

	return subs, nil
}
