package session

import (
	"sync/atomic"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
)

// Delete removes session from store.
func (ss *SessionStore) Delete(s *Session) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	delete(ss.sessCache, s.sid)
	if s.proto == models.LPOLL {
		ss.lru.Remove(s.lpTracker)
	}

	ss.stats.IntStatsSet(constants.StatsLiveSessions, int64(len(ss.sessCache)))
}

// cleanUp is called when the session is terminated to perform resource cleanup.
func (s *Session) cleanUp(expired bool, ss *SessionStore) {
	atomic.StoreInt32(&s.terminating, 1) // mark the session as being terminated
	s.uulit.purgeChannels()
	s.inflightReqs.Wait()
	s.inflightReqs = nil

	if !expired {
		s.sessionStoreLock.Lock()
		ss.Delete(s)
		s.sessionStoreLock.Unlock()
	}

	s.background = false
	s.bkgTimer.Stop()
	// s.
}

func (s *Session) purgeChannels() {
	for len(s.send) > 0 {
		<-s.send
	}

	for len(s.stop) > 0 {
		<-s.stop
	}

	for len(s.detach) > 0 {
		<-s.detach
	}
}

// Inform topics that the session is being terminated.
//
// No need to check for s.multi because it's not called for PROXY sessions.
func (s *Session) unsubAll() {
	s.subLock.RLock()
	defer s.subLock.Unlock()

	for _, sub := range s.subs {
		// sub.done is the same as topic.unreg, The whole session is being dropped; ClientComMessage is a wrapper
		// for session, ClientComMessage.init is false. keep redundant init: false so it can be searched for.
		sub.Done <- models.ClientComMessage{Init: false}
	}
}
