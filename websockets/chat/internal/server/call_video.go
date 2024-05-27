package server

import (
	"strconv"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// callPartySession returns a session to be stored in the call party data.
func callPartySession(sess *Session) *Session {
	if sess.isProxy() {
		// We are on the topic host node. Make a copy of the ephemeral proxy session.
		callSess := &Session{
			proto: PROXY,
			// Multiplexing session which actually handles the communication.
			multi: sess.multi,
			// Local parameters specific to this session.
			sid:         sess.sid,
			userAgent:   sess.userAgent,
			remoteAddr:  sess.remoteAddr,
			lang:        sess.lang,
			countryCode: sess.countryCode,
			proxyReq:    ProxyReqCall,
			background:  sess.background,
			uid:         sess.uid,
		}
		return callSess
	}
	return sess
}

func (s *Server) InitVideoCalls(c *config.WebrtcConfig) error {
	if !c.Enabled || len(c.IceServers) == 0 {
		return nil
	}

	if len(c.IceServers) > 0 {
		globals.iceServers = c.IceServers
	}

	globals.callEstablishmentTimeout = c.CallEstablishmentTimeout
	if globals.callEstablishmentTimeout <= 0 {
		globals.callEstablishmentTimeout = constants.DefaultCallEstablishmentTimeout
	}

	return nil
}

// Add webRTC-related headers to message Head. The original Head may
// already contain some entries, like 'sender', preserve them.
func (call *videoCall) messageHead(head map[string]any, newState string, duration int) map[string]any {
	if head == nil {
		head = map[string]any{}
	}

	head["replace"] = ":" + strconv.Itoa(call.seq)
	head["webrtc"] = newState

	if duration > 0 {
		head["webrtc-duration"] = duration
	} else {
		delete(head, "webrtc-duration")
	}
	if call.contentMime != nil {
		head["mime"] = call.contentMime
	}

	return head
}

// Generates server info message template for the video call event.
func (call *videoCall) infoMessage(event string) *ServerComMessage {
	return &ServerComMessage{
		Info: &MsgServerInfo{
			What:  "call",
			Event: event,
			SeqId: call.seq,
		},
	}
}

// Returns Uid and session of the present video call originator
// if a call is being established or in progress.
func (t *Topic) getCallOriginator() (types.Uid, *Session) {
	if t.currentCall == nil {
		return types.ZeroUid, nil
	}
	for _, p := range t.currentCall.parties {
		if p.isOriginator {
			return p.uid, p.sess
		}
	}
	return types.ZeroUid, nil
}

// Handles video call invite (initiation)
// (in response to msg = {pub head=[mime: application/x-tiniode-webrtc]}).
func (t *Topic) handleCallInvite(msg *ClientComMessage, asUid types.Uid) {
	// Call being establshed.
	t.currentCall = &videoCall{
		parties:     make(map[string]callPartyData),
		seq:         t.lastID,
		content:     msg.Pub.Content,
		contentMime: msg.Pub.Head["mime"],
	}

	t.currentCall.parties[msg.sess.sid] = callPartyData{
		uid:          asUid,
		isOriginator: true,
		sess:         callPartySession(msg.sess),
	}

	// Wait for constCallEstablishmentTimeout for the other side to accept the call.
	t.callEstablishmentTimer.Reset(time.Duration(globals.callEstablishmentTimeout) * time.Second)
}

// Handles events on existing video call (acceptance, termination, metadata exchange).
// (in response to msg = {note what=call}).
func (t *Topic) handleCallEvent(msg *ClientComMessage) {
	if t.currentCall == nil {
		// Must initiate call first.
		globals.l.Sugar().Warnf("topic[%s]: No call in progress", t.name)
		return
	}

	if t.isInactive() {
		// Topic is paused or being deleted.
		return
	}

	call := msg.Note
	if t.currentCall.seq != call.SeqId {
		// Call not found.
		globals.l.Sugar().Infof("topic[%s]: invalid seq id - current call (%d) vs received (%d)", t.name, t.currentCall.seq, call.SeqId)
		return
	}

	asUid := types.ParseUserId(msg.AsUser)
	if _, userFound := t.perUser[asUid]; !userFound {
		// User not found in topic.
		globals.l.Sugar().Infof("topic[%s]: could not find user %s", t.name, asUid.UserId())
		return
	}

	switch call.Event {
	case constCallEventRinging, constCallEventAccept:
		// Invariants:
		// 1. Call has been initiated but not been established yet.
		if len(t.currentCall.parties) != 1 {
			return
		}

		originatorUid, originator := t.getCallOriginator()
		if originator == nil {
			// No originator session: terminating.
			t.terminateCallInProgress(false)
			return
		}

		// 2. These events may only arrive from the callee.
		if originator.sid == msg.sess.sid || originatorUid == asUid {
			return
		}

		// Prepare a {info} message to forward to the call originator.
		forwardMsg := t.currentCall.infoMessage(call.Event)
		forwardMsg.Info.From = msg.AsUser
		forwardMsg.Info.Topic = t.original(originatorUid)

		if call.Event == constCallEventAccept {
			// The call has been accepted: Send a replacement {data} message to the topic.
			msgCopy := *msg
			msgCopy.AsUser = originatorUid.UserId()
			replaceWith := constCallMsgAccepted

			var origHead map[string]any
			if msgCopy.Pub != nil {
				origHead = msgCopy.Pub.Head
			}

			// else fetch the original message from store and use its head.
			head := t.currentCall.messageHead(origHead, replaceWith, 0)
			if err := t.saveAndBroadcastMessage(
				&msgCopy,
				originatorUid,
				false,
				nil,
				head,
				t.currentCall.content,
			); err != nil {
				return
			}

			// Add callee data to t.currentCall.
			t.currentCall.parties[msg.sess.sid] = callPartyData{
				uid:          asUid,
				isOriginator: false,
				sess:         callPartySession(msg.sess),
			}
			t.currentCall.acceptedAt = time.Now()

			// Notify other clients that the call has been accepted.
			t.infoCallSubsOffline(msg.AsUser, asUid, call.Event, t.currentCall.seq, call.Payload, msg.sess.sid, false)
			t.callEstablishmentTimer.Stop()
		}
		originator.queueOut(forwardMsg)

	case constCallEventOffer, constCallEventAnswer, constCallEventIceCandidate:
		// Invariants:
		// 1. Call has been estabslied (2 participants).
		if len(t.currentCall.parties) != 2 {
			globals.l.Sugar().Warnf("topic[%s]: call participants expected 2 vs found %d", t.name, len(t.currentCall.parties))
			return
		}

		// 2. Event is coming from a call participant session.
		if _, ok := t.currentCall.parties[msg.sess.sid]; !ok {
			globals.l.Sugar().Warnf("topic[%s]: call event from non-party session %s", t.name, msg.sess.sid)
			return
		}

		// Call metadata exchange. Either side of the call may send these events.
		// Simply forward them to the other session.
		var otherUid types.Uid
		var otherEnd *Session
		for sid, p := range t.currentCall.parties {
			if sid != msg.sess.sid {
				otherUid = p.uid
				otherEnd = p.sess
				break
			}
		}

		if otherEnd == nil {
			globals.l.Sugar().Warnf("topic[%s]: could not find call peer for session %s", t.name, msg.sess.sid)
			return
		}

		// All is good. Send {info} message to the otherEnd.
		forwardMsg := t.currentCall.infoMessage(call.Event)
		forwardMsg.Info.From = msg.AsUser
		forwardMsg.Info.Topic = t.original(otherUid)
		forwardMsg.Info.Payload = call.Payload
		otherEnd.queueOut(forwardMsg)

	case constCallEventHangUp:
		switch len(t.currentCall.parties) {
		case 2:
			// If it's a call in progress, hangup may arrive only from a call participant session.
			if _, ok := t.currentCall.parties[msg.sess.sid]; !ok {
				return
			}
		case 1:
			// Call hasn't been established yet.
			originatorUid, originator := t.getCallOriginator()
			// Hangup may come from either the originating session or any callee user session.
			if asUid == originatorUid && originator.sid != msg.sess.sid {
				return
			}

		default:
			break
		}

		t.maybeEndCallInProgress(msg.AsUser, msg, false)

	default:
		globals.l.Sugar().Warnf("topic[%s]: video call (seq %d) received unexpected call event: %s", t.name, t.currentCall.seq, call.Event)
	}
}

// Ends current call in response to a client hangup request (msg).
func (t *Topic) maybeEndCallInProgress(from string, msg *ClientComMessage, callDidTimeout bool) {
	if t.currentCall == nil {
		return
	}

	t.callEstablishmentTimer.Stop()
	originatorUid, _ := t.getCallOriginator()
	var replaceWith string
	var callDuration int64

	if from != "" && len(t.currentCall.parties) == 2 {
		// This is a call in progress.
		replaceWith = constCallMsgFinished
		callDuration = time.Since(t.currentCall.acceptedAt).Milliseconds()
	} else {
		if from != "" {
			// User originated hang-up.
			if from == originatorUid.UserId() {
				// Originator/caller requested event.
				replaceWith = constCallMsgMissed
			} else {
				// Callee requested event.
				replaceWith = constCallMsgDeclined
			}
		} else {
			// Server initiated disconnect.
			// Call hasn't been established. Just drop it.
			if callDidTimeout {
				replaceWith = constCallMsgMissed
			} else {
				replaceWith = constCallMsgDisconnected
			}
		}
	}

	// Send a message indicating the call has ended.
	msgCopy := *msg
	msgCopy.AsUser = originatorUid.UserId()
	var origHead map[string]any
	if msgCopy.Pub != nil {
		origHead = msgCopy.Pub.Head
	} // else fetch the original message from store and use its head.

	head := t.currentCall.messageHead(origHead, replaceWith, int(callDuration))
	if err := t.saveAndBroadcastMessage(&msgCopy, originatorUid, false, nil, head, t.currentCall.content); err != nil {
		globals.l.Sugar().Errorf("topic[%s]: failed to write finalizing message for call seq id %d - '%s'", t.name, t.currentCall.seq, err)
	}

	// Send {info} hangup event to the subscribed sessions.
	t.broadcastToSessions(t.currentCall.infoMessage(constCallEventHangUp))

	// Let all other sessions know the call is over.
	for tgt := range t.perUser {
		t.infoCallSubsOffline(from, tgt, constCallEventHangUp, t.currentCall.seq, nil, "", true)
	}
	t.currentCall = nil
}

// Server initiated call termination.
func (t *Topic) terminateCallInProgress(callDidTimeout bool) {
	if t.currentCall == nil {
		return
	}

	uid, sess := t.getCallOriginator()
	if sess == nil || uid.IsZero() {
		// Just drop the call.
		globals.l.Sugar().Infof("topic[%s]: video call seq %d has no originator, terminating.", t.name, t.currentCall.seq)
		t.currentCall = nil
		return
	}

	// Dummy hangup request.
	dummy := &ClientComMessage{
		Original:  t.original(uid),
		RcptTo:    uid.UserId(),
		AsUser:    uid.UserId(),
		Timestamp: types.TimeNow(),
		sess:      sess,
	}

	globals.l.Sugar().Infof("topic[%s]: terminating call seq %d, timeout: %t", t.name, t.currentCall.seq, callDidTimeout)

	t.maybeEndCallInProgress("", dummy, callDidTimeout)
}
