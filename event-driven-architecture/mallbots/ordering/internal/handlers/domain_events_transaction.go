package handlers

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/ddd"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/di"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/ordering/internal/constants"
)

func RegisterDomainEventHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		domainHandlers := di.Get(ctx, constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.Event])

		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterDomainEventHandlers(subscriber, handlers)
}
