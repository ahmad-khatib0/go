package shared

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5"
)

type Shared struct{}

func NewShared() *Shared {
	return &Shared{}
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

// Check if Postgres error is an Error Code: 1062. Duplicate entry ... for key ...
func (s *Shared) isDupe(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "SQLSTATE 23505")
}
