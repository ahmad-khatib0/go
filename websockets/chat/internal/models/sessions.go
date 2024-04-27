package models

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// SessionProto is the type of the wire transport.
type SessionProto int

// Constants defining individual types of wire transports.
const (
	// NONE is undefined/not set.
	NONE SessionProto = iota
	// WEBSOCK represents websocket connection.
	WEBSOCK
	// LPOLL represents a long polling connection.
	LPOLL
	// GRPC is a gRPC connection
	GRPC
	// PROXY is temporary session used as a proxy at master node.
	PROXY
	// MULTIPLEX is a multiplexing session representing a connection from proxy topic to master.
	MULTIPLEX
)

type Session interface {
	ProxyReq() ProxyReqType
	Sid() string
	Uid() types.Uid
}

type Subscription interface{}

type SessionUpdate interface{}

type SessionStore interface {
	Delete(s Session)
	EvictUser(uid types.Uid, skipSid string)
	Get(sid string) *Session
	NewSession(conn any, sid string) (*Session, int)
	NodeRestarted(nodeName string, fingerprint int64)
	Range(f func(sid string, s *Session) bool)
	Shutdown()
}
