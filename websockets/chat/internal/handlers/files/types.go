package files

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	mt "github.com/ahmad-khatib0/go/websockets/chat/internal/media/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type FilesHandler struct {
	db     types.Adapter
	media  mt.Handler
	logger *logger.Logger
}
