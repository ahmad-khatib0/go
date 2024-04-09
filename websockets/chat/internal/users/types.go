package users

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type Users interface {
	// GetUnvalidated(lastUpdatedBefore time.Time, limit int) ([]types.Uid, error)
	// Delete(id types.Uid, hard bool) error
	// InitUsersGarbageCollection() runs every 'period' and deletes up to 'blockSize'
	//
	// stale unvalidated user accounts which have been last updated at least 'minAccountAgeHours' hours.
	//
	// Returns channel which can be used to stop the process.
	InitUsersGarbageCollection(period time.Duration, blockSize, minAccountAgeHours int)
}

type users struct {
	logger *logger.Logger
	db     db.Adapter
}

// CredValidator holds additional config params for a credential validator.
type CredValidator struct {
	// AuthLevel(s) which require this validator.
	RequiredAuthLvl []auth.Level
	AddToTags       bool
}

func NewUser(adapter db.Adapter, l *logger.Logger) Users {
	return &users{db: adapter, logger: l}
}
