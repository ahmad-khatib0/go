package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/customers/customerspb"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/ordering/internal/domain"
)

type CustomerRepository struct {
	client customerspb.CustomersServiceClient
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(conn *grpc.ClientConn) CustomerRepository {
	return CustomerRepository{client: customerspb.NewCustomersServiceClient(conn)}
}

func (r CustomerRepository) Authorize(ctx context.Context, customerID string) error {
	_, err := r.client.AuthorizeCustomer(ctx, &customerspb.AuthorizeCustomerRequest{Id: customerID})
	return err
}
