package shared

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type Shared struct {
	utils *utils.Utils
	uGen  *types.UidGenerator
}

type SharedArgs struct {
	Utils *utils.Utils
	UGen  *types.UidGenerator
}

func NewShared(sa SharedArgs) *Shared {
	return &Shared{utils: sa.Utils, uGen: sa.UGen}
}

func (s *Shared) AddTags(ctx context.Context, tx pgx.Tx, table, keyName string, keyVal any, tags []string, ignoreDups bool) error {

	if len(table) == 0 {
		return nil
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s, tag) VALUES($1, $2)", table, keyName)
	for t := range tags {
		if _, err := tx.Exec(ctx, sql, keyVal, t); err != nil {
			if s.IsDupe(err) {
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
		res, err = tx.Exec(ctx, "DELETE FROM devices WHERE user_id = $1", s.uGen.DecodeUid(uid))
	} else {
		stmt := "DELETE FROM devices WHERE user_id = $1 AND hash = $2"
		res, err = tx.Exec(ctx, stmt, s.uGen.DecodeUid(uid), s.DeviceHasher(deviceID))
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
		stmt := "DELETE FROM subscriptions WHERE user_id = $1;"
		_, err = tx.Exec(ctx, stmt, s.uGen.DecodeUid(uid))
	} else {
		now := types.TimeNow()
		stmt := "UPDATE subscriptions SET updated_at = $1, deleted_at = $2 WHERE user_id = $3 AND deleted_at IS NULL;"
		_, err = tx.Exec(ctx, stmt, now, now, s.uGen.DecodeUid(uid))
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
	args := []any{s.uGen.DecodeUid(uid)}

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

func (s *Shared) MessageDeleteList(ctx context.Context, tx pgx.Tx, topic string, toDel *types.DelMessage) error {
	var err error

	if toDel == nil {
		// Whole topic is being deleted, thus also deleting all messages.
		_, err = tx.Exec(ctx, "DELETE FROM delete_logs WHERE topic = $1", topic)

		if err == nil {
			_, err = tx.Exec(ctx, "DELETE FROM messages WHERE topic = $1", topic)
		}
		// file_messages_links will be deleted also because of ON DELETE CASCADE

	} else {
		// Only some messages are being deleted

		// Start with making log entries
		forUser := s.DecodeUidString(toDel.DeletedFor)

		// Counter of deleted messages
		for _, rng := range toDel.SeqIdRanges {
			if rng.Hi == 0 {
				// Dellog must contain valid Low and *Hi*.
				rng.Hi = rng.Low + 1
			}

			stmt := `INSERT INTO delete_logs(topic, deleted_for, del_id, low, hi) VALUES($1, $2, $3, $4, $5)`
			if _, err = tx.Exec(ctx, stmt, topic, forUser, toDel.DelId, rng.Low, rng.Hi); err != nil {
				break
			}
		}

		if err == nil && toDel.DeletedFor == "" {
			// Hard-deleting messages requires updates to the messages table

			where := " m.topic = ? AND "
			args := []any{topic}

			if len(toDel.SeqIdRanges) > 1 || toDel.SeqIdRanges[0].Hi == 0 {
				seqRange := []int{}

				for _, r := range toDel.SeqIdRanges {
					if r.Hi == 0 {
						seqRange = append(seqRange, r.Low)

					} else {
						for i := r.Low; i < r.Hi; i++ {
							seqRange = append(seqRange, i)
						}
					}
				}

				args = append(args, seqRange)
				where += " m.seq_id IN (?) "

			} else {
				// Optimizing for a special case of single range low..hi.
				where += " m.seq_id BETWEEN ? AND ? "
				// MySQL's BETWEEN is inclusive-inclusive thus decrement Hi by 1.
				args = append(args, toDel.SeqIdRanges[0].Low, toDel.SeqIdRanges[0].Hi-1)
			}

			where += " AND m.deleted_at IS NULL "
			stmt := ` 
				DELETE FROM file_messages_links AS fml USING messages AS m 
				WHERE m.id = fml.msg_id AND 
			` + where

			query, newargs := s.ExpandQuery(stmt, args...)
			_, err = tx.Exec(ctx, query, newargs...)
			if err != nil {
				return err
			}

			stmt = `UPDATE messages AS m SET deleted_at = ?, del_id = ?, head = NULL, content = NULL WHERE ` + where
			query, newargs = s.ExpandQuery(stmt, types.TimeNow(), toDel.DelId, args)
			_, err = tx.Exec(ctx, query, newargs...)
		}
	}

	return err
}

// ExpandQuery() converts query strings like: 'where arg = ? AND arg2 = ? ...' to
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

	if val := update["tags"]; val != nil {
		tags, _ = val.(types.StringSlice)
	}

	return []string(tags)
}

func (s *Shared) DecodeUidString(str string) int64 {
	uid := types.ParseUid(str)
	return s.uGen.DecodeUid(uid)
}

func (s *Shared) DeviceHasher(devID string) string {
	// Generate custom key as [64-bit hash of device id] to ensure predictable length of the key
	hasher := fnv.New64()
	hasher.Write([]byte(devID))
	return strconv.FormatUint(uint64(hasher.Sum64()), 16)
}

// Check if Postgres error is an Error Code: 1062. Duplicate entry ... for key ...
func (s *Shared) IsDupe(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "SQLSTATE 23505")
}
