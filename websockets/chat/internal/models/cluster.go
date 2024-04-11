package models

type Cluster interface{}

// // ProxyReqType is the type of proxy requests.
// type ProxyReqType int

// // Cluster is the representation of the cluster.
// type Cluster struct {
// 	// Cluster nodes with RPC endpoints (excluding current node).
// 	Nodes map[string]*ClusterNode

// 	// Name of the local node
// 	ThisNodeName string

// 	// Fingerprint of the local node
// 	Fingerprint int64

// 	// Resolved address to listed on
// 	ListenOn string

// 	// Socket for inbound connections
// 	Inbound *net.TCPListener
// 	// Ring hash for mapping topic names to nodes
// 	// ring *rh.Ring

// 	// Failover parameters. Could be nil if failover is not enabled
// 	Fo *ClusterFailover

// 	// Thread pool to use for running proxy session (write) event processing logic.
// 	// The number of proxy sessions grows as O(number of topics x number of cluster nodes).
// 	// In large Tinode deployments (10s of thousands of topics, tens of nodes),
// 	// running a separate event processing goroutine for each proxy session
// 	// leads to a rather large memory usage and excessive scheduling overhead.
// 	// proxyEventQueue *concurrency.GoRoutinePool
// }

// // ClusterNode is a client's connection to another node.
// type ClusterNode struct {
// 	Lock sync.Mutex
// 	// RPC endpoint
// 	Endpoint *rpc.Client
// 	// True if the endpoint is believed to be connected
// 	Connected bool
// 	// True if a go routine is trying to reconnect the node
// 	Reconnecting bool
// 	// TCP address in the form host:port
// 	Address string
// 	// Name of the node
// 	Name string
// 	// Fingerprint of the node: unique value which changes when the node restarts.
// 	Fingerprint int64
// 	// A number of times this node has failed in a row
// 	FailCount int
// 	// Channel for shutting down the runner; buffered, 1.

// 	Done chan bool
// 	// IDs of multiplexing sessions belonging to this node.
// 	Msess map[string]struct{}
// 	// Default channel for receiving responses to RPC calls issued by this node.
// 	// Buffered, clusterRpcCompletionBuffer * number_of_nodes.
// 	RpcDone chan *rpc.Call
// 	// Channel for sending proxy to master requests; buffered, clusterProxyToMasterBuffer.
// 	p2mSender chan *ClusterReq
// }

// // ClusterReq is either a Proxy to Master or Topic Proxy to Topic Master or intra-cluster routing request message.
// type ClusterReq struct {
// 	// Name of the node sending this request
// 	Node string
// 	// Ring hash signature of the node sending this request Signature must match
// 	// the signature of the receiver, otherwise the Cluster is desynchronized.
// 	Signature string
// 	// Fingerprint of the node sending this request.
// 	// Fingerprint changes when the node is restarted.
// 	Fingerprint int64

// 	// Type of request.
// 	ReqType ProxyReqType

// 	// Client message. Set for C2S requests.
// 	CliMsg *ClientComMessage
// 	// Message to be routed. Set for intra-cluster route requests.
// 	SrvMsg *ServerComMessage

// 	// Expanded (routable) topic name
// 	RcptTo string
// 	// Originating session
// 	Sess *ClusterSess
// 	// True when the topic proxy is gone.
// 	Gone bool
// }

// // ClusterSess is a basic info on a remote session where the message was created.
// type ClusterSess struct {
// 	// IP address of the client. For long polling this is the IP of the last poll
// 	RemoteAddr string

// 	// User agent, a string provived by an authenticated client in {login} packet
// 	UserAgent string

// 	// ID of the current user or 0
// 	Uid types.Uid

// 	// User's authentication level
// 	AuthLvl AuthLevel

// 	// Protocol version of the client: ((major & 0xff) << 8) | (minor & 0xff)
// 	Ver int

// 	// Human language of the client
// 	Lang string
// 	// Country of the client
// 	CountryCode string

// 	// Device ID
// 	DeviceID string

// 	// Device platform: "web", "ios", "android"
// 	Platform string

// 	// Session ID
// 	Sid string

// 	// Background session
// 	Background bool
// }

// /*
// 																+----------------+
// 																| CLUSTER LEADER |
// 																+----------------+
// 	*****************************************************************************
// 	*****************************************************************************
// 	*****************************************************************************
// */

// // Failover config.
// type ClusterFailover struct {
// 	// Current leader
// 	Leader string
// 	// Current election term
// 	Term int
// 	// Hearbeat interval
// 	HeartBeat time.Duration
// 	// Vote timeout: the number of missed heartbeats before a new election is initiated.
// 	VoteTimeout int

// 	// The list of nodes the leader considers active
// 	ActiveNodes     []string
// 	ActiveNodesLock sync.RWMutex
// 	// The number of heartbeats a node can fail before being declared dead
// 	NodeFailCountLimit int

// 	// Channel for processing leader health checks.
// 	HealthCheck chan *ClusterHealth
// 	// Channel for processing election votes.
// 	ElectionVote chan *ClusterVote
// 	// Channel for stopping the failover runner.
// 	Done chan bool
// }

// // ClusterHealth is content of a leader's health check of a follower node.
// type ClusterHealth struct {
// 	// Name of the leader node
// 	Leader string
// 	// Election term
// 	Term int
// 	// Ring hash signature that represents the cluster
// 	Signature string
// 	// Names of nodes currently active in the cluster
// 	Nodes []string
// }

// // ClusterVote is a vote request and a response in leader election.
// type ClusterVote struct {
// 	req  *ClusterVoteRequest
// 	resp chan ClusterVoteResponse
// }

// // ClusterVoteRequest is a request from a leader candidate to a node to vote for the candidate.
// type ClusterVoteRequest struct {
// 	// Candidate node which issued this request
// 	Node string
// 	// Election term
// 	Term int
// }

// // ClusterVoteResponse is a vote from a node.
// type ClusterVoteResponse struct {
// 	// Actual vote
// 	Result bool
// 	// Node's term after the vote
// 	Term int
// }
