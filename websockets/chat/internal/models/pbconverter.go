package models

import (
	"encoding/json"
	"log"
	"time"

	pbx "github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
)

// Convert ServerComMessage to pbx.ServerMsg
func PbServSerialize(msg *ServerComMessage) *pbx.ServerMsg {
	var pkt pbx.ServerMsg

	switch {
	case msg.Ctrl != nil:
		pkt.Message = pbServCtrlSerialize(msg.Ctrl)
	case msg.Data != nil:
		pkt.Message = pbServDataSerialize(msg.Data)
	case msg.Pres != nil:
		pkt.Message = pbServPresSerialize(msg.Pres)
	case msg.Info != nil:
		pkt.Message = pbServInfoSerialize(msg.Info)
	case msg.Meta != nil:
		pkt.Message = pbServMetaSerialize(msg.Meta)
	}

	pkt.Topic = msg.RcptTo

	return &pkt
}

func pbServCtrlSerialize(ctrl *MsgServerCtrl) *pbx.ServerMsg_Ctrl {
	var params map[string][]byte
	if ctrl.Params != nil {
		if in, ok := ctrl.Params.(map[string]any); ok {
			params = interfaceMapToByteMap(in)
		}
	}

	return &pbx.ServerMsg_Ctrl{
		Ctrl: &pbx.ServerCtrl{
			Id:     ctrl.Id,
			Topic:  ctrl.Topic,
			Code:   int32(ctrl.Code),
			Text:   ctrl.Text,
			Params: params,
		},
	}
}

func pbServDataSerialize(data *MsgServerData) *pbx.ServerMsg_Data {
	return &pbx.ServerMsg_Data{
		Data: &pbx.ServerData{
			Topic:      data.Topic,
			FromUserId: data.From,
			Timestamp:  timeToInt64(&data.Timestamp),
			DeletedAt:  timeToInt64(data.DeletedAt),
			SeqId:      int32(data.SeqId),
			Head:       interfaceMapToByteMap(data.Head),
			Content:    interfaceToBytes(data.Content),
		},
	}
}

func pbServPresSerialize(pres *MsgServerPres) *pbx.ServerMsg_Pres {
	var what pbx.ServerPres_What
	switch pres.What {
	case "on":
		what = pbx.ServerPres_ON
	case "off":
		what = pbx.ServerPres_OFF
	case "ua":
		what = pbx.ServerPres_UA
	case "upd":
		what = pbx.ServerPres_UPD
	case "gone":
		what = pbx.ServerPres_GONE
	case "acs":
		what = pbx.ServerPres_ACS
	case "term":
		what = pbx.ServerPres_TERM
	case "msg":
		what = pbx.ServerPres_MSG
	case "read":
		what = pbx.ServerPres_READ
	case "recv":
		what = pbx.ServerPres_RECV
	case "del":
		what = pbx.ServerPres_DEL
	case "tags":
		what = pbx.ServerPres_TAGS
	default:
		log.Println("Unknown pres.what value", pres.What)
	}
	return &pbx.ServerMsg_Pres{
		Pres: &pbx.ServerPres{
			Topic:        pres.Topic,
			Src:          pres.Src,
			What:         what,
			UserAgent:    pres.UserAgent,
			SeqId:        int32(pres.SeqId),
			DelId:        int32(pres.DelId),
			DelSeq:       pbDelQuerySerialize(pres.DelSeq),
			TargetUserId: pres.AcsTarget,
			ActorUserId:  pres.AcsActor,
			Acs:          pbAccessModeSerialize(pres.Acs),
		},
	}
}

func pbServInfoSerialize(info *MsgServerInfo) *pbx.ServerMsg_Info {
	return &pbx.ServerMsg_Info{
		Info: &pbx.ServerInfo{
			Topic:      info.Topic,
			FromUserId: info.From,
			Src:        info.Src,
			What:       pbInfoNoteWhatSerialize(info.What),
			SeqId:      int32(info.SeqId),
			Event:      pbCallEventSerialize(info.Event),
			Payload:    info.Payload,
		},
	}
}

func pbServMetaSerialize(meta *MsgServerMeta) *pbx.ServerMsg_Meta {
	return &pbx.ServerMsg_Meta{
		Meta: &pbx.ServerMeta{
			Id:    meta.Id,
			Topic: meta.Topic,
			Desc:  pbTopicDescSerialize(meta.Desc),
			Sub:   pbTopicSubSliceSerialize(meta.Sub),
			Del:   pbDelValuesSerialize(meta.Del),
			Tags:  meta.Tags,
			Cred:  pbServerCredsSerialize(meta.Cred),
		},
	}
}

// interfaceMapToByteMap() marshales map values
func interfaceMapToByteMap(in map[string]any) map[string][]byte {
	out := make(map[string][]byte, len(in))
	for key, val := range in {
		if val != nil {
			out[key], _ = json.Marshal(val)
		}
	}
	return out
}

func interfaceToBytes(in any) []byte {
	if in != nil {
		out, _ := json.Marshal(in)
		return out
	}
	return nil
}

func timeToInt64(ts *time.Time) int64 {
	if ts != nil {
		return ts.UnixNano() / int64(time.Millisecond)
	}
	return 0
}

