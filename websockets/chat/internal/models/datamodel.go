package models

import (
	"encoding/json"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// ClientComMessage is a wrapper for client messages.
type ClientComMessage struct {
	Hi    *MsgClientHi    `json:"hi"`
	Acc   *MsgClientAcc   `json:"acc"`
	Login *MsgClientLogin `json:"login"`
	Sub   *MsgClientSub   `json:"sub"`
	Leave *MsgClientLeave `json:"leave"`
	Pub   *MsgClientPub   `json:"pub"`
	Get   *MsgClientGet   `json:"get"`
	Set   *MsgClientSet   `json:"set"`
	Del   *MsgClientDel   `json:"del"`
	Note  *MsgClientNote  `json:"note"`
	// Optional data.
	Extra *MsgClientExtra `json:"extra"`

	// Internal fields, routed only within the cluster.

	// Message ID de-normalized
	ID string `json:"-"`
	// Un-routable (original) topic name de-normalized from XXX.Topic.
	Original string `json:"-"`
	// Routable (expanded) topic name.
	RcptTo string `json:"-"`
	// Sender's UserId as string.
	AsUser string `json:"-"`
	// Sender's authentication level.
	AuthLvl int `json:"-"`
	// De-normalized 'what' field of meta messages (set, get, del).
	MetaWhat int `json:"-"`
	// Timestamp when this message was received by the server.
	Timestamp time.Time `json:"-"`

	// Originating session to send an acknowledgement to.
	Sess Session
	// The message is initialized (true) as opposite to being used as a wrapper for session.
	Init bool
}

/****************************************************************
 * Client to Server (C2S) messages.
 ****************************************************************/

// MsgClientHi is a handshake {hi} message
type MsgClientHi struct {
	ID        string `json:"id"`
	UserAgent string `json:"user_agent"`
	Version   int    `json:"version"`   // Protocol version, i.e. "0.13"
	DeviceID  string `json:"device_id"` // Client's unique device ID
	Lang      string `json:"lang"`      // ISO 639-1 human language of the connected device
	Platform  string `json:"platform"`  // Platform code: ios, android, web.
	// Session is initially in non-iteractive, i.e. issued by a service. Presence notifications are delayed.
	Background bool `json:"background"`
}

// MsgClientAcc is an {acc} message for creating or updating a user account.
type MsgClientAcc struct {
	// message id
	ID string
	// "newXYZ" to create a new user or UserId to update a user; default: current user.
	User string
	// Temporary authentication parameters for one-off actions, like password reset.
	TmpSchema string
	TmpSecret []byte
	// Account state: normal, suspended.
	State string
	// Authentication level of the user when UserID is set and not equal to the current user.
	// Either "", "auth" or "anon". Default: ""
	AuthLevel string
	// The initial authentication scheme the account can use
	Scheme string
	// Shared secret
	Secret []byte
	// Authenticate session with the newly created account
	Login bool
	// Indexable tags for user discovery
	Tags []string
	// User initialization data when creating a new user, otherwise ignored
	Desc *MsgSetDesc
	// Credentials to verify (email or phone or captcha)
	Cred []MsgCredClient
}

// MsgSetDesc is a C2S in set.what == "desc", acc, sub message.
type MsgSetDesc struct {
	DFM     *MsgDefaultAcsMode `json:"dfm"`     // default access mode
	Public  any                `json:"public"`  // description of the user or topic
	Trusted any                `json:"trusted"` // trusted (system-provided) user or topic data
	Private any                `json:"private"` // per-subscription private data
}

// MsgDefaultAcsMode is a C2S in set.what == "desc", acc, sub message.
type MsgDefaultAcsMode struct {
	Auth string `mapstructure:"auth"`
	Anon string `mapstructure:"anon"`
}

// MsgCredClient is an account credential such as email or phone number.
type MsgCredClient struct {
	// Credential type, i.e. `email` or `tel`.
	Method string
	// Value to verify, i.e. `user@example.com` or `+18003287448`
	Value string
	// Verification response
	Response string
	// Request parameters, such as preferences. Passed to valiator without interpretation.
	Params map[string]any
}

// MsgClientLogin is a login {login} message.
type MsgClientLogin struct {
	// message id
	ID string
	// Authentication scheme
	Scheme string
	// Shared secret
	Secret []byte
	// Credntials being verified (email or phone or captcha etc.)
	Cred []MsgCredClient
}

// MsgClientSub  is a subscription request {sub} message.
type MsgClientSub struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`

	// Mirrors {set}.
	Set *MsgSetQuery `json:"set"`
	// Mirrors {get}.
	Get *MsgGetQuery `json:"get"`

	// Intra-cluster fields.

	// True if this subscription created a new topic. In case of p2p topics, it's true if the other
	// user's subscription was created (as a part of new topic creation or just alone).
	Created bool `json:"-"`
	// True if this is a new subscription.
	NewSub bool `json:"-"`
}

// MsgSetQuery is an update to topic or user metadata: description, subscriptions, tags, credentials.
type MsgSetQuery struct {
	// Topic/user description, new object & new subscriptions only
	Desc *MsgSetDesc `json:"desc"`
	// Subscription parameters
	Sub *MsgSetSub `json:"sub"`
	// Indexable tags for user discovery
	Tags []string `json:"tags"`
	// Update to account credentials.
	Cred *MsgCredClient `json:"cred"`
}

// MsgSetSub is a payload in set.sub request to update current
// subscription or invite another user, {sub.what} == "sub".
type MsgSetSub struct {
	// User affected by this request. Default (empty): current user
	User string `json:"user"`
	// Access mode change, either Given or Want depending on context
	Mode string `json:"mode"`
}

// MsgClientLeave is an unsubscribe {leave} request message.
type MsgClientLeave struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`
	Unsub bool   `json:"unsub"`
}

