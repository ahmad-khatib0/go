package session

import (
	"container/list"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	st "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

type SessionArgs struct {
	Lifetime time.Duration
	Stats    *stats.Stats
	UGen     *st.UidGenerator
}

// NewSessionStore initializes a session store.
func NewSessionStore(sa SessionArgs) models.SessionStore {
	sa.Stats.RegisterInt(constants.StatsLiveSessions)
	sa.Stats.RegisterInt(constants.StatsTotalSessions)

	return SessionStore{
		lru:       list.New(),
		lifeTime:  sa.Lifetime,
		sessCache: make(map[string]*Session),
		stats:     sa.Stats,
	}

}
