package ports

import (
	"context"
	"github.com/ahmad-khatib0/go/grpc/microservices/payment/internal/application/core/domain"
)

type APIPort interface {
	Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error)
}
