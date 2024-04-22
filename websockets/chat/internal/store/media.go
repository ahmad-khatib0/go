package store

import (
	"errors"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/media/types"
)

// GetMediaHandler returns default media handler.
func (s *Store) GetMediaHandler() types.Handler {
	return s.mediaHandler
}

// UseMediaHandler sets specified media handler as default.
func (s *Store) SetDefaultMediaHandler(name string, cfg interface{}) error {
	s.mediaHandler = s.mediaHandlers[name]
	if s.mediaHandler == nil {
		return errors.New("unknown media handler SetDefaultMediaHandler ")
	}
	return s.mediaHandler.Init(cfg)
}
