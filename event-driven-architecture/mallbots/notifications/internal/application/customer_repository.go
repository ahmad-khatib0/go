package application

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/notifications/internal/models"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}
