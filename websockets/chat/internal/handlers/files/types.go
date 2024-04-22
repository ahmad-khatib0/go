package files

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type FilesHandler struct {
	db     types.Adapter
	logger *logger.Logger
}
