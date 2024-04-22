package push

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/push/fcm"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/push/types"
)

type Push struct {
	handlers      map[string]types.Handler
	HandlersNames []string
}

func NewPush(cfg config.PushConfig) (*Push, error) {
	var p Push
	p.handlers["fcm"] = fcm.NewFcm()

	var enabled []string
	all := map[string]interface{}{}

	if cfg.FCM != nil {
		all["fcm"] = cfg.FCM
	}

	for name, cc := range all {
		if hand := p.handlers[name]; hand != nil {
			if ok, err := hand.Init(cc); err != nil {
				return nil, err
			} else if ok {
				enabled = append(enabled, name)
			}
		}
	}

	p.HandlersNames = enabled
	return &p, nil
}

// Stop all pushes
func (p *Push) Stop() {
	if p.handlers == nil {
		return
	}

	for _, hand := range p.handlers {
		if hand.IsReady() {
			hand.Stop()
		}
	}
}
