package store

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type StoreArgs struct {
	Cfg      *config.Config
	Logger   *logger.Logger
	WorkerID uint
}

type Store struct {
	logger *logger.Logger
	// Logical auth handler names (supplied by config)
	authHandlerNames map[string]string
	authHandlers     map[string]auth.AuthHandler
	adp              db.Adapter
	cfg              *config.Config
}