// MsgClientPub is client's request to publish data to topic subscribers {pub}.
type MsgClientPub struct {
	ID      string         `json:"id"`
	Topic   string         `json:"topic"`
	NoEcho  bool           `json:"no_echo"`
	Head    map[string]any `json:"head"`
	Content any            `json:"content"`
}

// MsgClientGet is a query of topic state {get}.
type MsgClientGet struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`
	MsgGetQuery
}

// MsgClientSet is an update of topic state {set}.
type MsgClientSet struct {
	ID    string `json:"id,omitempty"`
	Topic string `json:"topic"`
	MsgSetQuery
}

// MsgClientDel delete messages or topic {del}.
type MsgClientDel struct {
	ID    string `json:"id,omitempty"`
	Topic string `json:"topic,omitempty"`
	// What to delete:
	// * "msg" to delete messages (default)
	// * "topic" to delete the topic
	// * "sub" to delete a subscription to topic.
	// * "user" to delete or disable user.
	// * "cred" to delete credential (email or phone)
	What string `json:"what"`
	// Delete messages with these IDs (either one by one or a set of ranges)
	DelSeq []MsgDelRange `json:"del_seq"`
	// User ID of the user or subscription to delete
	User string `json:"user"`
	// Credential to delete
	Cred *MsgCredClient `json:"cred"`
	// Request to hard-delete objects (i.e. delete messages for all users), if such option is available.
	Hard bool `json:"hard"`
}

// MsgClientNote is a client-generated notification for topic subscribers {note}.
type MsgClientNote struct {
	// There is no Id -- server will not akn {ping} packets, they are "fire and forget"
	Topic string `json:"topic"`
	// what is being reported: "recv" - message received, "read" - message read, "kp" - typing notification
	What string `json:"what"`
	// Server-issued message ID being reported
	SeqId int `json:"seq,omitempty"`
	// Client's count of unread messages to report back to the server. Used in push notifications on iOS.
	Unread int `json:"unread,omitempty"`
	// Call event.
	Event string `json:"event,omitempty"`
	// Arbitrary json payload (used in video calls).
	Payload json.RawMessage `json:"payload,omitempty"`
}

// MsgClientExtra is not a stand-alone message but extra data which augments the main payload.
type MsgClientExtra struct {
	// Array of out-of-band attachments which have to be exempted from GC.
	Attachments []string `json:"attachments,omitempty"`
	// Alternative user ID set by the root user (obo = On Behalf Of).
	AsUser string `json:"obo,omitempty"`
	// Altered authentication level set by the root user.
	AuthLevel string `json:"authlevel,omitempty"`
}

// MsgGetQuery is a topic metadata or data query.
type MsgGetQuery struct {
	What string
	// Parameters of "desc" request: IfModifiedSince
	Desc *MsgGetOpts `json:"desc,omitempty"`
	// Parameters of "sub" request: User, Topic, IfModifiedSince, Limit.
	Sub *MsgGetOpts `json:"sub,omitempty"`
	// Parameters of "data" request: Since, Before, Limit.
	Data *MsgGetOpts `json:"data,omitempty"`
	// Parameters of "del" request: Since, Before, Limit.
	Del *MsgGetOpts `json:"del,omitempty"`
}

// MsgGetOpts defines Get query parameters.
type MsgGetOpts struct {
	// Optional User ID to return result(s) for one user.
	User string `json:"user"`
	// Optional topic name to return result(s) for one topic.
	Topic string `json:"topic"`
	// Return results modified since this timespamp.
	IfModifiedSince *time.Time `json:"if_modified_since"`
	// Load messages/ranges with IDs equal or greater than this (inclusive or closed)
	SinceID int `json:"since_id"`
	// Load messages/ranges with IDs lower than this (exclusive or open)
	BeforeID string `json:"before_id"`
	// Limit the number of messages loaded
	Limit int `json:"limit"`
}

// MsgDelRange is either an individual ID (HiID=0) or a randge of deleted IDs,
//
// low end inclusive (closed), high-end exclusive (open): [LowId .. HiId), e.g. 1..5 -> 1, 2, 3, 4.
type MsgDelRange struct {
	LowID int `json:"low,omitempty"`
	HiID  int `json:"hi,omitempty"`
}

// MsgAccessMode is a definition of access mode.
type MsgAccessMode struct {
	// Access mode requested by the user
	Want string
	// Access mode granted to the user by the admin
	Given string
	// Cumulative access mode want & given
	Mode string
}

// *******************************************************
// *******************************************************
// *******************************************************
// *******************************************************
// *******************************************************
// *******************************************************
// *******************************************************

/****************************************************************
 * Server to client messages.
 ****************************************************************/

// ServerComMessage is a wrapper for server-side messages.
type ServerComMessage struct {
	Ctrl *MsgServerCtrl `json:"ctrl,omitempty"`
	Data *MsgServerData `json:"data,omitempty"`
	Meta *MsgServerMeta `json:"meta,omitempty"`
	Pres *MsgServerPres `json:"pres,omitempty"`
	Info *MsgServerInfo `json:"info,omitempty"`

	// Internal fields.

	// MsgServerData has no Id field, copying it here for use in {ctrl} aknowledgements
	Id string `json:"-"`
	// Routable (expanded) name of the topic.
	RcptTo string `json:"-"`
	// User ID of the sender of the original message.
	AsUser string `json:"-"`
	// Timestamp for consistency of timestamps in {ctrl} messages
	// (corresponds to originating client message receipt timestamp).
	Timestamp time.Time `json:"-"`
	// Originating session to send an aknowledgement to. Could be nil.
	// sess *Session
	// Session ID to skip when sendng packet to sessions. Used to skip sending to original session.
	// Could be either empty.
	SkipSid string `json:"-"`
	// User id affected by this message.
	uid types.Uid
}

// MsgServerCtrl is a server control message {ctrl}.
type MsgServerCtrl struct {
	Id        string    `json:"id,omitempty"`
	Topic     string    `json:"topic,omitempty"`
	Params    any       `json:"params,omitempty"`
	Code      int       `json:"code"`
	Text      string    `json:"text,omitempty"`
	Timestamp time.Time `json:"ts"`
}

// MsgServerData is a server {data} message.
type MsgServerData struct {
	Topic string `json:"topic"`
	// ID of the user who originated the message as {pub}, could be empty if sent by the system
	From      string         `json:"from,omitempty"`
	Timestamp time.Time      `json:"ts"`
	DeletedAt *time.Time     `json:"deleted,omitempty"`
	SeqId     int            `json:"seq"`
	Head      map[string]any `json:"head,omitempty"`
	Content   any            `json:"content"`
}

// MsgServerPres is presence notification {pres} (authoritative update).
type MsgServerPres struct {
	Topic     string        `json:"topic"`
	Src       string        `json:"src,omitempty"`
	What      string        `json:"what"`
	UserAgent string        `json:"ua,omitempty"`
	SeqId     int           `json:"seq,omitempty"`
	DelId     int           `json:"clear,omitempty"`
	DelSeq    []MsgDelRange `json:"delseq,omitempty"`
	AcsTarget string        `json:"tgt,omitempty"`
	AcsActor  string        `json:"act,omitempty"`
	// Acs or a delta Acs. Need to marshal it to json under a name different than 'acs'
	// to allow different handling on the client
	Acs *MsgAccessMode `json:"dacs,omitempty"`

	// UN-routable params. All marked with `json:"-"` to exclude from json marshaling.
	// They are still serialized for intra-cluster communication.

	// Flag to break the reply loop
	WantReply bool `json:"-"`

	// Additional access mode filters when sending to topic's online members. Both filter conditions must be true.
	// send only to those who have this access mode.
	FilterIn int `json:"-"`
	// skip those who have this access mode.
	FilterOut int `json:"-"`

	// When sending to 'me', skip sessions subscribed to this topic.
	SkipTopic string `json:"-"`

	// Send to sessions of a single user only.
	SingleUser string `json:"-"`

	// Exclude sessions of a single user.
	ExcludeUser string `json:"-"`
}

// MsgServerMeta is a topic metadata {meta} update.
type MsgServerMeta struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`

	Timestamp *time.Time `json:"ts,omitempty"`

	// Topic description
	Desc *MsgTopicDesc `json:"desc,omitempty"`
	// Subscriptions as an array of objects
	Sub []MsgTopicSub `json:"sub,omitempty"`
	// Delete ID and the ranges of IDs of deleted messages
	Del *MsgDelValues `json:"del,omitempty"`
	// User discovery tags
	Tags []string `json:"tags,omitempty"`
	// Account credentials, 'me' only.
	Cred []*MsgCredServer `json:"cred,omitempty"`
}

