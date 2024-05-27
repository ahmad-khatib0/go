package server

import (
	"encoding/json"
)

func (s *Session) addSub(topic string, sub *Subscription) {
	if s.multi != nil {
		s.multi.addSub(topic, sub)
		return
	}
	s.subsLock.Lock()

	// Sessions that serve as an interface between proxy topics and their masters (proxy sessions)
	// may have only one subscription, that is, to its master topic.
	// Normal sessions may be subscribed to multiple topics.

	if !s.isMultiplex() || s.countSub() == 0 {
		s.subs[topic] = sub
	}
	s.subsLock.Unlock()
}

func (s *Session) getSub(topic string) *Subscription {
	// Don't check s.multi here. Let it panic if called for proxy session.

	s.subsLock.RLock()
	defer s.subsLock.RUnlock()

	return s.subs[topic]
}

func (s *Session) delSub(topic string) {
	if s.multi != nil {
		s.multi.delSub(topic)
		return
	}
	s.subsLock.Lock()
	delete(s.subs, topic)
	s.subsLock.Unlock()
}

func (s *Session) countSub() int {
	if s.multi != nil {
		return s.multi.countSub()
	}
	return len(s.subs)
}

// Indicates whether this session is a local interface for a remote proxy topic.
// It multiplexes multiple sessions.
func (s *Session) isMultiplex() bool {
	return s.proto == MULTIPLEX
}

// Indicates whether this session is a short-lived proxy for a remote session.
func (s *Session) isProxy() bool {
	return s.proto == PROXY
}

// Cluster session: either a proxy or a multiplexing session.
func (s *Session) isCluster() bool {
	return s.isProxy() || s.isMultiplex()
}

func (s *Session) serialize(msg *ServerComMessage) (int, any) {
	if s.proto == GRPC {
		msg := pbServSerialize(msg)
		// TODO: calculate and return the size of `msg`.
		return -1, msg
	}

	if s.isMultiplex() {
		// No need to serialize the message to bytes within the cluster.
		return -1, msg
	}

	out, _ := json.Marshal(msg)
	return len(out), out
}
