package ports

import (
	"context"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/payment/internal/application/core/domain"
)

type DBPort interface {
	Get(ctx context.Context, id string) (domain.Payment, error)
	Save(ctx context.Context, payment *domain.Payment) error
}
