package store

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	dt "github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	mt "github.com/ahmad-khatib0/go/websockets/chat/internal/media/types"
	st "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/validate"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type StoreArgs struct {
	Cfg      *config.Config
	Logger   *logger.Logger
	WorkerID uint
}

type Store struct {
	adp              dt.Adapter
	logger           *logger.Logger
	cfg              *config.Config
	UidGen           *st.UidGenerator
	authHandlerNames map[string]string
	authHandlers     map[string]types.AuthHandler
	validators       map[string]validate.Validator
	mediaHandlers    map[string]mt.Handler
	mediaHandler     mt.Handler
}
