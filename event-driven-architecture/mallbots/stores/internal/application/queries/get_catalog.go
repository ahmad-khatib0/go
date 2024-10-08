package queries

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/stores/internal/domain"
)

type GetCatalog struct {
	StoreID string
}

type GetCatalogHandler struct {
	catalog domain.CatalogRepository
}

func NewGetCatalogHandler(catalog domain.CatalogRepository) GetCatalogHandler {
	return GetCatalogHandler{catalog: catalog}
}

func (h GetCatalogHandler) GetCatalog(ctx context.Context, query GetCatalog) ([]*domain.CatalogProduct, error) {
	return h.catalog.GetCatalog(ctx, query.StoreID)
}
