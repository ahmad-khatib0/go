package server

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"google.golang.org/grpc"
)

func pluginsInit(configString json.RawMessage) {
	// Check if any plugins are defined
	if len(configString) == 0 {
		return
	}

	var config []config.PluginConfig
	if err := json.Unmarshal(configString, &config); err != nil {
		globals.l.Sugar().Fatalf("", err)
	}

	nameIndex := make(map[string]bool)
	globals.plugins = make([]Plugin, len(config))
	count := 0
	for i := range config {
		conf := &config[i]
		if !conf.Enabled {
			continue
		}

		if nameIndex[conf.Name] {
			globals.l.Sugar().Fatalf("plugins: duplicate name '%s'", conf.Name)
		}

		globals.plugins[count] = Plugin{
			name:        conf.Name,
			timeout:     time.Duration(conf.Timeout) * time.Microsecond,
			failureCode: conf.FailureCode,
			failureText: conf.FailureMessage,
		}

		var err error
		if globals.plugins[count].filterFireHose, err =
			ParsePluginFilter(conf.Filters.FireHose, plgFilterByTopicType|plgFilterByPacket); err != nil {
			globals.l.Sugar().Fatalf("plugins: bad FireHose filter", err)
		}
		if globals.plugins[count].filterAccount, err =
			ParsePluginFilter(conf.Filters.Account, plgFilterByAction); err != nil {
			globals.l.Sugar().Fatalf("plugins: bad Account filter", err)
		}
		if globals.plugins[count].filterTopic, err =
			ParsePluginFilter(conf.Filters.Topic, plgFilterByTopicType|plgFilterByAction); err != nil {
			globals.l.Sugar().Fatalf("plugins: bad Topic filter", err)
		}
		if globals.plugins[count].filterSubscription, err =
			ParsePluginFilter(conf.Filters.Subscription, plgFilterByTopicType|plgFilterByAction); err != nil {
			globals.l.Sugar().Fatalf("plugins: bad Subscription filter", err)
		}
		if globals.plugins[count].filterMessage, err =
			ParsePluginFilter(conf.Filters.Message, plgFilterByTopicType|plgFilterByAction); err != nil {
			globals.l.Sugar().Fatalf("plugins: bad Message filter", err)
		}

		globals.plugins[count].filterFind = conf.Filters.Find

		if parts := strings.SplitN(conf.ServiceAddr, "://", 2); len(parts) < 2 {
			globals.l.Sugar().Fatalf("plugins: invalid server address format", conf.ServiceAddr)
		} else {
			globals.plugins[count].network = parts[0]
			globals.plugins[count].addr = parts[1]
		}

		globals.plugins[count].conn, err = grpc.Dial(globals.plugins[count].addr, grpc.WithInsecure())
		if err != nil {
			globals.l.Sugar().Fatalf("plugins: connection failure %v", err)
		}

		globals.plugins[count].client = chat.NewPluginClient(globals.plugins[count].conn)

		nameIndex[conf.Name] = true
		count++
	}

	globals.plugins = globals.plugins[:count]
	if len(globals.plugins) == 0 {
		globals.l.Sugar().Infof("plugins: no active plugins found")
		globals.plugins = nil
	} else {
		var names []string
		for i := range globals.plugins {
			names = append(names, globals.plugins[i].name+"("+globals.plugins[i].addr+")")
		}

		globals.l.Sugar().Infof("plugins: active", "'"+strings.Join(names, "', '")+"'")
	}
}

