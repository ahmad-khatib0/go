package api

import (
	"context"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/application/core/domain"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := a.db.Save(ctx, &order)
	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(ctx, &order)
	if paymentErr != nil {
		st, _ := status.FromError(paymentErr)
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: st.Message(),
		}

		// This example assumes the Payment service returns a simple status object with code and a message, but if it
		// returns the message with details, we may need to extract those field violations separately. In that case,
		// we can use the built-in status.Convert() instead of status.FromError():
		// st := status.Convert(paymentErr)
		// var allErrors []string
		// for _, detail := range st.Details() {
		// 	switch t := detail.(type) {
		// 	case *errdetails.BadRequest:
		// 		for _, violation := range t.GetFieldViolations() {
		// 			allErrors = append(allErrors, violation.Description)
		// 		}
		// 	}
		// }
		// fieldErr := &errdetails.BadRequest_FieldViolation{
		// 	Field:       "payment",
		// 	Description: strings.Join(allErrors, "\n"),
		// }

		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDetails, _ := orderStatus.WithDetails(badReq)

		// grpcurl -d '{"user_id": 123, "order_items":[{"product_code":"sku1", "unit_price": 0.12, "quantity":1}]}' -plaintext localhost:3000
		// ➥ Order/Create
		// ERROR:
		// Code: InvalidArgument
		// Message: order creation failed
		// Details:
		//   {
		//     "@type":"type.googleapis.com/google.rpc.BadRequest","fieldViolations":
		//      [{"field":"payment","description":"failed to charge. invalid billing address "}]
		//   }
		return domain.Order{}, statusWithDetails.Err()
	}
	return order, nil
}

func (a Application) GetOrder(ctx context.Context, id int64) (domain.Order, error) {
	return a.db.Get(ctx, id)
}
