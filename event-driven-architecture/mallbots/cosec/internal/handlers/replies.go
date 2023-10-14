package handlers

import (
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/cosec/internal"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/cosec/internal/models"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/am"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/registry"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/sec"
)

func NewReplyHandlers(reg registry.Registry, orchestrator sec.Orchestrator[*models.CreateOrderData], mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewReplyHandler(reg, orchestrator, mws...)
}

func RegisterReplyHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(internal.CreateOrderReplyChannel, handlers, am.GroupName("cosec-replies"))
	return err
}
