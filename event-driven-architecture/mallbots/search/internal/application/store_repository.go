package application

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/search/internal/models"
)

type StoreRepository interface {
	Find(ctx context.Context, storeID string) (*models.Store, error)
}

type StoreCacheRepository interface {
	Add(ctx context.Context, storeID, name string) error
	Rename(ctx context.Context, storeID, name string) error
	StoreRepository
}
