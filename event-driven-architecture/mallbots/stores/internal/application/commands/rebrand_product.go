package commands

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/ddd"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/stores/internal/domain"
)

type RebrandProduct struct {
	ID          string
	Name        string
	Description string
}

type RebrandProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRebrandProductHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) RebrandProductHandler {
	return RebrandProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h RebrandProductHandler) RebrandProduct(ctx context.Context, cmd RebrandProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.Rebrand(cmd.Name, cmd.Description)
	if err != nil {
		return err
	}

	err = h.products.Save(ctx, product)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
