package handlers

import (
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/ddd"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/ordering/internal/application"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/ordering/internal/domain"
)

func RegisterNotificationHandlers(notificationHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderCreated{}, notificationHandlers.OnOrderCreated)
	domainSubscriber.Subscribe(domain.OrderReadied{}, notificationHandlers.OnOrderReadied)
	domainSubscriber.Subscribe(domain.OrderCanceled{}, notificationHandlers.OnOrderCanceled)
}
