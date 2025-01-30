package api

import (
	"context"
	"errors"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

type mockedPayment struct {
	mock.Mock // Embeds to track the activity of the payment
}

func (p *mockedPayment) Charge(ctx context.Context, order *domain.Order) error {
	args := p.Called(ctx, order) // Tracks the function call with arguments
	return args.Error(0)         // Tracks the function return values
}

type mockedDb struct {
	mock.Mock
}

func (d *mockedDb) Save(ctx context.Context, order *domain.Order) error {
	args := d.Called(ctx, order)
	return args.Error(0)
}

func (d *mockedDb) Get(ctx context.Context, id int64) (domain.Order, error) {
	args := d.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(1)
}

func TestPlaceOrder(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDb)

	payment.On("Charge", mock.Anything, mock.Anything).Return(nil) // There is no error on payment.Charge
	db.On("Save", mock.Anything, mock.Anything).Return(nil)        // There is no error on db.Save

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "camera",
				UnitPrice:   12.3,
				Quantity:    3,
			},
		},
		CreatedAt: 0,
	})
	assert.Nil(t, err)
}

func Test_Should_Return_Error_When_Db_Persistence_Fail(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDb)
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(errors.New("connection error")) // db.Save() returns a connection error

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "phone",
				UnitPrice:   12.3,
				Quantity:    3,
			},
		},
		CreatedAt: 0,
	})
	assert.EqualError(t, err, "connection error")
}

// There could be an error on the payment.Charge() call, and solving this would be a bit
// complex because it contains a validation error message. Since the message comes
// from the Payment service, we get only the fields we need and return them to the end user
func Test_Should_Return_Error_When_Payment_Fail(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDb)

	payment.On("Charge", mock.Anything).Return(errors.New("insufficient balance"))
	db.On("Save", mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "bag",
				UnitPrice:   2.5,
				Quantity:    6,
			},
		},
		CreatedAt: 0,
	})

	st, _ := status.FromError(err)
	assert.Equal(t, st.Message(), "order creation failed")
	assert.Equal(t, st.Details()[0].(*errdetails.BadRequest).FieldViolations[0].Description, "insufficient balance")
	assert.Equal(t, st.Code(), codes.InvalidArgument)
}
