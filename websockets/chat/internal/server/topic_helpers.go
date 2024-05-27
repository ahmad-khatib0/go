package server

import (
	"sync/atomic"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// isInactive checks if topic is paused or being deleted.
func (t *Topic) isInactive() bool {
	return (atomic.LoadInt32(&t.status) & (topicStatusPaused | topicStatusMarkedDeleted)) != 0
}

func (t *Topic) isReadOnly() bool {
	return (atomic.LoadInt32(&t.status) & topicStatusReadOnly) != 0
}

func (t *Topic) isLoaded() bool {
	return (atomic.LoadInt32(&t.status) & topicStatusLoaded) != 0
}

func (t *Topic) isDeleted() bool {
	return (atomic.LoadInt32(&t.status) & topicStatusMarkedDeleted) != 0
}

// Get topic name suitable for the given client
func (t *Topic) original(uid types.Uid) string {
	if t.cat == types.TopicCatP2P {
		if pud, ok := t.perUser[uid]; ok {
			return pud.topicName
		}
		panic("Invalid P2P topic")
	}

	if t.cat == types.TopicCatGrp && t.isChan {
		if t.perUser[uid].isChan {
			// This is a channel reader.
			return types.GrpToChn(t.xoriginal)
		}
	}
	return t.xoriginal
}

// passesPresenceFilters applies presence filters to `msg`
//
// depending on per-user want and given acls for the provided `uid`.
func (t *Topic) passesPresenceFilters(pres *MsgServerPres, uid types.Uid) bool {
	modeWant, modeGiven := t.getPerUserAcs(uid)

	// "gone" and "acs" notifications are sent even if the topic is muted.
	return ((modeGiven & modeWant).IsPresencer() || pres.What == "gone" || pres.What == "acs") &&
		(pres.FilterIn == 0 || int(modeGiven&modeWant)&pres.FilterIn != 0) &&
		(pres.FilterOut == 0 || int(modeGiven&modeWant)&pres.FilterOut == 0)
}

// getPerUserAcs returns `want` and `given` permissions for the given user id.
func (t *Topic) getPerUserAcs(uid types.Uid) (types.AccessMode, types.AccessMode) {
	if uid.IsZero() {
		// For zero uids (typically for proxy sessions), return the union of all permissions.
		return t.modeWantUnion, t.modeGivenUnion
	}
	pud := t.perUser[uid]
	return pud.modeWant, pud.modeGiven
}

// userIsReader returns true if the user (specified by `uid`) may read the given topic.
func (t *Topic) userIsReader(uid types.Uid) bool {
	modeWant, modeGiven := t.getPerUserAcs(uid)

	return (modeGiven & modeWant).IsReader()
}

// prepareBroadcastableMessage sets the topic field in `msg` depending on the uid and subscription type.
func (t *Topic) prepareBroadcastableMessage(msg *ServerComMessage, uid types.Uid, isChanSub bool) {
	// We are only interested in broadcastable messages.
	if msg.Data == nil && msg.Pres == nil && msg.Info == nil {
		return
	}

	if (t.cat == types.TopicCatP2P && !uid.IsZero()) || (t.cat == types.TopicCatGrp && t.isChan) {
		// For p2p topics topic name is dependent on receiver.
		// Channel topics may be presented as grpXXX or chnXXX.

		var topicName string
		if isChanSub {
			topicName = types.GrpToChn(t.xoriginal)
		} else {
			topicName = t.original(uid)
		}

		switch {
		case msg.Data != nil:
			msg.Data.Topic = topicName
		case msg.Pres != nil:
			msg.Pres.Topic = topicName
		case msg.Info != nil:
			msg.Info.Topic = topicName
		}
	}

	// Send channel messages anonymously.
	if isChanSub && msg.Data != nil {
		msg.Data.From = ""
	}
}

// Verifies if topic can be access by the provided name: access any topic
//
// as non-channel, access channel as channel. Returns true if access is for channel,
//
// false if not and error if access is invalid.
func (t *Topic) verifyChannelAccess(asTopic string) (bool, error) {
	if !types.IsChannel(asTopic) {
		return false, nil
	}
	if t.isChan {
		return true, nil
	}
	return false, types.ErrNotFound
}

// Check if topic has any online (non-background) users.
func (t *Topic) isOnline() bool {
	// Find at least one non-background session.
	for s, pssd := range t.sessions {
		if s.isMultiplex() && len(pssd.muids) > 0 {
			return true
		}
		if !s.background {
			return true
		}
	}
	return false
}

// Infer topic category from name.
func topicCat(name string) types.TopicCat {
	return types.GetTopicCat(name)
}

// Generate random string as a name of the group topic
func genTopicName() string {
	return "grp" + globals.store.UidGen.GetStr()
}

// Convert expanded (routable) topic name into name suitable for sending to the user.
// For example p2pAbCDef123 -> usrAbCDef
func topicNameForUser(name string, uid types.Uid, isChan bool) string {
	switch topicCat(name) {
	case types.TopicCatMe:
		return "me"
	case types.TopicCatFnd:
		return "fnd"
	case types.TopicCatP2P:
		topic, _ := types.P2PNameForUser(uid, name)
		return topic
	case types.TopicCatGrp:
		if isChan {
			return types.GrpToChn(name)
		}
	}
	return name
}

// subsCount returns the number of topic subscribers
func (t *Topic) subsCount() int {
	if t.cat == types.TopicCatP2P {
		count := 0
		for uid := range t.perUser {
			if !t.perUser[uid].deleted {
				count++
			}
		}
		return count
	}
	return len(t.perUser)
}
