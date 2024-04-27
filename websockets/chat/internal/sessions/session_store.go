/******************************************************************************
 *
 *  Description:
 *
 *  Session management.
 *
 *****************************************************************************/
package session

import (
	"net/http"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// NewSession creates a new session and saves it to the session store.
func (ss *SessionStore) NewSession(conn any, sid string) (*Session, int) {
	var s Session

	if sid == "" {
		s.sid = ss.uGen.GetUidString()
	} else {
		s.sid = sid
	}

	ss.lock.Lock()
	if _, found := ss.sessCache[s.sid]; found {
		log.Fatal().Msg("ERROR! duplicate session ID" + s.sid)
	}
	ss.lock.Unlock()

	switch c := conn.(type) {
	case *websocket.Conn:
		s.proto = models.WEBSOCK
		s.ws = c
	case http.ResponseWriter:
		s.proto = models.LPOLL
		// no need to store c for long polling, it changes with every request
	case models.ClusterNode:
		s.proto = models.MULTIPLEX
		s.clnode = c
	case chat.Node_MessageLoopServer:
		s.proto = models.GRPC
		s.grpcCNode = c
	default:
		log.Panic().Msgf("session: unknown connection type %+v", conn)
	}

	s.subs = make(map[string]*Subscription)
	s.send = make(chan any, sendQueueLimit+32)
	s.stop = make(chan any, 1) // Buffered by 1 just to make it non-blocking
	s.detach = make(chan string, 64)

	s.bkgTimer = time.NewTimer(time.Hour)
	s.bkgTimer.Stop()

	// Make sure at most 1 request is modifying session/topic state at any time.
	// TODO: use Mutex & CondVar?
	s.inflightReqs = newBoundedWaitGroup(1)
	s.lastTouched = time.Now()
	ss.lock.Lock()

	if s.proto == models.LPOLL {
		// Only LP sessions need to be sorted by last active
		s.lpTracker = ss.lru.PushFront(&s)
	}

	ss.sessCache[s.sid] = &s

	// Expire stale long polling sessions: ss.lru contains only long polling sessions.
	// If ss.lru is empty this is a noop.
	var expired []*Session
	expire := s.lastTouched.Add(-ss.lifeTime)
	for el := ss.lru.Back(); el != nil; el = ss.lru.Back() {
		sess := el.Value.(*Session)
		if sess.lastTouched.Before(expire) {
			ss.lru.Remove(el)
			delete(ss.sessCache, sess.sid)
			expired = append(expired, sess)
		} else {
			break // don't need to traverse further
		}
	}

	sessCount := len(ss.sessCache)
	ss.stats.IntStatsSet(constants.StatsLiveSessions, int64(sessCount))
	ss.stats.IntStatsInc(constants.StatsTotalSessions, 1)

	ss.lock.Unlock()

	// Deleting long polling sessions.
	for _, sess := range expired {
		// This locks the session. Thus cleaning up outside of the
		// sessionStore lock. Otherwise deadlock.
		// sess.cl
	}

}
