package server

import (
	"sync/atomic"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// Delete removes session from store.
func (ss *SessionStore) Delete(s *Session) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	delete(ss.sessCache, s.sid)
	if s.proto == LPOLL {
		ss.lru.Remove(s.lpTracker)
	}

	ss.stats.IntStatsSet(constants.StatsLiveSessions, int64(len(ss.sessCache)))
}

// Shutdown terminates sessionStore. No need to clean up.
// Don't send to clustered sessions, their servers are not being shut down.
func (ss *SessionStore) Shutdown() {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	shutdown := NoErrShutdown(types.TimeNow())
	for _, s := range ss.sessCache {
		if !s.isMultiplex() {
			_, data := s.serialize(shutdown)
			s.stopSession(data)
		}
	}

	// TODO: Consider broadcasting shutdown to other cluster nodes.
	ss.stats.Logger.Sugar().Infof("SessionStore shut down, sessions terminated: %d", len(ss.sessCache))
}

// EvictUser terminates all sessions of a given user.
func (ss *SessionStore) EvictUser(uid types.Uid, skipSid string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	// FIXME: this probably needs to be optimized. This may take very long
	// time if the node hosts 100000 sessions.
	evicted := NoErrEvicted("", "", types.TimeNow())
	evicted.AsUser = uid.UserId()
	for _, s := range ss.sessCache {
		if s.uid == uid && !s.isMultiplex() && s.sid != skipSid {
			_, data := s.serialize(evicted)
			s.stopSession(data)
			delete(ss.sessCache, s.sid)
			if s.proto == LPOLL {
				ss.lru.Remove(s.lpTracker)
			}
		}
	}

	ss.stats.IntStatsSet(constants.StatsLiveSessions, int64(len(ss.sessCache)))
}

// NodeRestarted removes stale sessions from a restarted cluster node.
//   - nodeName is the name of affected node
//   - fingerprint is the new fingerprint of the node.
func (ss *SessionStore) NodeRestarted(nodeName string, fingerprint int64) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	for _, s := range ss.sessCache {
		if !s.isMultiplex() || s.clnode.name != nodeName {
			continue
		}

		if s.clnode.fingerprint != fingerprint {
			s.stopSession(nil)
			delete(ss.sessCache, s.sid)
		}
	}

	ss.stats.IntStatsSet(constants.StatsLiveSessions, int64(len(ss.sessCache)))
}

// cleanUp is called when the session is terminated to perform resource cleanup.
func (s *Session) cleanUp(expired bool, ss *SessionStore) {
	atomic.StoreInt32(&s.terminating, 1) // mark the session as being terminated
	s.multi.purgeChannels()
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
		sub.done <- &ClientComMessage{Init: false}
	}
}

func (s *Session) stopSession(data any) {
	s.stop <- data
	s.maybeScheduleClusterWriteLoop()
}

func (s *Session) maybeScheduleClusterWriteLoop() {
	if s.multi != nil {
		s.multi.scheduleClusterWriteLoop()
		return
	}
	if s.isMultiplex() {
		s.scheduleClusterWriteLoop()
	}
}

func (s *Session) scheduleClusterWriteLoop() {
	if s.cluster != nil && s.cluster.proxyEventQueue != nil {
		s.cluster.proxyEventQueue.Schedule(func() { s.clusterWriteLoop(s.proxiedTopic) })
	}
}

// clusterWriteLoop implements write loop for multiplexing (proxy) session at a node which hosts master topic.
//
// The session is a multiplexing session, i.e. it handles requests for multiple sessions at origin.
func (s *Session) clusterWriteLoop(forTopic string) {
	terminate := true

	defer func() {
		if terminate {
			s.closeRPC()
			s.sessStore.Delete(s)
			s.inflightReqs = nil
			s.unsubAll()
		}
	}()

	for {
		select {
		case msg, ok := <-s.send:
			if !ok || s.clnode.endpoint == nil {
				// channel closed
				return
			}

			srvMsg := msg.(*ServerComMessage)
			response := &ClusterResp{SrvMsg: srvMsg}

			if srvMsg.sess == nil {
				response.OrigSid = "*"
			} else {
				response.OrigReqType = srvMsg.sess.proxyReq
				response.OrigSid = srvMsg.sess.sid
				srvMsg.AsUser = srvMsg.sess.uid.UserId()

				switch srvMsg.sess.proxyReq {
				case
					ProxyReqJoin,
					ProxyReqLeave,
					ProxyReqMeta,
					ProxyReqBgSession,
					ProxyReqMeUserAgent,
					ProxyReqCall:
					// Do nothing
				case ProxyReqBroadcast, ProxyReqNone:
					if srvMsg.Data != nil || srvMsg.Pres != nil || srvMsg.Info != nil {
						response.OrigSid = "*"
					} else if srvMsg.Ctrl == nil {
						s.logger.Warn(
							"session: request type not set in clusterWriteLoop: " +
								s.sid +
								srvMsg.Describe() +
								"src_sid:" +
								srvMsg.sess.sid,
						)
					}

				default:
					s.logger.Sugar().Panicf("cluster: unknown request type in clusterWriteLoop %+v", srvMsg.sess.proxyReq)
				}
			}

			srvMsg.RcptTo = forTopic
			response.RcptTo = forTopic
			if err := s.clnode.masterToProxyAsync(response); err != nil {
				s.logger.Sugar().Infof("cluster: response to proxy failed \"%s\": %s", s.sid, err.Error())
				return
			}

		case msg := <-s.stop:
			if msg == nil {
				// Terminating multiplexing session.
				return
			}

		// There are two cases of msg != nil:
		//  * user is being deleted
		//  * node shutdown
		// In both cases the msg does not need to be forwarded to the proxy.

		case <-s.detach:
			return

		default:
			terminate = false
			return
		}
	}
}

// Proxied session is being closed at the Master node.
func (s *Session) closeRPC() {
	if s.isMultiplex() {
		s.logger.Info("cluster: session proxy closed: " + s.sid)
	}
}
