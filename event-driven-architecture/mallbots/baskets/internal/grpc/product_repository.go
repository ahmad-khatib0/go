package grpc

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/domain"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/stores/storespb"
	"github.com/stackus/errors"
	"google.golang.org/grpc"
)

type ProductRepository struct {
	client storespb.StoresServiceClient
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(conn *grpc.ClientConn) ProductRepository {
	return ProductRepository{client: storespb.NewStoresServiceClient(conn)}
}

func (r ProductRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	resp, err := r.client.GetProduct(ctx, &storespb.GetProductRequest{
		Id: productID,
	})

	if err != nil {
		return nil, errors.Wrap(err, "requesting product")
	}

	return r.productToDomain(resp.Product), nil
}

func (r ProductRepository) productToDomain(product *storespb.Product) *domain.Product {
	return &domain.Product{
		ID:      product.GetId(),
		StoreID: product.GetStoreId(),
		Name:    product.GetName(),
		Price:   product.GetPrice(),
	}
}