// MsgServerInfo is the server-side copy of MsgClientNote with From and optionally Src added (non-authoritative).
type MsgServerInfo struct {
	// Topic to send event to.
	Topic string `json:"topic"`
	// Topic where the event has occurred (set only when Topic='me').
	Src string `json:"src,omitempty"`
	// ID of the user who originated the message.
	From string `json:"from,omitempty"`
	// The event being reported: "rcpt" - message received, "read" - message read, "kp" - typing notification, "call" - video call.
	What string `json:"what"`
	// Server-issued message ID being reported.
	SeqId int `json:"seq,omitempty"`
	// Call event.
	Event string `json:"event,omitempty"`
	// Arbitrary json payload (used by video calls).
	Payload json.RawMessage `json:"payload,omitempty"`

	// UNroutable params. All marked with `json:"-"` to exclude from json marshaling.
	// They are still serialized for intra-cluster communication.

	// When sending to 'me', skip sessions subscribed to this topic.
	SkipTopic string `json:"-"`
}

// MsgTopicDesc is a topic description, S2C in Meta message.
type MsgTopicDesc struct {
	CreatedAt *time.Time `json:"created,omitempty"`
	UpdatedAt *time.Time `json:"updated,omitempty"`
	// Timestamp of the last message
	TouchedAt *time.Time `json:"touched,omitempty"`

	// Account state, 'me' topic only.
	State string `json:"state,omitempty"`

	// If the group topic is online.
	Online bool `json:"online,omitempty"`

	// If the topic can be accessed as a channel
	IsChan bool `json:"chan,omitempty"`

	// P2P other user's last online timestamp & user agent
	LastSeen *MsgLastSeenInfo `json:"seen,omitempty"`

	DefaultAcs *MsgDefaultAcsMode `json:"defacs,omitempty"`
	// Actual access mode
	Acs *MsgAccessMode `json:"acs,omitempty"`
	// Max message ID
	SeqId     int `json:"seq,omitempty"`
	ReadSeqId int `json:"read,omitempty"`
	RecvSeqId int `json:"recv,omitempty"`
	// Id of the last delete operation as seen by the requesting user
	DelId   int `json:"clear,omitempty"`
	Public  any `json:"public,omitempty"`
	Trusted any `json:"trusted,omitempty"`
	// Per-subscription private data
	Private any `json:"private,omitempty"`
}

