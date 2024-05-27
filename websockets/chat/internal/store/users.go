package store

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// UpdateLastSeen updates LastSeen and UserAgent.
func (s *Store) UsersUpdateLastSeen(uid types.Uid, userAgent string, when time.Time) error {
	return s.adp.Users().Update(uid, map[string]interface{}{"LastSeen": when, "UserAgent": userAgent})
}
