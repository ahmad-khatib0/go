package shared

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type Shared struct {
	utils *utils.Utils
}

type SharedArgs struct {
	Utils *utils.Utils
}

func NewShared(sa SharedArgs) *Shared {
	return &Shared{utils: sa.Utils}
}

func (s *Shared) AddTags(ctx context.Context, tx pgx.Tx, table, keyName string, keyVal any, tags []string, ignoreDups bool) error {

	if len(table) == 0 {
		return nil
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s, tag) VALUES($1, $2)", table, keyName)
	for t := range tags {
		if _, err := tx.Exec(ctx, sql, keyVal, t); err != nil {
			if s.isDupe(err) {
				if ignoreDups {
					continue
				}
				return types.ErrDuplicate
			}
			return err
		}
	}

	return nil
}

func (s *Shared) DeviceDelete(ctx context.Context, tx pgx.Tx, uid types.Uid, deviceID string) error {
	var err error
	var res pgconn.CommandTag
	if deviceID == "" {
		res, err = tx.Exec(ctx, "DELETE FROM devices WHERE user_id = $1", store.DecodeUid(uid))
	} else {
		stmt := "DELETE FROM devices WHERE user_id = $1 AND hash = $2"
		res, err = tx.Exec(ctx, stmt, store.DecodeUid(uid), s.deviceHasher(deviceID))
	}

	if err == nil {
		if count := res.RowsAffected(); count == 0 {
			err = types.ErrNotFound
		}
	}

	return err
}

// SubDelForUser marks user's subscriptions as deleted.
func (s *Shared) SubDelForUser(ctx context.Context, tx pgx.Tx, uid types.Uid, hard bool) error {
	var err error

	if hard {
		stmt := "DELETE FROM subscriptions WHERE userid = $1;"
		_, err = tx.Exec(ctx, stmt, store.DecodeUid(uid))
	} else {
		now := types.TimeNow()
		stmt := "UPDATE subscriptions SET updated_at = $1, deleted_at = $2 WHERE user_id = $3 AND deleted_at IS NULL;"
		_, err = tx.Exec(ctx, stmt, now, now, store.DecodeUid(uid))
	}

	return err
}

// CredDel deletes given validation method or all methods of the given user.
//
// 1. If user is being deleted, hard-delete all records (method == "")
//
// 2. If one value is being deleted:
//
// 2.1 Delete it if it's valiated or if there were no attempts at validation
//
//	(otherwise it could be used to circumvent the limit on validation attempts).
//
// 2.2 In that case mark it as soft-deleted.
func (s *Shared) CredDel(ctx context.Context, tx pgx.Tx, uid types.Uid, method, value string) error {
	constraints := " WHERE user_id = ? "
	args := []any{store.DecodeUid(uid)}

	if method != "" {
		constraints += " AND method = ? "
		args = append(args, method)

		if value != "" {
			constraints += " AND value = ? "
			args = append(args, value)
		}
	}

	whereStmt, _ := s.ExpandQuery(constraints, args...)

	var err error
	var res pgconn.CommandTag

	if method == "" {
		// Case 1
		res, err = tx.Exec(ctx, "DELETE FROM credentials"+whereStmt, args...)
		if err == nil {
			if count := res.RowsAffected(); count == 0 {
				err = types.ErrNotFound
			}
		}
		return err
	}

	// Case 2.1
	res, err = tx.Exec(ctx, "DELETE FROM credentials "+whereStmt+" AND (done = TRUE OR retries = 0)", args...)
	if err != nil {
		return err
	}
	if count := res.RowsAffected(); count > 0 {
		return nil
	}

	// Case 2.2
	// (note the order: types.TimeNow and then the rest of args)
	query, args := s.ExpandQuery("UPDATE credentials SET deleted_at = ? "+constraints, types.TimeNow(), args)
	res, err = tx.Exec(ctx, query, args...)
	if err == nil {
		// FIX: is it count >= 0   or  count == 0 ?
		if count := res.RowsAffected(); count >= 0 {
			err = types.ErrNotFound
		}
	}

	return err
}

// AppendQuery() converts query strings like: 'where arg = ? AND arg2 = ? ...' to
//
// a query that can be used with pgx library like: 'where arg = $1 AND arg2 = $2 ...'
func (s *Shared) ExpandQuery(query string, args ...any) (string, []any) {
	var expandedArgs []any
	var expandedQuery string

	if len(args) != strings.Count(query, "?") {
		// INFO: flatten the slice of args, why? because i.e maybe one of the args is slice
		// and it contains more than one value that corresponds to => ?, so it's required
		// in this case to extract each => ? to its value in a flat way
		args = s.utils.FlattenMap(args)
	}

	// prepare the query to be used
	expandedQuery, expandedArgs, _ = sqlx.In(query, args)

	placeholders := make([]string, len(expandedArgs))
	for i := range expandedArgs {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		expandedQuery = strings.Replace(expandedQuery, "?", placeholders[i], 1)
	}
	return expandedQuery, expandedArgs
}

// Convert update to a list of columns and arguments.
func (s *Shared) UpdateByMap(update map[string]any) (cols []string, args []any) {
	for col, arg := range update {
		col = strings.ToLower(col)
		if col == "public" || col == "trusted" || col == "private" {
			arg = s.utils.ToJSON(arg)
		}
		cols = append(cols, col+"=?")
		args = append(args, arg)
	}
	return
}

// If Tags field is updated, get the tags so tags table cab be updated too.
func (s *Shared) ExtractTags(update map[string]any) []string {
	var tags []string
	if val := update["Tags"]; val != nil {
		tags, _ = val.(types.StringSlice)
	}

	return []string(tags)
}

func (s *Shared) deviceHasher(devID string) string {
	// Generate custom key as [64-bit hash of device id] to ensure predictable length of the key
	hasher := fnv.New64()
	hasher.Write([]byte(devID))
	return strconv.FormatUint(uint64(hasher.Sum64()), 16)
}

// Check if Postgres error is an Error Code: 1062. Duplicate entry ... for key ...
func (s *Shared) isDupe(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "SQLSTATE 23505")
}
