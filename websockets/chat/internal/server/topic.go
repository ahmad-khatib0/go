package server

import (
	"errors"
	"sync/atomic"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"go.uber.org/zap"
)

// Saves a new message (defined by head, content and attachments) in the topic
//
// in response to a client request (msg, asUid) and broadcasts it to the attached sessions.
func (t *Topic) saveAndBroadcastMessage(
	msg *ClientComMessage,
	asUid types.Uid,
	noEcho bool,
	attachments []string,
	head map[string]any,
	content any,
) error {
	pud, userFound := t.perUser[asUid]

	// Anyone is allowed to post to 'sys' topic.
	if t.cat != types.TopicCatSys {
		// If it's not 'sys' check write permission.
		if !(pud.modeWant & pud.modeGiven).IsWriter() {
			msg.sess.queueOut(ErrPermissionDenied(msg.Id, t.original(asUid), msg.Timestamp))
			return types.ErrPermissionDenied
		}
	}

	if msg.sess != nil && msg.sess.uid != asUid {
		// The "sender" header contains ID of the user who sent the message on behalf of asUid.
		if head == nil {
			head = map[string]any{}
		}
		head["sender"] = msg.sess.uid.UserId()
	} else if head != nil {
		// Make sure the received Head does not include a fake "sender" header.
		delete(head, "sender")
	}

	markedReadBySender := false
	if err, unreadUpdated := globals.store.MsgSave(
		&types.Message{
			ObjHeader: types.ObjHeader{CreatedAt: msg.Timestamp},
			SeqId:     t.lastID + 1,
			Topic:     t.name,
			From:      asUid.String(),
			Head:      head,
			Content:   content,
		},
		attachments,
		(pud.modeGiven & pud.modeWant).IsReader(),
	); err != nil {

		globals.l.Sugar().Warnf("topic[%s]: failed to save message: %v", t.name, err)
		msg.sess.queueOut(ErrUnknown(msg.Id, t.original(asUid), msg.Timestamp))

		return err
	} else {
		markedReadBySender = unreadUpdated
	}

	t.lastID++
	t.touched = msg.Timestamp

	if userFound {
		pud.readID = t.lastID
		pud.recvID = t.lastID
		t.perUser[asUid] = pud
	}

	if msg.Id != "" && msg.sess != nil {
		reply := NoErrAccepted(msg.Id, t.original(asUid), msg.Timestamp)
		reply.Ctrl.Params = map[string]any{"seq": t.lastID}
		msg.sess.queueOut(reply)
	}

	data := &ServerComMessage{
		Data: &MsgServerData{
			Topic:     msg.Original,
			From:      msg.AsUser,
			Timestamp: msg.Timestamp,
			SeqId:     t.lastID,
			Head:      head,
			Content:   content,
		},
		// Internal-only values.
		Id:        msg.Id,
		RcptTo:    msg.RcptTo,
		AsUser:    msg.AsUser,
		Timestamp: msg.Timestamp,
		sess:      msg.sess,
	}

	if noEcho {
		data.SkipSid = msg.sess.sid
	}

	// Message sent: notify offline 'R' subscrbers on 'me'.
	t.presSubsOffline(
		"msg",
		&presParams{seqID: t.lastID, actor: msg.AsUser},
		&presFilters{filterIn: types.ModeRead}, nilPresFilters,
		"",
		true,
	)

	// Tell the plugins that a message was accepted for delivery
	pluginMessage(data.Data, plgActCreate)

	t.broadcastToSessions(data)

	// sendPush will update unread message count and send push notification.
	if pushRcpt := t.pushForData(asUid, data.Data, markedReadBySender); pushRcpt != nil {
		sendPush(pushRcpt)
	}
	return nil
}

