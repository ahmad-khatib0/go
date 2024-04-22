package store

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	dt "github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/validate"
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
	authHandlers     map[string]types.AuthHandler
	adp              dt.Adapter
	cfg              *config.Config
	validators       map[string]validate.Validator
}
