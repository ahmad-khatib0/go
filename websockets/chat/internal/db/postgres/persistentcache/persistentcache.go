package persistentcache

import (
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersistentCache struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
	uGen   *types.UidGenerator
}

type PersistentCacheArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
	UGen   *types.UidGenerator
}

func NewPersistentCache(pc PersistentCacheArgs) *PersistentCache {
	return &PersistentCache{db: pc.DB, utils: pc.Utils, cfg: pc.Cfg, shared: pc.Shared}
}

// Get reads a persistet cache entry.
func (pc *PersistentCache) Get(key string) (string, error) {
	ctx, cancel := pc.utils.GetContext(time.Duration(pc.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	var value string
	stmt := `SELECT "value" FROM kvmeta WHERE "key" = $1 LIMIT 1`

	if err := pc.db.QueryRow(ctx, stmt, key).Scan(&value); err != nil {
		if err == pgx.ErrNoRows {
			return "", types.ErrNotFound
		}
		return "", err
	}
	return value, nil
}

// Upsert creates or updates a persistent cache entry.
func (pc *PersistentCache) Upsert(key string, value string, failOnDuplicate bool) error {
	if strings.Contains(key, "%") {
		// Do not allow % in keys: it interferes with LIKE query.
		return types.ErrMalformed
	}

	ctx, cancel := pc.utils.GetContext(time.Duration(pc.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	var action string

	if failOnDuplicate {
		action = "INSERT"
	} else {
		action = "REPLACE"
	}

	stmt := action + ` INTO kvmeta("key", created_at, "value") VALUES($1, $2, $3)`

	_, err := pc.db.Exec(ctx, stmt, key, types.TimeNow(), value)
	if pc.shared.IsDupe(err) {
		return types.ErrDuplicate
	}

	return err
}

// Delete deletes one persistent cache entry.
func (pc *PersistentCache) Delete(key string) error {
	ctx, cancel := pc.utils.GetContext(time.Duration(pc.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	_, err := pc.db.Exec(ctx, `DELETE FROM kvmeta WHERE "key" = $1`, key)
	return err
}

// Expire expires old entries with the given key prefix.
func (pc *PersistentCache) Expire(keyPrefix string, olderThan time.Time) error {
	if keyPrefix == "" {
		return types.ErrMalformed
	}

	ctx, cancel := pc.utils.GetContext(time.Duration(pc.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	stmt := `DELETE FROM kvmeta WHERE "key" LIKE $1 AND created_at < $2 `
	_, err := pc.db.Exec(ctx, stmt, keyPrefix+"%", olderThan)

	return err
}
