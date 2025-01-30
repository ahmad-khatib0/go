package ports

import (
	"context"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/application/core/domain"
)

// The payment port has only one functionality: charge . Simply pass the actual order
// object, and it charges the customer based on order details:
type PaymentPort interface {
	Charge(context.Context, *domain.Order) error
}