// broadcastToSessions writes message to attached sessions.
func (t *Topic) broadcastToSessions(msg *ServerComMessage) {
	// List of sessions to be dropped.
	var dropSessions []*Session
	// Broadcast the message. Only {data}, {pres}, {info} are broadcastable.
	// {meta} and {ctrl} are sent to the session only
	for sess, pssd := range t.sessions {
		// Send all messages to multiplexing session.
		if !sess.isMultiplex() {
			if sess.sid == msg.SkipSid {
				continue
			}

			if msg.Pres != nil {
				// Skip notifying - already notified on topic.
				if msg.Pres.SkipTopic != "" && sess.getSub(msg.Pres.SkipTopic) != nil {
					continue
				}

				// Notification addressed to a single user only.
				if msg.Pres.SingleUser != "" && pssd.uid.UserId() != msg.Pres.SingleUser {
					continue
				}
				// Notification should skip a single user.
				if msg.Pres.ExcludeUser != "" && pssd.uid.UserId() == msg.Pres.ExcludeUser {
					continue
				}

				// Check presence filters
				if !t.passesPresenceFilters(msg.Pres, pssd.uid) {
					continue
				}

			} else {
				if msg.Info != nil {
					// Don't forward read receipts and key presses to channel readers and those without the R permission.
					// OK to forward with Src != "" because it's sent from another topic to 'me', permissions already
					// checked there.
					if msg.Info.Src == "" && (pssd.isChanSub || !t.userIsReader(pssd.uid)) {
						continue
					}

					// Skip notifying - already notified on topic.
					if msg.Info.SkipTopic != "" && sess.getSub(msg.Info.SkipTopic) != nil {
						continue
					}

					// Don't send key presses from one user's session to the other sessions of the same user.
					if msg.Info.What == "kp" && msg.Info.From == pssd.uid.UserId() {
						continue
					}

				} else if !t.userIsReader(pssd.uid) && !pssd.isChanSub {
					// Skip {data} if the user has no Read permission and not a channel reader.
					continue
				}
			}
		} else if pssd.isChanSub && types.IsChannel(sess.sid) {
			// If it's a chnX multiplexing session, check if there's a corresponding
			// grpX multiplexing session as we don't want to send the message to both.
			grpSid := types.ChnToGrp(sess.sid)
			if grpSess := globals.sessionStore.Get(grpSid); grpSess != nil && grpSess.isMultiplex() {
				// If grpX multiplexing session's attached to topic, skip this chnX session
				// (message will be routed to the topic proxy via the grpX session).
				if _, attached := t.sessions[grpSess]; attached {
					continue
				}
			}
		}

		// Make a copy of msg since messages sent to sessions differ.
		msgCopy := msg.copy()
		// Topic name may be different depending on the user to which the `sess` belongs.
		t.prepareBroadcastableMessage(msgCopy, pssd.uid, pssd.isChanSub)
		// Send message to session.
		if !sess.queueOut(msgCopy) {
			globals.l.Sugar().Warnf("topic[%s]: connection stuck, detaching - %s", t.name, sess.sid)
			dropSessions = append(dropSessions, sess)
		}
	}

	// Drop "bad" sessions.
	for _, sess := range dropSessions {
		// The whole session is being dropped, so ClientComMessage.init is false.
		// keep redundant init: false so it can be searched for.
		t.unregisterSession(&ClientComMessage{sess: sess, init: false})
	}
}

// unregisterSession implements all logic following receipt of a leave
// request via the Topic.unreg channel.
func (t *Topic) unregisterSession(msg *ClientComMessage) {
	if t.currentCall != nil {
		shouldTerminateCall := false
		if msg.sess.isMultiplex() {
			// Check if any of the call party sessions is multiplexed over msg.sess.
			for _, p := range t.currentCall.parties {
				if p.sess.isProxy() && p.sess.multi == msg.sess {
					shouldTerminateCall = true
					break
				}
			}
		} else if _, found := t.currentCall.parties[msg.sess.sid]; found {
			// Normal session disconnecting from topic. Just terminate the call.
			shouldTerminateCall = true
		}

		if shouldTerminateCall {
			t.terminateCallInProgress(false)
		}
	}

	t.handleLeaveRequest(msg, msg.sess)
	if msg.init && msg.sess.inflightReqs != nil {
		// If it's a client initiated request.
		msg.sess.inflightReqs.Done()
	}

	// If there are no more subscriptions to this topic, start a kill timer
	if len(t.sessions) == 0 && t.cat != types.TopicCatSys {
		t.killTimer.Reset(constants.IdleMasterTopicTimeout)
	}
}

