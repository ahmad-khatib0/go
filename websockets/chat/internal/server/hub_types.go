package server

import (
	"sync"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// RequestLatencyDistribution is an array of request latency distribution bounds (in milliseconds).
// "var" because Go does not support array constants.
var requestLatencyDistribution = []float64{
	1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30,
	40, 50, 65, 80, 100, 130, 160, 200, 250, 300,
	400, 500, 650, 800, 1000, 2000, 5000, 10000,
	20000, 50000, 100000,
}

// OutgoingMessageSizeDistribution is an array of outgoing message size distribution bounds (in bytes).
var outgoingMessageSizeDistribution = []float64{
	1, 2, 4, 8, 16, 32, 64, 128, 256, 512,
	1024, 2048, 4096, 16384, 65536, 262144,
	1048576, 4194304, 16777216, 67108864,
	268435456, 1073741824, 4294967296,
}

// Request to hub to remove the topic
type topicUnreg struct {
	// Original request, could be nil.
	pkt *ClientComMessage
	// Session making the request, could be nil.
	sess *Session
	// Routable name of the topic to drop. Duplicated here because pkt could be nil.
	rcptTo string
	// UID of the user being deleted. Duplicated here because pkt could be nil.
	forUser types.Uid
	// Unregister then delete the topic.
	del bool
	// Channel for reporting operation completion when deleting topics for a user.
	done chan<- bool
}

type userStatusReq struct {
	// UID of the user being affected.
	forUser types.Uid
	// New topic state value. Only types.StateSuspended is supported at this time.
	state types.ObjState
}

// Hub is the core structure which holds topics.
type Hub struct {

	// topics must be indexed by name
	topics *sync.Map

	// Current number of loaded topics
	numTopics int

	// Channel for routing client-side messages, buffered at 4096
	routeCli chan *ClientComMessage

	// Process get.info requests for topic not subscribed to, buffered 128.
	meta chan *ClientComMessage

	// Channel for routing server-generated messages, buffered at 4096
	routeSrv chan *ServerComMessage

	// subscribe session to topic, possibly creating a new topic, buffered at 256
	join chan *ClientComMessage

	// Remove topic from hub, possibly deleting it afterwards, buffered at 32
	unreg chan *topicUnreg

	// Channel for suspending/resuming users, buffered 128.
	userStatus chan *userStatusReq

	// Cluster request to rehash topics, unbuffered
	rehash chan bool

	// Request to shutdown, unbuffered
	shutdown chan chan<- bool
}
