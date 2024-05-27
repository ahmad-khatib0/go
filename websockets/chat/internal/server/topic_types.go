package server

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

const (
	// Topic is fully initialized.
	topicStatusLoaded = 0x1
	// Topic is paused: all packets are rejected.
	topicStatusPaused = 0x2

	// Topic is in the process of being deleted. This is irrecoverable.
	topicStatusMarkedDeleted = 0x10
	// Topic is suspended: read-only mode.
	topicStatusReadOnly = 0x20
)

var (
	nilPresParams  = &presParams{}
	nilPresFilters = &presFilters{}
)

// Topic is an isolated communication channel
type Topic struct {
	// Еxpanded/unique name of the topic.
	name string
	// For single-user topics session-specific topic name, such as 'me',
	// otherwise the same as 'name'.
	xoriginal string

	// Topic category
	cat types.TopicCat

	// Name of the master node for this topic if isProxy is true.
	masterNode string

	// Time when the topic was first created.
	created time.Time
	// Time when the topic was last updated.
	updated time.Time
	// Time of the last outgoing message.
	touched time.Time

	// Server-side ID of the last data message
	lastID int
	// ID of the deletion operation. Not an ID of the message.
	delID int

	// Last published userAgent ('me' topic only)
	userAgent string

	// User ID of the topic owner/creator. Could be zero.
	owner types.Uid

	// Default access mode
	accessAuth types.AccessMode
	accessAnon types.AccessMode

	// Topic discovery tags
	tags []string

	// Topic's public data
	public any
	// Topic's trusted data
	trusted any

	// Topic's per-subscriber data
	perUser map[types.Uid]perUserData
	// Union of permissions across all users (used by proxy sessions with uid = 0).
	// These are used by master topics only (in the proxy-master topic context)
	// as a coarse-grained attempt to perform acs checks since proxy sessions "impersonate"
	// multiple normal sessions (uids) which may have different uids.
	modeWantUnion  types.AccessMode
	modeGivenUnion types.AccessMode

	// User's contact list (not nil for 'me' topic only).
	// The map keys are UserIds for P2P topics and grpXXX for group topics.
	perSubs map[string]perSubsData

	// Sessions attached to this topic. The UID kept here may not match Session.uid if session is
	// subscribed on behalf of another user.
	sessions map[*Session]perSessionData

	// Present video call data. Null when there's no call in progress or being established.
	// Only available for p2p topics.
	currentCall *videoCall

	// Channel for receiving client messages from sessions or other topics, buffered = 256.
	clientMsg chan *ClientComMessage
	// Channel for receiving server messages generated on the server or received from other cluster nodes, buffered = 64.
	serverMsg chan *ServerComMessage
	// Channel for receiving {get}/{set}/{del} requests, buffered = 64
	meta chan *ClientComMessage
	// Subscribe requests from sessions, buffered = 256
	reg chan *ClientComMessage
	// Unsubscribe requests from sessions, buffered = 256
	unreg chan *ClientComMessage
	// Session updates: background sessions coming online, User Agent changes. Buffered = 32
	supd chan *sessionUpdate
	// Channel to terminate topic  -- either the topic is deleted or system is being shut down. Buffered = 1.
	exit chan *shutDown
	// Channel to receive topic master responses (used only by proxy topics).
	proxy chan *ClusterResp
	// Channel to receive topic proxy service requests, e.g. sending deferred notifications.
	master chan *ClusterSessUpdate

	// Flag which tells topic lifecycle status: new, ready, paused, marked for deletion.
	status int32

	// Channel functionality is enabled for the group topic.
	isChan bool

	// If isProxy == true, the actual topic is hosted by another cluster member.
	// The topic should:
	// 1. forward all messages to master
	// 2. route replies from the master to sessions.
	// 3. disconnect sessions at master's request.
	// 4. shut down the topic at master's request.
	// 5. aggregate access permissions on behalf of attached sessions.
	isProxy bool

	// Countdown timer for destroying the topic when there are no more attached sessions to it.
	killTimer *time.Timer

	// Countdown timer for terminating iniatated (but not established) calls.
	callEstablishmentTimer *time.Timer
}

// perUserData holds topic's cache of per-subscriber data
type perUserData struct {
	// Count of subscription online and announced (presence not deferred).
	online int

	// Last t.lastId reported by user through {pres} as received or read
	recvID int
	readID int
	// ID of the latest Delete operation
	delID int

	private any

	modeWant  types.AccessMode
	modeGiven types.AccessMode

	// P2P only:
	public   any
	trusted  any
	lastSeen *time.Time
	lastUA   string

	topicName string
	deleted   bool

	// The user is a channel subscriber.
	isChan bool
}

// perSubsData holds user's (on 'me' topic) cache of subscription data
type perSubsData struct {
	// The other user's/topic's online status as seen by this user.
	online bool
	// True if we care about the updates from the other user/topic: (want&given).IsPresencer().
	// Does not affect sending notifications from this user to other users.
	enabled bool
}

// Data related to a subscription of a session to a topic.
type perSessionData struct {
	// ID of the subscribed user (asUid); not necessarily the session owner.
	// Could be zero for multiplexed sessions in cluster.
	uid types.Uid
	// This is a channel subscription
	isChanSub bool
	// IDs of subscribed users in a multiplexing session.
	muids []types.Uid
}

// Session update: user agent change or background session becoming normal.
//
// If sess is nil then user agent change, otherwise bg to fg update.
type sessionUpdate struct {
	sess      *Session
	userAgent string
}

// Topic shutdown
type shutDown struct {
	// Channel to report back completion of topic shutdown. Could be nil
	done chan<- bool
	// Topic is being deleted as opposite to total system shutdown
	reason int
}