func pbDelQuerySerialize(in []MsgDelRange) []*pbx.SeqRange {
	if in == nil {
		return nil
	}

	out := make([]*pbx.SeqRange, len(in))
	for i, dq := range in {
		out[i] = &pbx.SeqRange{Low: int32(dq.LowID), Hi: int32(dq.HiID)}
	}

	return out
}

func pbAccessModeSerialize(acs *MsgAccessMode) *pbx.AccessMode {
	if acs == nil {
		return nil
	}

	return &pbx.AccessMode{
		Want:  acs.Want,
		Given: acs.Given,
	}
}

func pbInfoNoteWhatSerialize(what string) pbx.InfoNote {
	var out pbx.InfoNote
	switch what {
	case "kp":
		out = pbx.InfoNote_KP
	case "read":
		out = pbx.InfoNote_READ
	case "recv":
		out = pbx.InfoNote_RECV
	case "call":
		out = pbx.InfoNote_CALL
	default:
		log.Println("unknown info-note.what", what)
	}
	return out

}

func pbCallEventSerialize(event string) pbx.CallEvent {
	var out pbx.CallEvent
	switch event {
	case "accept":
		out = pbx.CallEvent_ACCEPT
	case "answer":
		out = pbx.CallEvent_ANSWER
	case "hang-up":
		out = pbx.CallEvent_HANG_UP
	case "ice-candidate":
		out = pbx.CallEvent_ICE_CANDIDATE
	case "invite":
		out = pbx.CallEvent_INVITE
	case "offer":
		out = pbx.CallEvent_OFFER
	case "ringing":
		out = pbx.CallEvent_RINGING
	case "":
		out = pbx.CallEvent_X2
	default:
		log.Println("unknown call event", event)
	}
	return out
}

func pbTopicDescSerialize(desc *MsgTopicDesc) *pbx.TopicDesc {
	if desc == nil {
		return nil
	}
	out := &pbx.TopicDesc{
		CreatedAt: timeToInt64(desc.CreatedAt),
		UpdatedAt: timeToInt64(desc.UpdatedAt),
		TouchedAt: timeToInt64(desc.TouchedAt),
		State:     desc.State,
		Online:    desc.Online,
		IsChan:    desc.IsChan,
		Defacs:    pbDefaultAcsSerialize(desc.DefaultAcs),
		Acs:       pbAccessModeSerialize(desc.Acs),
		SeqId:     int32(desc.SeqId),
		ReadId:    int32(desc.ReadSeqId),
		RecvId:    int32(desc.RecvSeqId),
		DelId:     int32(desc.DelId),
		Public:    interfaceToBytes(desc.Public),
		Trusted:   interfaceToBytes(desc.Trusted),
		Private:   interfaceToBytes(desc.Private),
	}
	if desc.LastSeen != nil {
		out.LastSeenTime = timeToInt64(desc.LastSeen.When)
		out.LastSeenUserAgent = desc.LastSeen.UserAgent
	}
	return out
}

func pbDefaultAcsSerialize(defacs *MsgDefaultAcsMode) *pbx.DefaultAcsMode {
	if defacs == nil {
		return nil
	}
	return &pbx.DefaultAcsMode{Auth: defacs.Auth, Anon: defacs.Anon}
}

func pbTopicSubSliceSerialize(subs []MsgTopicSub) []*pbx.TopicSub {
	if len(subs) == 0 {
		return nil
	}

	out := make([]*pbx.TopicSub, len(subs))
	for i := 0; i < len(subs); i++ {
		out[i] = pbTopicSubSerialize(&subs[i])
	}
	return out
}

func pbTopicSubSerialize(sub *MsgTopicSub) *pbx.TopicSub {
	out := &pbx.TopicSub{
		UpdatedAt: timeToInt64(sub.UpdatedAt),
		DeletedAt: timeToInt64(sub.DeletedAt),
		Online:    sub.Online,
		Acs:       pbAccessModeSerialize(&sub.Acs),
		ReadId:    int32(sub.ReadSeqId),
		RecvId:    int32(sub.RecvSeqId),
		Public:    interfaceToBytes(sub.Public),
		Trusted:   interfaceToBytes(sub.Trusted),
		Private:   interfaceToBytes(sub.Private),
		UserId:    sub.User,
		Topic:     sub.Topic,
		TouchedAt: timeToInt64(sub.TouchedAt),
		SeqId:     int32(sub.SeqId),
		DelId:     int32(sub.DelId),
	}
	if sub.LastSeen != nil {
		out.LastSeenTime = timeToInt64(sub.LastSeen.When)
		out.LastSeenUserAgent = sub.LastSeen.UserAgent
	}
	return out
}

func pbDelValuesSerialize(in *MsgDelValues) *pbx.DelValues {
	if in == nil {
		return nil
	}

	return &pbx.DelValues{
		DelId:  int32(in.DelId),
		DelSeq: pbDelQuerySerialize(in.DelSeq),
	}
}

func pbServerCredsSerialize(in []*MsgCredServer) []*pbx.ServerCred {
	if in == nil {
		return nil
	}

	out := make([]*pbx.ServerCred, len(in))
	for i, cr := range in {
		out[i] = &pbx.ServerCred{Method: cr.Method, Value: cr.Value}
	}

	return out
}
