package files

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Files struct {
	db *pgxpool.Pool
}

// DeleteUnused implements db.Files.
func (f *Files) DeleteUnused(olderThan time.Time, limit int) error {
	panic("unimplemented")
}

// LinkAttachments implements db.Files.
func (f *Files) LinkAttachments(topic string, userId types.Uid, msgId types.Uid, fids []string) error {
	panic("unimplemented")
}

type FilesArgs struct {
	DB *pgxpool.Pool
}

func NewFiles(fa FilesArgs) *Files {
	return &Files{db: fa.DB}
}
