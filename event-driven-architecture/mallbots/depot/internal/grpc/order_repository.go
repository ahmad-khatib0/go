package grpc

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/depot/internal/domain"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/ordering/orderingpb"
	"google.golang.org/grpc"
)

type OrderRepository struct {
	client orderingpb.OrderingServiceClient
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(conn *grpc.ClientConn) OrderRepository {
	return OrderRepository{client: orderingpb.NewOrderingServiceClient(conn)}
}

func (r OrderRepository) Ready(ctx context.Context, orderID string) error {
	_, err := r.client.ReadyOrder(ctx, &orderingpb.ReadyOrderRequest{Id: orderID})
	return err
}
