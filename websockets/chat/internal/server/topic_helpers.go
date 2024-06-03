package server

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

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
