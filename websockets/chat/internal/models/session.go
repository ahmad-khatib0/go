package models

import (
	"container/list"
	"sync"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

	"github.com/gorilla/websocket"
)

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
	// MULTIPLEX is a multiplexing session reprsenting a connection from proxy topic to master.
	MULTIPLEX
)

// Session represents a single WS connection or a long polling session.
// A user may have multiple sessions.
type Session struct {
	// protocol - NONE (unset), WEBSOCK, LPOLL, GRPC, PROXY, MULTIPLEX
	Proto SessionProto
	// Session ID
	SID string

	// Websocket. Set only for websocket sessions.
	Ws *websocket.Conn

	// Pointer to session's record in sessionStore. Set only for Long Poll sessions.
	LpTracker *list.Element

	// gRPC handle. Set only for gRPC clients.
	GrpcCNode chat.Node_MessageLoopServer

	// Reference to the cluster node where the session has originated. Set only for cluster RPC sessions.
	// CLNode *ClusterNode

	// Reference to multiplexing session. Set only for proxy sessions.
	Mulit        *Session
	ProxiedTopic string

	// IP address of the client. For long polling this is the IP of the last poll.
	RemoteAddr string

	// User agent, a string provived by an authenticated client in {login} packet.
	UserAgent string

	// Protocol version of the client: ((major & 0xff) << 8) | (minor & 0xff).
	Ver int

	// Device ID of the client
	DeviceID string
	// Platform: web, ios, android
	Platform string

	// Human language of the client
	Lang string

	// Country code of the client
	CountryCode string

	// ID of the current user. Could be zero if session is not authenticated
	// or for multiplexing sessions.
	Uid types.Uid

	// Authentication level - NONE (unset), ANON, AUTH, ROOT.
	AuthLvl AuthLevel

	// Time when the long polling session was last refreshed
	LastTouched time.Time

	// Time when the session received any packer from client
	LastAction int64

	// Timer which triggers after some seconds to mark background session as foreground.
	BkgTimer *time.Time

	// Number of subscribe/unsubscribe requests in flight.
	InflightReqs *BoundedWaitGroup

	// Synchronizes access to session store in cluster mode:
	// subscribe/unsubscribe replies are asynchronous.
	SessionStoreLock sync.Mutex

	// Indicates that the session is terminating.
	// After this flag's been flipped to true, there must not be any more writes
	// into the session's send channel.
	// Read/written atomically.
	// 0 = false
	// 1 = true
	Terminating int32

	// Background session: subscription presence notifications and online status are delayed.
	Background bool

	// Outbound mesages, buffered.
	// The content must be serialized in format suitable for the session.
	Send chan any

	// Channel for shutting down the session, buffer 1.
	// Content in the same format as for 'send'
	Stop chan any

	// detach - channel for detaching session from topic, buffered.
	// Content is topic name to detach from.
	Detach chan string

	// Map of topic subscriptions, indexed by topic name.
	// Don't access directly. Use getters/setters.
	Subs map[string]*Subscription

	// Mutex for subs access: both topic go routines and network go routines access subs concurrently.
	SubLock sync.RWMutex

	// Needed for long polling and grpc.
	Lock sync.Mutex

	// Field used only in cluster mode by topic master node.

	// ProxyReq ProxyReqType
}

// SessionStore holds live sessions. Long polling sessions are stored in a linked list with
//
// most recent sessions on top. In addition all sessions are stored in a map indexed by session ID.
type SessionStore struct {
	Lock sync.Mutex

	// Support for long polling sessions: a list of sessions sorted by last access time.
	// Needed for cleaning abandoned sessions.
	Lru      *list.List
	LifeTime time.Duration

	// All sessions indexed by session ID
	SessCache map[string]*Session
}

// Subscription is a mapper of sessions to topics.
type Subscription struct {
	// Channel to communicate with the topic, copy of Topic.clientMsg
	Broadcast chan<- *ClientComMessage

	// Session sends a signal to Topic when this session is unsubscribed This is a copy of Topic.unreg
	Done chan<- *ClientComMessage

	// Channel to send {meta} requests, copy of Topic.meta
	Meta chan<- *ClientComMessage

	// Channel to ping topic with session's updates, copy of Topic.supd
	Supd chan<- *SessionUpdate
}

// Session update: user agent change or background session becoming normal.
// If sess is nil then user agent change, otherwise bg to fg update.
type SessionUpdate struct {
	Sess      *Session
	UserAgent string
}

type BoundedWaitGroup struct {
	Wg  sync.WaitGroup
	Sem chan struct{}
}
