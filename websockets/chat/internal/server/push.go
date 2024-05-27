package server

/******************************************************************************
 *
 *  Description:
 *    Push notifications handling.
 *
 *****************************************************************************/

import (
	"time"

	pt "github.com/ahmad-khatib0/go/websockets/chat/internal/push/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// Subscribe or unsubscribe user to/from FCM topic (channel).
func (t *Topic) channelSubUnsub(uid types.Uid, sub bool) {
	globals.push.ChannelSub(&pt.ChannelReq{
		Uid:     uid,
		Channel: types.GrpToChn(t.name),
		Unsub:   !sub,
	})
}

// Prepares a payload to be delivered to a mobile device as a push notification in response to a {data} message.
func (t *Topic) pushForData(fromUid types.Uid, data *MsgServerData, msgMarkedAsReadBySender bool) *pt.Receipt {
	// Passing `Topic` as `t.name` for group topics and P2P topics. The p2p topic name is later rewritten for
	// each recipient then the payload is created: p2p recipient sees the topic as the ID of the other user.

	// Initialize the push receipt.
	contentType, _ := data.Head["mime"].(string)
	receipt := pt.Receipt{
		To: make(map[types.Uid]pt.Recipient, t.subsCount()),
		Payload: pt.Payload{
			What:        pt.ActMsg,
			Silent:      false,
			Topic:       t.name,
			From:        data.From,
			Timestamp:   data.Timestamp,
			SeqId:       data.SeqId,
			ContentType: contentType,
			Content:     data.Content,
		},
	}

	if webrtc, found := data.Head["webrtc"].(string); found {
		receipt.Payload.Webrtc = webrtc
		audioOnly, _ := data.Head["aonly"].(bool)
		receipt.Payload.AudioOnly = audioOnly
	}

	if replace, found := data.Head["replace"].(string); found {
		receipt.Payload.Replace = replace
	}

	if t.isChan {
		// Channel readers should get a push on a channel name (as an FCM topic push).
		receipt.Channel = types.GrpToChn(t.name)
	}

	for uid, pud := range t.perUser {
		online := pud.online
		if uid == fromUid && online == 0 {
			// Make sure the sender's devices receive a silent push.
			online = 1
		}

		// Send only to those who have notifications enabled.
		mode := pud.modeWant & pud.modeGiven
		if mode.IsPresencer() && mode.IsReader() && !pud.deleted && !pud.isChan {
			receipt.To[uid] = pt.Recipient{
				// Number of attached sessions the data message will be delivered to.
				// Push notifications sent to users with non-zero online sessions will be marked silent.
				Delivered: online,
				// Unread counts are incremented for all recipients,
				// and for sender only if the message wasnt't marked 'read' by the sender
				ShouldIncrementUnreadCountInCache: uid != fromUid || !msgMarkedAsReadBySender,
			}
		}
	}
	if len(receipt.To) > 0 || receipt.Channel != "" {
		return &receipt
	}
	// If there are no recipient there is no need to send the push notification.
	return nil
}

func (t *Topic) preparePushForSubReceipt(fromUid types.Uid, now time.Time) *pt.Receipt {
	// The `Topic` in the push receipt is `t.xoriginal` for group topics, `fromUid` for p2p topics,
	// not the t.original(fromUid) because it's the topic name as seen by the recipient, not by the sender.
	topic := t.xoriginal
	if t.cat == types.TopicCatP2P {
		topic = fromUid.UserId()
	}

	// Initialize the push receipt.
	receipt := &pt.Receipt{
		To: make(map[types.Uid]pt.Recipient, t.subsCount()),
		Payload: pt.Payload{
			What:      pt.ActSub,
			Silent:    false,
			Topic:     topic,
			From:      fromUid.UserId(),
			Timestamp: now,
			SeqId:     t.lastID,
		},
	}
	return receipt
}

// Prepares payload to be delivered to a mobile device as a push notification in response to a new subscription in a p2p topic.
func (t *Topic) pushForP2PSub(fromUid, toUid types.Uid, want, given types.AccessMode, now time.Time) *pt.Receipt {
	receipt := t.preparePushForSubReceipt(fromUid, now)
	receipt.Payload.ModeWant = want
	receipt.Payload.ModeGiven = given

	receipt.To[toUid] = pt.Recipient{}

	return receipt
}

// Prepares payload to be delivered to a mobile device as a push notification in response to a new subscription in a group topic.
func (t *Topic) pushForGroupSub(fromUid types.Uid, now time.Time) *pt.Receipt {
	receipt := t.preparePushForSubReceipt(fromUid, now)
	if pud, ok := t.perUser[fromUid]; ok {
		receipt.Payload.ModeWant = pud.modeWant
		receipt.Payload.ModeGiven = pud.modeGiven
	} else {
		// Sender is not a subscriber (BUG?)
		return nil
	}

	for uid, pud := range t.perUser {
		// Send only to those who have notifications enabled.
		mode := pud.modeWant & pud.modeGiven
		if mode.IsPresencer() && mode.IsReader() && !pud.deleted && !pud.isChan {
			receipt.To[uid] = pt.Recipient{}
		}
	}

	if len(receipt.To) > 0 || receipt.Channel != "" {
		return receipt
	}
	return nil
}

// Prepares payload to be delivered to a mobile device as a push notification in response to owner deleting a channel.
func pushForChanDelete(topicName string, now time.Time) *pt.Receipt {
	topicName = types.GrpToChn(topicName)
	// Initialize the push receipt.
	return &pt.Receipt{
		Payload: pt.Payload{
			What:      pt.ActSub,
			Silent:    true,
			Topic:     topicName,
			Timestamp: now,
			ModeWant:  types.ModeNone,
			ModeGiven: types.ModeNone,
		},
		Channel: topicName,
	}
}

// Prepares payload to be delivered to a mobile device as a push notification in response to receiving "read" notification.
func (t *Topic) pushForReadRcpt(uid types.Uid, seq int, now time.Time) *pt.Receipt {
	// The `Topic` in the push receipt is `t.xoriginal` for group topics, `fromUid` for p2p topics,
	// not the t.original(fromUid) because it's the topic name as seen by the recipient, not by the sender.
	topic := t.xoriginal
	if t.cat == types.TopicCatP2P {
		topic = uid.UserId()
	}

	// Initialize the push receipt.
	receipt := &pt.Receipt{
		To: make(map[types.Uid]pt.Recipient, 1),
		Payload: pt.Payload{
			What:      pt.ActRead,
			Silent:    true,
			Topic:     topic,
			From:      uid.UserId(),
			Timestamp: now,
			SeqId:     seq,
		},
	}
	receipt.To[uid] = pt.Recipient{}
	return receipt
}

// Process push notification.
func sendPush(rcpt *pt.Receipt) {
	if rcpt == nil || globals.usersUpdate == nil {
		return
	}

	var local *UserCacheReq

	// In case of a cluster pushes will be initiated at the nodes which own the users.
	// Sort users into local and remote.
	if globals.cluster != nil {
		local = &UserCacheReq{PushRcpt: &pt.Receipt{
			Payload: rcpt.Payload,
			Channel: rcpt.Channel,
			To:      make(map[types.Uid]pt.Recipient),
		}}
		remote := &UserCacheReq{PushRcpt: &pt.Receipt{
			Payload: rcpt.Payload,
			Channel: rcpt.Channel,
			To:      make(map[types.Uid]pt.Recipient),
		}}

		for uid, recipient := range rcpt.To {
			if globals.cluster.isRemoteTopic(uid.UserId()) {
				remote.PushRcpt.To[uid] = recipient
			} else {
				local.PushRcpt.To[uid] = recipient
			}
		}

		if len(remote.PushRcpt.To) > 0 || remote.PushRcpt.Channel != "" {
			globals.cluster.routeUserReq(remote)
		}
	} else {
		local = &UserCacheReq{PushRcpt: rcpt}
	}

	if len(local.PushRcpt.To) > 0 || local.PushRcpt.Channel != "" {
		select {
		case globals.usersUpdate <- local:
		default:
		}
	}
}
