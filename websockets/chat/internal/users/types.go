package users

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	dt "github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type Users struct {
	logger *logger.Logger
	db     dt.Adapter
}

// CredValidator holds additional config params for a credential validator.
type CredValidator struct {
	// AuthLevel(s) which require this validator.
	RequiredAuthLvl []types.Level
	AddToTags       bool
}

func NewUser(adapter dt.Adapter, l *logger.Logger) *Users {
	return &Users{db: adapter, logger: l}
}
