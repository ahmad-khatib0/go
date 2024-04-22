package files

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

func NewFilesHandler(db types.Adapter, l *logger.Logger) *FilesHandler {
	return &FilesHandler{logger: l, db: db}
}
