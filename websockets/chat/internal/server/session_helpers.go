package server

import (
	"encoding/json"
)

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
		msg := PbServSerialize(msg)
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
