package ports

import (
	"context"
	"github.com/ahmad-khatib0/go/grpc/microservices/order/internal/application/core/domain"
)

type DBPort interface {
	Get(ctx context.Context, id int64) (domain.Order, error)
	Save(context.Context, *domain.Order) error
}
