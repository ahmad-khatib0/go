package models

import (
	"net/rpc"

	at "github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// ProxyReqType is the type of proxy requests.
type ProxyReqType int

// Individual request types.
const (
	ProxyReqNone      ProxyReqType = iota
	ProxyReqJoin                   // {sub}.
	ProxyReqLeave                  // {leave}
	ProxyReqMeta                   // {meta set|get}
	ProxyReqBroadcast              // {pub}, {note}
	ProxyReqBgSession
	ProxyReqMeUserAgent
	ProxyReqCall // Used in video call proxy sessions for routing call events.
)

// ClusterSess is a basic info on a remote session where the message was created.
type ClusterSess struct {
	// IP address of the client. For long polling this is the IP of the last poll
	RemoteAddr string

	// User agent, a string provived by an authenticated client in {login} packet
	UserAgent string

	// ID of the current user or 0
	Uid types.Uid

	// User's authentication level
	AuthLvl at.Level

	// Protocol version of the client: ((major & 0xff) << 8) | (minor & 0xff)
	Ver int

	// Human language of the client
	Lang string
	// Country of the client
	CountryCode string

	// Device ID
	DeviceID string

	// Device platform: "web", "ios", "android"
	Platform string

	// Session ID
	Sid string

	// Background session
	Background bool
}

// ClusterSessUpdate represents a request to update a session.
// User Agent change or background session comes to foreground.
type ClusterSessUpdate struct {
	// User this session represents.
	Uid types.Uid
	// Session id.
	Sid string
	// Session user agent.
	UserAgent string
}

// ClusterReq is either a Proxy to Master or Topic Proxy to Topic Master or intra-cluster routing request message.
type ClusterReq struct {
	// Name of the node sending this request
	Node string

	// Ring hash signature of the node sending this request
	// Signature must match the signature of the receiver, otherwise the
	// Cluster is desynchronized.
	Signature string

	// Fingerprint of the node sending this request.
	// Fingerprint changes when the node is restarted.
	Fingerprint int64

	// Type of request.
	ReqType ProxyReqType

	// Client message. Set for C2S requests.
	CliMsg *ClientComMessage
	// Message to be routed. Set for intra-cluster route requests.
	SrvMsg *ServerComMessage

	// Expanded (routable) topic name
	RcptTo string
	// Originating session
	Sess *ClusterSess
	// True when the topic proxy is gone.
	Gone bool
}

// ClusterRoute is intra-cluster routing request message.
type ClusterRoute struct {
	// Name of the node sending this request
	Node string

	// Ring hash signature of the node sending this request
	// Signature must match the signature of the receiver, otherwise the
	// Cluster is desynchronized.
	Signature string

	// Fingerprint of the node sending this request.
	// Fingerprint changes when the node is restarted.
	Fingerprint int64

	// Message to be routed. Set for intra-cluster route requests.
	SrvMsg *ServerComMessage

	// Originating session
	Sess *ClusterSess
}

// ClusterResp is a Master to Proxy response message.
type ClusterResp struct {
	// Server message with the response.
	SrvMsg *ServerComMessage
	// Originating session ID to forward response to, if any.
	OrigSid string
	// Expanded (routable) topic name
	RcptTo string

	// Parameters sent back by the topic master in response a topic proxy request.

	// Original request type.
	OrigReqType ProxyReqType
}

// ClusterPing is used to detect node restarts.
type ClusterPing struct {
	// Name of the node sending this request.
	Node string

	// Fingerprint of the node sending this request.
	// Fingerprint changes when the node restarts.
	Fingerprint int64
}

type Cluster interface {
	ProxyEventQueue() GoRoutinePool
}

type ClusterNode interface {
	Endpoint() *rpc.Client
	MasterToProxyAsync(msg *ClusterResp) error
	Name() string
	Fingerprint() int64
}

type ClusterFailover interface{}

type ClusterHealth interface{}

type ClusterVote interface{}

type ClusterVoteRequest interface{}

type ClusterVoteResponse interface{}
