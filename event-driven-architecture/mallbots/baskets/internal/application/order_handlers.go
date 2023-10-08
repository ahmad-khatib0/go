package application

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/domain"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/ddd"
)

type OrderHandlers struct {
	orders domain.OrderRepository
	ignoreUnimplementedDomainEvents
}

var _ DomainEventHandlers = (*OrderHandlers)(nil)

func NewOrderHandlers(orders domain.OrderRepository) OrderHandlers {
	return OrderHandlers{
		orders: orders,
	}
}

func (h OrderHandlers) OnBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	checkedOut := event.(*domain.BasketCheckedOut)
	_, err := h.orders.Save(ctx, checkedOut.Basket)
	return err
}