// handleLeaveRequest processes a session leave request.
func (t *Topic) handleLeaveRequest(msg *ClientComMessage, sess *Session) {
	// Remove connection from topic; session may continue to function
	now := types.TimeNow()

	var asUid types.Uid
	var asChan bool
	if msg.init {
		asUid = types.ParseUserId(msg.AsUser)
		var err error

		asChan, err = t.verifyChannelAccess(msg.Original)
		if err != nil {
			// Group topic cannot be addressed as channel unless channel functionality is enabled.
			sess.queueOut(ErrNotFoundReply(msg, now))
		}
	}

	if t.isInactive() {
		if !asUid.IsZero() && msg.init {
			sess.queueOut(ErrLockedReply(msg, now))
		}
		return
	}

	// User wants to leave and unsubscribe.
	if msg.init && msg.Leave.Unsub {
		// asUid must not be Zero.
		if err := t.replyLeaveUnsub(sess, msg, asUid); err != nil {
			globals.l.Sugar().Errorf("failed to unsub", err, sess.sid)
		}
		return
	}

	// User wants to leave without unsubscribing.
	if pssd, _ := t.remSession(sess, asUid); pssd != nil {
		if !sess.isProxy() {
			sess.delSub(t.name)
		}
		if pssd.isChanSub != asChan {
			// Cannot address non-channel subscription as channel and vice versa.
			if msg.init {
				// Group topic cannot be addressed as channel unless channel functionality is enabled.
				sess.queueOut(ErrNotFoundReply(msg, now))
			}
			return
		}

		var uid types.Uid
		if sess.isProxy() {
			// Multiplexing session, multiple UIDs.
			uid = asUid
		} else {
			// Simple session, single UID.
			uid = pssd.uid
		}

		var pud perUserData
		// uid may be zero when a proxy session is trying to terminate (it called unsubAll).
		if !uid.IsZero() {
			// UID not zero: one user removed.
			pud = t.perUser[uid]
			if !sess.background {
				pud.online--
				t.perUser[uid] = pud
			}
		} else if len(pssd.muids) > 0 {
			// UID is zero: multiplexing session is dropped altogether.
			// Using new 'uid' and 'pud' variables.
			for _, uid := range pssd.muids {
				pud := t.perUser[uid]
				pud.online--
				t.perUser[uid] = pud
			}
		} else if !sess.isCluster() {
			globals.l.Sugar().Panicf("cannot determine uid: leave req", msg, sess)
		}

		switch t.cat {
		case types.TopicCatMe:
			mrs := t.mostRecentSession()
			if mrs == nil {
				// Last session
				mrs = sess
			} else {
				// Change UA to the most recent live session and announce it. Don't block.
				select {
				case t.supd <- &sessionUpdate{userAgent: mrs.userAgent}:
				default:
				}
			}

			meUid := uid
			if meUid.IsZero() && len(pssd.muids) > 0 {
				// The entire multiplexing session is being dropped. Need to find owner's UID.
				// len(pssd.muids) could be zero if the session was a background session.
				meUid = pssd.muids[0]
			}
			if !meUid.IsZero() {
				// Update user's last online timestamp & user agent. Only one user can be subscribed to 'me' topic.
				if err := globals.store.UsersUpdateLastSeen(meUid, mrs.userAgent, now); err != nil {
					globals.l.Error("", zap.Error(err))
				}
			}

		case types.TopicCatFnd:
			// FIXME: this does not work correctly in case of a multiplexing session.
			// 	Remove ephemeral query.
			t.fndRemovePublic(sess)
		case types.TopicCatGrp:
			// Subscriber is going offline in the topic: notify other subscribers who are currently online.
			readFilter := &presFilters{filterIn: types.ModeRead}
			if !uid.IsZero() {
				if pud.online == 0 {
					if asChan {
						// Simply delete record from perUserData
						delete(t.perUser, uid)
					} else {
						t.presSubsOnline("off", uid.UserId(), nilPresParams, readFilter, "")
					}
				}
			} else if len(pssd.muids) > 0 {
				for _, uid := range pssd.muids {
					if t.perUser[uid].online == 0 {
						if asChan {
							// delete record from perUserData
							delete(t.perUser, uid)
						} else {
							t.presSubsOnline("off", uid.UserId(), nilPresParams, readFilter, "")
						}
					}
				}
			}
		}

		if !uid.IsZero() {
			// Respond if contains an id.
			if msg.init {
				sess.queueOut(NoErrReply(msg, now))
			}
		}
	}
}

