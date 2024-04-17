package files

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type FilesHandler struct {
	db     db.Adapter
	logger *logger.Logger
}
