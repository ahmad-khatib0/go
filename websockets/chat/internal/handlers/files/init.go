package files

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

func NewFilesHandler(db db.Adapter, l *logger.Logger) *FilesHandler {
	return &FilesHandler{logger: l, db: db}
}