// replyLeaveUnsub is request to unsubscribe user and detach all user's sessions from topic.
func (t *Topic) replyLeaveUnsub(sess *Session, msg *ClientComMessage, asUid types.Uid) error {
	now := types.TimeNow()

	if asUid.IsZero() {
		panic("replyLeaveUnsub: zero asUid")
	}

	if t.owner == asUid {
		if msg.init {
			sess.queueOut(ErrPermissionDeniedReply(msg, now))
		}
		return errors.New("replyLeaveUnsub: owner cannot unsubscribe")
	}

	var err error
	var asChan bool
	if msg.init {
		asChan, err = t.verifyChannelAccess(msg.Original)
		if err != nil {
			sess.queueOut(ErrNotFoundReply(msg, now))
			return errors.New("replyLeaveUnsub: incorrect addressing of channel")
		}
	}

	pud := t.perUser[asUid]
	// Delete user's subscription from the database; msg could be nil, so cannot use msg.Original.
	if pud.isChan {
		// Handle channel reader.
		err = globals.store.SubsDelete(types.GrpToChn(t.name), asUid)
	} else {
		// Handle subscriber.
		err = globals.store.SubsDelete(t.name, asUid)
	}

	if err != nil {
		if msg.init {
			if err == types.ErrNotFound {
				sess.queueOut(InfoNoActionReply(msg, now))
				err = nil
			} else {
				sess.queueOut(ErrUnknownReply(msg, now))
			}
		}
		return err
	}

	if msg.init {
		sess.queueOut(NoErrReply(msg, now))
	}

	var oldWant types.AccessMode
	var oldGiven types.AccessMode
	if !asChan {
		// Update cached unread count: negative value
		if (pud.modeWant & pud.modeGiven).IsReader() {
			usersUpdateUnread(asUid, pud.readID-t.lastID, true)
		}
		oldWant, oldGiven = pud.modeWant, pud.modeGiven
	} else {
		oldWant, oldGiven = types.ModeCChnReader, types.ModeCChnReader
		// Unsubscribe user's devices from the channel (FCM topic).
		t.channelSubUnsub(asUid, false)
	}

	// Send prsence notifictions to admins, other users, and user's other sessions.
	t.notifySubChange(asUid, asUid, asChan, oldWant, oldGiven, types.ModeUnset, types.ModeUnset, sess.sid)

	// Evict all user's sessions, clear cached data, send notifications.
	t.evictUser(asUid, true, sess.sid)

	// Notify plugins.
	pluginSubscription(&types.Subscription{Topic: t.name, User: asUid.String()}, plgActDel)

	// If all P2P users were deleted, suspend the topic to let it shut down.
	if t.cat == types.TopicCatP2P && t.subsCount() == 0 {
		t.markPaused(true)
		globals.hub.unreg <- &topicUnreg{del: true, sess: nil, rcptTo: t.name, pkt: nil}
	}

	return nil
}

// FIXME: this won't work correctly with multiplexing sessions.
func (t *Topic) mostRecentSession() *Session {
	var sess *Session
	var latest int64
	for s := range t.sessions {
		sessionLastAction := atomic.LoadInt64(&s.lastAction)
		if sessionLastAction > latest {
			sess = s
			latest = sessionLastAction
		}
	}
	return sess
}

// Remove per-session value of fnd.Public.
func (t *Topic) fndRemovePublic(sess *Session) {
	if t.public == nil {
		return
	}

	// FIXME: case of a multiplexing session won't work correctly.
	// Maybe handle it at the proxy topic.
	if pubmap, ok := t.public.(map[string]any); ok {
		delete(pubmap, sess.sid)
		return
	}
	panic("Invalid Fnd.Public type")
}
