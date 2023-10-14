package application

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/search/internal/models"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}

type CustomerCacheRepository interface {
	Add(ctx context.Context, customerID, name string) error
	CustomerRepository
}
