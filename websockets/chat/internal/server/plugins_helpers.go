package server

import "github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"

func pluginActionToCrud(action int) chat.Crud {
	switch action {
	case plgActCreate:
		return chat.Crud_CREATE
	case plgActUpd:
		return chat.Crud_UPDATE
	case plgActDel:
		return chat.Crud_DELETE
	}
	panic("plugin: unknown action")
}

// Returns false to skip, true to process
func pluginDoFiltering(filter *PluginFilter, msg *ClientComMessage) bool {
	filterByTopic := func(topic string, flt int) bool {
		if topic == "" || flt == plgTopicCatMask {
			return true
		}

		tt := topic
		if len(tt) > 3 {
			tt = topic[:3]
		}
		switch tt {
		case "me":
			return flt&plgTopicMe != 0
		case "fnd":
			return flt&plgTopicFnd != 0
		case "usr":
			return flt&plgTopicP2P != 0
		case "grp":
			return flt&plgTopicGrp != 0
		case "new":
			return flt&plgTopicNew != 0
		}
		return false
	}

	// Check if plugin has any filters for this call
	if filter == nil || filter.byPacket == 0 {
		return false
	}
	// Check if plugin wants all the messages
	if filter.byPacket == plgClientMask && filter.byTopicType == plgTopicCatMask {
		return true
	}
	// Check individual bits
	if msg.Hi != nil {
		return filter.byPacket&plgHi != 0
	}
	if msg.Acc != nil {
		return filter.byPacket&plgAcc != 0
	}
	if msg.Login != nil {
		return filter.byPacket&plgLogin != 0
	}
	if msg.Sub != nil {
		return filter.byPacket&plgSub != 0 && filterByTopic(msg.Sub.Topic, filter.byTopicType)
	}
	if msg.Leave != nil {
		return filter.byPacket&plgLeave != 0 && filterByTopic(msg.Leave.Topic, filter.byTopicType)
	}
	if msg.Pub != nil {
		return filter.byPacket&plgPub != 0 && filterByTopic(msg.Pub.Topic, filter.byTopicType)
	}
	if msg.Get != nil {
		return filter.byPacket&plgGet != 0 && filterByTopic(msg.Get.Topic, filter.byTopicType)
	}
	if msg.Set != nil {
		return filter.byPacket&plgSet != 0 && filterByTopic(msg.Set.Topic, filter.byTopicType)
	}
	if msg.Del != nil {
		return filter.byPacket&plgDel != 0 && filterByTopic(msg.Del.Topic, filter.byTopicType)
	}
	if msg.Note != nil {
		return filter.byPacket&plgNote != 0 && filterByTopic(msg.Note.Topic, filter.byTopicType)
	}
	return false
}

// pluginIDAndTopic extracts message ID and topic name.
func pluginIDAndTopic(msg *ClientComMessage) (string, string) {
	if msg.Hi != nil {
		return msg.Hi.Id, ""
	}
	if msg.Acc != nil {
		return msg.Acc.Id, ""
	}
	if msg.Login != nil {
		return msg.Login.Id, ""
	}
	if msg.Sub != nil {
		return msg.Sub.Id, msg.Sub.Topic
	}
	if msg.Leave != nil {
		return msg.Leave.Id, msg.Leave.Topic
	}
	if msg.Pub != nil {
		return msg.Pub.Id, msg.Pub.Topic
	}
	if msg.Get != nil {
		return msg.Get.Id, msg.Get.Topic
	}
	if msg.Set != nil {
		return msg.Set.Id, msg.Set.Topic
	}
	if msg.Del != nil {
		return msg.Del.Id, msg.Del.Topic
	}
	if msg.Note != nil {
		return "", msg.Note.Topic
	}
	return "", ""
}
