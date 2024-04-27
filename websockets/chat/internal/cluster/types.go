package cluster

import (
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type ClusterArgs struct {
	Cfg    *config.ClusterConfig
	Logger *logger.Logger
	Stats  *stats.Stats
}

// ProxyReqType is the type of proxy requests.
type ProxyReqType int

// Cluster is the representation of the cluster.
type Cluster struct {
	Logger *logger.Logger

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
	// ring *rh.Ring

	// Failover parameters. Could be nil if failover is not enabled
	fo *ClusterFailover

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
	// Name of the node sending this request
	node string
	// Ring hash signature of the node sending this request Signature must match
	// the signature of the receiver, otherwise the Cluster is desynchronized.
	signature string
	// Fingerprint of the node sending this request.
	// Fingerprint changes when the node is restarted.
	fingerprint int64

	// Type of request.
	feqType ProxyReqType

	// Client message. Set for C2S requests.
	// cliMsg *ClientComMessage
	// Message to be routed. Set for intra-cluster route requests.
	// srvMsg *ServerComMessage

	// Expanded (routable) topic name
	rcptTo string
	// Originating session
	sess *ClusterSess
	// True when the topic proxy is gone.
	gone bool
}

// ClusterSess is a basic info on a remote session where the message was created.
type ClusterSess struct {
	// IP address of the client. For long polling this is the IP of the last poll
	remoteAddr string

	// User agent, a string provided by an authenticated client in {login} packet
	userAgent string

	// ID of the current user or 0
	uid types.Uid

	// User's authentication level
	// authLvl AuthLevel

	// Protocol version of the client: ((major & 0xff) << 8) | (minor & 0xff)
	ver int

	// Human language of the client
	lang string
	// Country of the client
	countryCode string

	// Device ID
	deviceID string

	// Device platform: "web", "ios", "android"
	platform string

	// Session ID
	sid string

	// Background session
	background bool
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

// ClusterHealth is content of a leader's health check of a follower node.
type ClusterHealth struct {
	// Name of the leader node
	leader string
	// Election term
	term int
	// Ring hash signature that represents the cluster
	signature string
	// Names of nodes currently active in the cluster
	nodes []string
}

// ClusterVote is a vote request and a response in leader election.
type ClusterVote struct {
	req  *ClusterVoteRequest
	resp chan ClusterVoteResponse
}

// ClusterVoteRequest is a request from a leader candidate to a node to vote for the candidate.
type ClusterVoteRequest struct {
	// Candidate node which issued this request
	node string
	// Election term
	term int
}

// ClusterVoteResponse is a vote from a node.
type ClusterVoteResponse struct {
	// Actual vote
	result bool
	// Node's term after the vote
	term int
}