// ParsePluginFilter parses filter config string.
func ParsePluginFilter(s *string, filterBy int) (*PluginFilter, error) {
	if s == nil {
		return nil, nil
	}

	parseByName := func(parts []string, options []string, def int) (int, error) {
		var result int

		// Iterate over filter parts
		for _, inp := range parts {
			if inp != "" {
				inp = strings.ToLower(inp)
				// Split string like "hi,login,pres" or "me,p2p,fnd"
				values := strings.Split(inp, ",")
				// For each value in the input string, try to find it in the options set
				for _, val := range values {
					i := 0
					// Iterate over the options, i.e find "hi" in the slice of packet names
					for i = range options {
						if options[i] == val {
							result |= 1 << uint(i)
							break
						}
					}

					if result != 0 && i == len(options) {
						// Mix of known and unknown options in the input
						return 0, errors.New("plugin: unknown value in filter " + val)
					}
				}

				if result != 0 {
					// Found and parsed the right part
					break
				}
			}
		}

		// If the filter value is not defined, use default.
		if result == 0 {
			result = def
		}

		return result, nil
	}

	parseAction := func(parts []string) int {
		var result int
		for _, inp := range parts {
		Loop:
			for _, char := range inp {
				switch char {
				case 'c', 'C':
					result |= plgActCreate
				case 'u', 'U':
					result |= plgActUpd
				case 'd', 'D':
					result |= plgActDel
				default:
					// Unknown symbol means this is not an action string.
					result = 0
					break Loop
				}
			}

			if result != 0 {
				// Found and parsed actions.
				break
			}
		}
		if result == 0 {
			result = plgActMask
		}
		return result
	}

	filter := PluginFilter{}
	parts := strings.Split(*s, ";")
	var err error

	if filterBy&plgFilterByPacket != 0 {
		if filter.byPacket, err = parseByName(parts, plgPacketNames, plgClientMask); err != nil {
			return nil, err
		}
	}

	if filterBy&plgFilterByTopicType != 0 {
		if filter.byTopicType, err = parseByName(parts, plgTopicCatNames, plgTopicCatMask); err != nil {
			return nil, err
		}
	}

	if filterBy&plgFilterByAction != 0 {
		filter.byAction = parseAction(parts)
	}

	return &filter, nil
}

// Message accepted for delivery
func pluginMessage(data *MsgServerData, action int) {
	if globals.plugins == nil || action != plgActCreate {
		return
	}

	var event *chat.MessageEvent
	for i := range globals.plugins {
		p := &globals.plugins[i]
		if p.filterMessage == nil || p.filterMessage.byAction&action == 0 {
			// Plugin is not interested in Message actions
			continue
		}

		if event == nil {
			event = &chat.MessageEvent{
				Action: pluginActionToCrud(action),
				Msg:    pbServDataSerialize(data).Data,
			}
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if p.timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
			defer cancel()
		} else {
			ctx = context.Background()
		}

		if _, err := p.client.Message(ctx, event); err != nil {
			globals.l.Sugar().Warnf("plugins: Message call failed", p.name, err)
		}
	}
}

func pluginsShutdown() {
	if globals.plugins == nil {
		return
	}

	for i := range globals.plugins {
		globals.plugins[i].conn.Close()
	}
}

func pluginGenerateClientReq(sess *Session, msg *ClientComMessage) *chat.ClientReq {
	cmsg := pbCliSerialize(msg)
	if cmsg == nil {
		return nil
	}
	return &chat.ClientReq{
		Msg: cmsg,
		Sess: &chat.Session{
			SessionId:  sess.sid,
			UserId:     sess.uid.UserId(),
			AuthLevel:  chat.AuthLevel(sess.authLvl),
			UserAgent:  sess.userAgent,
			RemoteAddr: sess.remoteAddr,
			DeviceId:   sess.deviceID,
			Language:   sess.lang,
		},
	}
}

func pluginFireHose(sess *Session, msg *ClientComMessage) (*ClientComMessage, *ServerComMessage) {
	if globals.plugins == nil {
		// Return the original message to continue processing without changes
		return msg, nil
	}

	var req *chat.ClientReq

	id, topic := pluginIDAndTopic(msg)
	ts := time.Now().UTC().Round(time.Millisecond)
	for i := range globals.plugins {
		p := &globals.plugins[i]
		if !pluginDoFiltering(p.filterFireHose, msg) {
			// Plugin is not interested in FireHose
			continue
		}

		if req == nil {
			// Generate request only if needed
			req = pluginGenerateClientReq(sess, msg)
			if req == nil {
				// Failed to serialize message. Most likely the message is invalid.
				break
			}
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if p.timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
			defer cancel()
		} else {
			ctx = context.Background()
		}
		if resp, err := p.client.FireHose(ctx, req); err == nil {
			respStatus := resp.GetStatus()
			// CONTINUE means default processing
			if respStatus == chat.RespCode_CONTINUE {
				continue
			}
			// DROP means stop processing of the message
			if respStatus == chat.RespCode_DROP {
				return nil, nil
			}
			// REPLACE: ClientMsg was updated by the plugin. Use the new one for further processing.
			if respStatus == chat.RespCode_REPLACE {
				return pbCliDeserialize(resp.GetClmsg()), nil
			}

			// RESPOND: Plugin provided an alternative response message. Use it
			return nil, pbServDeserialize(resp.GetSrvmsg())

		} else if p.failureCode != 0 {
			// Plugin failed and it's configured to stop further processing.
			globals.l.Sugar().Errorf("plugin: failed,", p.name, err)
			return nil, &ServerComMessage{
				Ctrl: &MsgServerCtrl{
					Id:        id,
					Code:      p.failureCode,
					Text:      p.failureText,
					Topic:     topic,
					Timestamp: ts,
				},
			}
		} else {
			// Plugin failed but configured to ignore failure.
			globals.l.Sugar().Warnf("plugin: failure ignored,", p.name, err)
		}
	}

	return msg, nil
}

