package server

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// Video call constants.
const (
	// Events for call between users A and B.
	//
	// B has received the call but hasn't picked it up yet.
	constCallEventRinging = "ringing"
	// B has accepted the call.
	constCallEventAccept = "accept"
	// WebRTC SDP & ICE data exchange events.
	constCallEventOffer        = "offer"
	constCallEventAnswer       = "answer"
	constCallEventIceCandidate = "ice-candidate"
	// Call finished by either side or server.
	constCallEventHangUp = "hang-up"

	// Message headers representing call states.
	// Call is established.
	constCallMsgAccepted = "accepted"
	// Previously establied call has successfully finished.
	constCallMsgFinished = "finished"
	// Call is dropped (e.g. because of an error).
	constCallMsgDisconnected = "disconnected"
	// Call is missed (the callee didn't pick up the phone).
	constCallMsgMissed = "missed"
	// Call is declined (the callee hung up before picking up).
	constCallMsgDeclined = "declined"
)

// callPartyData describes a video call participant.
type callPartyData struct {
	// ID of the call participant (asUid); not necessarily the session owner.
	uid types.Uid
	// True if this session/user initiated the call.
	isOriginator bool
	// Call party session.
	sess *Session
}

type videoCall struct {
	// Call participants (session sid -> callPartyData).
	parties map[string]callPartyData
	// Call message seq ID.
	seq int
	// Call message content.
	content any
	// Call message content mime type.
	contentMime any
	// Time when the call was accepted.
	acceptedAt time.Time
}
