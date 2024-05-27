package server

import (
	"net"
	"net/rpc"
	"sync"
	"time"

	at "github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/ringhash"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

const (
	// Network connection timeout.
	clusterNetworkTimeout = 3 * time.Second
	// Default timeout before attempting to reconnect to a node.
	clusterDefaultReconnectTime = 200 * time.Millisecond
	// Number of replicas in ringhash.
	clusterHashReplicas = 20
	// Buffer size for sending requests from proxy to master.
	clusterProxyToMasterBuffer = 64
	// Buffer size for receiving responses from other nodes, per node.
	clusterRpcCompletionBuffer = 64
)

type ClusterArgs struct {
	Cfg    *config.ClusterConfig
	Logger *logger.Logger
	Stats  *stats.Stats
}

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

// ProxyReqType is the type of proxy requests.
type ProxyReqType int

// Cluster is the representation of the cluster.
type Cluster struct {
	// Cluster nodes with RPC endpoints (excluding current node).
	nodes map[string]*ClusterNode

	// Name of the local node
	thisNodeName string

	// Fingerprint of the local node
	fingerprint int64

	// Resolved address to listed on
	listenOn string

	// Socket for inbound connections
	inbound *net.TCPListener
	// Ring hash for mapping topic names to nodes
	ring *ringhash.Ring

	// Failover parameters. Could be nil if failover is not enabled
	fo *clusterFailover

	// Thread pool to use for running proxy session (write) event processing logic.
	// The number of proxy sessions grows as O(number of topics x number of cluster nodes).
	// In large Tinode deployments (10s of thousands of topics, tens of nodes),
	// running a separate event processing goroutine for each proxy session
	// leads to a rather large memory usage and excessive scheduling overhead.
	proxyEventQueue models.GoRoutinePool
}

// ClusterNode is a client's connection to another node.
type ClusterNode struct {
	lock sync.Mutex
	// RPC endpoint
	endpoint *rpc.Client
	// True if the endpoint is believed to be connected
	connected bool
	// True if a go routine is trying to reconnect the node
	reconnecting bool
	// TCP address in the form host:port
	address string
	// Name of the node
	name string
	// Fingerprint of the node: unique value which changes when the node restarts.
	fingerprint int64
	// A number of times this node has failed in a row
	failCount int
	// Channel for shutting down the runner; buffered, 1.

	done chan bool
	// IDs of multiplexing sessions belonging to this node.
	msess map[string]struct{}
	// Default channel for receiving responses to RPC calls issued by this node.
	// Buffered, clusterRpcCompletionBuffer * number_of_nodes.
	rpcDone chan *rpc.Call
	// Channel for sending proxy to master requests; buffered, clusterProxyToMasterBuffer.
	p2mSender chan *ClusterReq
}

// ClusterReq is either a Proxy to Master or Topic Proxy to Topic Master or intra-cluster routing request message.
type ClusterReq struct {
	// Name of the Node sending this request
	Node string
	// Ring hash Signature of the node sending this request Signature must match
	// the Signature of the receiver, otherwise the Cluster is desynchronized.
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
	// True when the topic proxy is Gone.
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

// ClusterPing is used to detect node restarts.
type ClusterPing struct {
	// Name of the node sending this request.
	Node string

	// Fingerprint of the node sending this request.
	// Fingerprint changes when the node restarts.
	Fingerprint int64
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

// ClusterSess is a basic info on a remote session where the message was created.
type ClusterSess struct {
	// IP address of the client. For long polling this is the IP of the last poll
	RemoteAddr string

	// User agent, a string provided by an authenticated client in {login} packet
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

	// Device Platform: "web", "ios", "android"
	Platform string

	// Session ID
	Sid string

	// Background session
	Background bool
}

/*
								+----------------+
								| CLUSTER LEADER |
								+----------------+
	*****************************************************************************
	*****************************************************************************
	*****************************************************************************
*/

// Failover config.
type ClusterFailover struct {
	// Current leader
	leader string
	// Current election term
	term int
	// Heartbeat interval
	heartBeat time.Duration
	// Vote timeout: the number of missed heartbeats before a new election is initiated.
	voteTimeout int

	// The list of nodes the leader considers active
	activeNodes     []string
	activeNodesLock sync.RWMutex
	// The number of heartbeats a node can fail before being declared dead
	nodeFailCountLimit int

	// Channel for processing leader health checks.
	healthCheck chan *ClusterHealth
	// Channel for processing election votes.
	electionVote chan *ClusterVote
	// Channel for stopping the failover runner.
	done chan bool
}

// Failover config.
type clusterFailover struct {
	// Current leader
	leader string
	// Current election term
	term int
	// Hearbeat interval
	heartBeat time.Duration
	// Vote timeout: the number of missed heartbeats before a new election is initiated.
	voteTimeout int

	// The list of nodes the leader considers active
	activeNodes     []string
	activeNodesLock sync.RWMutex
	// The number of heartbeats a node can fail before being declared dead
	nodeFailCountLimit int

	// Channel for processing leader health checks.
	healthCheck chan *ClusterHealth
	// Channel for processing election votes.
	electionVote chan *ClusterVote
	// Channel for stopping the failover runner.
	done chan bool
}

type clusterFailoverConfig struct {
	// Failover is enabled
	Enabled bool `json:"enabled"`
	// Time in milliseconds between heartbeats
	Heartbeat int `json:"heartbeat"`
	// Number of failed heartbeats before a leader election is initiated.
	VoteAfter int `json:"vote_after"`
	// Number of failures before a node is considered dead
	NodeFailAfter int `json:"node_fail_after"`
}

// ClusterHealth is content of a leader's health check of a follower node.
type ClusterHealth struct {
	// Name of the leader node
	Leader string
	// Election term
	Term int
	// Ring hash signature that represents the cluster
	Signature string
	// Names of nodes currently active in the cluster
	Nodes []string
}

// ClusterVoteRequest is a request from a leader candidate to a node to vote for the candidate.
type ClusterVoteRequest struct {
	// Candidate node which issued this request
	Node string
	// Election term
	Term int
}

// ClusterVoteResponse is a vote from a node.
type ClusterVoteResponse struct {
	// Actual vote
	Result bool
	// Node's term after the vote
	Term int
}

// ClusterVote is a vote request and a response in leader election.
type ClusterVote struct {
	req  *ClusterVoteRequest
	resp chan ClusterVoteResponse
}
