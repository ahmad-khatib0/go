package ports

import (
	"context"
	"github.com/ahmad-khatib0/go/grpc/microservices/order/internal/application/core/domain"
)

type APIPort interface {
	PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	GetOrder(ctx context.Context, id int64) (domain.Order, error)
}