// Ask plugin to perform search.
func pluginFind(user types.Uid, query string) (string, []types.Subscription, error) {
	if globals.plugins == nil {
		return query, nil, nil
	}

	find := &chat.SearchQuery{
		UserId: user.UserId(),
		Query:  query,
	}
	for i := range globals.plugins {
		p := &globals.plugins[i]
		if !p.filterFind {
			// Plugin cannot service Find requests
			continue
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if p.timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
			defer cancel()
		} else {
			ctx = context.Background()
		}
		resp, err := p.client.Find(ctx, find)
		if err != nil {
			globals.l.Sugar().Warnf("plugins: Find call failed", p.name, err)
			return "", nil, err
		}
		respStatus := resp.GetStatus()
		// CONTINUE means default processing
		if respStatus == chat.RespCode_CONTINUE {
			continue
		}
		// DROP means stop processing the request
		if respStatus == chat.RespCode_DROP {
			return "", nil, nil
		}
		// REPLACE: query string was changed. Use the new one for further processing.
		if respStatus == chat.RespCode_REPLACE {
			return resp.GetQuery(), nil, nil
		}
		// RESPOND: Plugin provided a specific response. Use it
		return "", pbSubSliceDeserialize(resp.GetResult()), nil
	}

	return query, nil, nil
}

func pluginAccount(user *types.User, action int) {
	if globals.plugins == nil {
		return
	}

	var event *chat.AccountEvent
	for i := range globals.plugins {
		p := &globals.plugins[i]
		if p.filterAccount == nil || p.filterAccount.byAction&action == 0 {
			// Plugin is not interested in Account actions
			continue
		}

		if event == nil {
			event = &chat.AccountEvent{
				Action: pluginActionToCrud(action),
				UserId: user.Uid().UserId(),
				DefaultAcs: pbDefaultAcsSerialize(&MsgDefaultAcsMode{
					Auth: user.Access.Auth.String(),
					Anon: user.Access.Anon.String(),
				}),
				Public: interfaceToBytes(user.Public),
				Tags:   user.Tags,
			}
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if p.timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
			defer cancel()
		} else {
			ctx = context.Background()
		}
		if _, err := p.client.Account(ctx, event); err != nil {
			globals.l.Sugar().Warnf("plugins: Account call failed", p.name, err)
		}
	}
}

func pluginTopic(topic *Topic, action int) {
	if globals.plugins == nil {
		return
	}

	var event *chat.TopicEvent
	for i := range globals.plugins {
		p := &globals.plugins[i]
		if p.filterTopic == nil || p.filterTopic.byAction&action == 0 {
			// Plugin is not interested in Message actions
			continue
		}

		if event == nil {
			event = &chat.TopicEvent{
				Action: pluginActionToCrud(action),
				Name:   topic.name,
				Desc:   pbTopicSerializeToDesc(topic),
			}
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if p.timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
			defer cancel()
		} else {
			ctx = context.Background()
		}
		if _, err := p.client.Topic(ctx, event); err != nil {
			globals.l.Sugar().Warnf("plugins: Topic call failed", p.name, err)
		}
	}
}

func pluginSubscription(sub *types.Subscription, action int) {
	if globals.plugins == nil {
		return
	}

	var event *chat.SubscriptionEvent
	for i := range globals.plugins {
		p := &globals.plugins[i]
		if p.filterSubscription == nil || p.filterSubscription.byAction&action == 0 {
			// Plugin is not interested in Message actions
			continue
		}

		if event == nil {
			event = &chat.SubscriptionEvent{
				Action: pluginActionToCrud(action),
				Topic:  sub.Topic,
				UserId: sub.User,

				DelId:  int32(sub.DelId),
				ReadId: int32(sub.ReadSeqId),
				RecvId: int32(sub.RecvSeqId),

				Mode: &chat.AccessMode{
					Want:  sub.ModeWant.String(),
					Given: sub.ModeGiven.String(),
				},

				Private: interfaceToBytes(sub.Private),
			}
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if p.timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
			defer cancel()
		} else {
			ctx = context.Background()
		}
		if _, err := p.client.Subscription(ctx, event); err != nil {
			globals.l.Sugar().Warnf("plugins: Subscription call failed", p.name, err)
		}
	}
}
