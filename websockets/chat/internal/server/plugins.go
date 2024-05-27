package server

import (
	"context"

	"github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
)

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
