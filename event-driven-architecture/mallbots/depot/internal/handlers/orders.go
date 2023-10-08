package handlers

import (
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/depot/internal/application"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/depot/internal/domain"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/ddd"
)

func RegisterOrderHandlers(orderHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.ShoppingListCompleted{}, orderHandlers.OnShoppingListCompleted)
}
