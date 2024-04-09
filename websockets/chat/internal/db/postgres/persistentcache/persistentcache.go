package persistentcache

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PersistentCache struct {
	db *pgxpool.Pool
}

// Delete implements db.PersistentCache.
func (p *PersistentCache) Delete(key string) error {
	panic("unimplemented")
}

// Expire implements db.PersistentCache.
func (p *PersistentCache) Expire(keyPrefix string, olderThan time.Time) error {
	panic("unimplemented")
}

// Get implements db.PersistentCache.
func (p *PersistentCache) Get(key string) (string, error) {
	panic("unimplemented")
}

// Upsert implements db.PersistentCache.
func (p *PersistentCache) Upsert(key string, value string, failOnDuplicate bool) error {
	panic("unimplemented")
}

type PersistentCacheArgs struct {
	DB *pgxpool.Pool
}

func NewPersistentCache(ua PersistentCacheArgs) *PersistentCache {
	return &PersistentCache{db: ua.DB}
}