// MsgTopicSub is topic subscription details, sent in Meta message.
type MsgTopicSub struct {
	// Fields common to all subscriptions

	// Timestamp when the subscription was last updated
	UpdatedAt *time.Time `json:"updated,omitempty"`
	// Timestamp when the subscription was deleted
	DeletedAt *time.Time `json:"deleted,omitempty"`

	// If the subscriber/topic is online
	Online bool `json:"online,omitempty"`

	// Access mode. Topic admins receive the full info, non-admins receive just the cumulative mode
	// Acs.Mode = want & given. The field is not a pointer because at least one value is always assigned.
	Acs MsgAccessMode `json:"acs,omitempty"`
	// ID of the message reported by the given user as read
	ReadSeqId int `json:"read,omitempty"`
	// ID of the message reported by the given user as received
	RecvSeqId int `json:"recv,omitempty"`
	// Topic's public data
	Public any `json:"public,omitempty"`
	// Topic's trusted public data
	Trusted any `json:"trusted,omitempty"`
	// User's own private data per topic
	Private any `json:"private,omitempty"`

	// Response to non-'me' topic

	// Uid of the subscribed user
	User string `json:"user,omitempty"`

	// The following sections makes sense only in context of getting
	// user's own subscriptions ('me' topic response)

	// Topic name of this subscription
	Topic string `json:"topic,omitempty"`
	// Timestamp of the last message in the topic.
	TouchedAt *time.Time `json:"touched,omitempty"`
	// ID of the last {data} message in a topic
	SeqId int `json:"seq,omitempty"`
	// Id of the latest Delete operation
	DelId int `json:"clear,omitempty"`

	// P2P topics in 'me' {get subs} response:

	// Other user's last online timestamp & user agent
	LastSeen *MsgLastSeenInfo `json:"seen,omitempty"`
}

// MsgDelValues describes request to delete messages.
type MsgDelValues struct {
	DelId  int           `json:"clear,omitempty"`
	DelSeq []MsgDelRange `json:"delseq,omitempty"`
}

// MsgLastSeenInfo contains info on user's appearance online - when & user agent.
type MsgLastSeenInfo struct {
	// Timestamp of user's last appearance online.
	When *time.Time `json:"when,omitempty"`
	// User agent of the device when the user was last online.
	UserAgent string `json:"ua,omitempty"`
}

// MsgCredServer is an account credential such as email or phone number.
type MsgCredServer struct {
	// Credential type, i.e. `email` or `tel`.
	Method string `json:"meth,omitempty"`
	// Credential value, i.e. `user@example.com` or `+18003287448`
	Value string `json:"val,omitempty"`
	// Indicates that the credential is validated.
	Done bool `json:"done,omitempty"`
}
