package server

import (
	"sync"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// Request to hub to remove the topic
type TopicUnreg struct {
	// Original request, could be nil.
	pkt *ClientComMessage
	// Session making the request, could be nil.
	sess *Session
	// Routable name of the topic to drop. Duplicated here because pkt could be nil.
	RcptTo string
	// UID of the user being deleted. Duplicated here because pkt could be nil.
	ForUser types.Uid
	// Unregister then delete the topic.
	Del bool
	// Channel for reporting operation completion when deleting topics for a user.
	Done chan<- bool
}

type UserStatusReq struct {
	// UID of the user being affected.
	ForUser types.Uid
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
	unreg chan *TopicUnreg

	// Channel for suspending/resuming users, buffered 128.
	userStatus chan *UserStatusReq

	// Cluster request to rehash topics, unbuffered
	rehash chan bool

	// Request to shutdown, unbuffered
	shutdown chan chan<- bool
}
