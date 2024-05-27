package server

import (
	"container/list"
	"sync"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
	at "github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
	"github.com/gorilla/websocket"
)

// Maximum number of queued messages before session is considered stale and dropped.
const sendQueueLimit = 128

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

type Session struct {
	// reference to the session store struct
	sessStore *SessionStore

	cluster *Cluster
	// Reference to the cluster node where the session has originated. Set only for cluster RPC sessions.
	clnode *ClusterNode
	logger *logger.Logger

	// protocol - NONE (unset), WEBSOCK, LPOLL, GRPC, PROXY, MULTIPLEX
	proto SessionProto
	// Session ID
	sid string

	// Websocket. Set only for websocket sessions.
	ws *websocket.Conn

	// Pointer to session's record in sessionStore. Set only for Long Poll sessions.
	lpTracker *list.Element

	// gRPC handle. Set only for gRPC clients.
	grpcCNode chat.Node_MessageLoopServer

	// Reference to multiplexing session. Set only for proxy sessions.
	multi        *Session
	proxiedTopic string

	// IP address of the client. For long polling this is the IP of the last poll.
	remoteAddr string

	// User agent, a string provided by an authenticated client in {login} packet.
	userAgent string

	// Protocol version of the client: ((major & 0xff) << 8) | (minor & 0xff).
	ver int

	// Device ID of the client
	deviceID string
	// Platform: web, ios, android
	platf string

	// Human language of the client
	lang string

	// Country code of the client
	countryCode string

	// ID of the current user. Could be zero if session is not authenticated
	// or for multiplexing sessions.
	uid types.Uid

	// Authentication level - NONE (unset), ANON, AUTH, ROOT.
	authLvl at.Level

	// Time when the long polling session was last refreshed
	lastTouched time.Time

	// Time when the session received any packer from client
	lastAction int64

	// Timer which triggers after some seconds to mark background session as foreground.
	bkgTimer *time.Timer

	// Number of subscribe/unsubscribe requests in flight.
	inflightReqs *boundedWaitGroup

	// Synchronizes access to session store in cluster mode:
	// subscribe/unsubscribe replies are asynchronous.
	sessionStoreLock sync.Mutex

	// Indicates that the session is terminating.
	// After this flag's been flipped to true, there must not be any more writes
	// into the session's send channel.
	// Read/written atomically.
	// 0 = false
	// 1 = true
	terminating int32

	// Background session: subscription presence notifications and online status are delayed.
	background bool

	// Outbound messages, buffered.
	// The content must be serialized in format suitable for the session.
	send chan any

	// Channel for shutting down the session, buffer 1.
	// Content in the same format as for 'send'
	stop chan any

	// detach - channel for detaching session from topic, buffered.
	// Content is topic name to detach from.
	detach chan string

	// Map of topic subscriptions, indexed by topic name.
	// Don't access directly. Use getters/setters.
	subs map[string]*Subscription

	// Mutex for subs access: both topic go routines and network go routines access subs concurrently.
	subsLock sync.RWMutex

	// Needed for long polling and grpc.
	lock sync.Mutex

	// Field used only in cluster mode by topic master node.
	proxyReq ProxyReqType
}

// SessionStore holds live sessions. Long polling sessions are stored in a linked list with
//
// most recent sessions on top. In addition all sessions are stored in a map indexed by session ID.
type SessionStore struct {
	logger *logger.Logger
	lock   sync.Mutex

	// Support for long polling sessions: a list of sessions sorted by last access time.
	// Needed for cleaning abandoned sessions.
	lru      *list.List
	lifeTime time.Duration

	// All sessions indexed by session ID
	sessCache map[string]*Session
	uGen      *types.UidGenerator
	stats     *stats.Stats
}

// Subscription is a mapper of sessions to topics.
type Subscription struct {
	// Channel to communicate with the topic, copy of Topic.clientMsg
	droadcast chan<- *ClientComMessage

	// Session sends a signal to Topic when this session is unsubscribed This is a copy of Topic.unreg
	done chan<- *ClientComMessage

	// Channel to send {meta} requests, copy of Topic.meta
	deta chan<- *ClientComMessage

	// Channel to ping topic with session's updates, copy of Topic.supd
	supd chan<- *SessionUpdate
}

// Session update: user agent change or background session becoming normal.
//
// If sess is nil then user agent change, otherwise bg to fg update.
type SessionUpdate struct {
	Sess      *Session
	UserAgent string
}
