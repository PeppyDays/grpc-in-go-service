package application

import (
	"errors"
	"testing"

	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockedPaymentGateway struct {
	mock.Mock
}

func (p *mockedPaymentGateway) Charge(order *domain.Order) error {
	args := p.Called(order)
	return args.Error(0)
}

type mockedOrderRepository struct {
	mock.Mock
}

func (r *mockedOrderRepository) Save(order *domain.Order) error {
	args := r.Called(order)
	return args.Error(0)
}

func (r *mockedOrderRepository) Get(id int64) (domain.Order, error) {
	args := r.Called(id)
	return args.Get(0).(domain.Order), args.Error(1)
}

func Test_Should_Place_Order(t *testing.T) {
	paymentGateway := new(mockedPaymentGateway)
	orderRepository := new(mockedOrderRepository)
	paymentGateway.On("Charge", mock.Anything).Return(nil)
	orderRepository.On("Save", mock.Anything).Return(nil)

	orderService := NewOrderService(orderRepository, paymentGateway)

	_, err := orderService.PlaceOrder(domain.Order{
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
	paymentGateway := new(mockedPaymentGateway)
	orderRepository := new(mockedOrderRepository)
	paymentGateway.On("Charge", mock.Anything).Return(nil)
	orderRepository.On("Save", mock.Anything).Return(errors.New("connection error"))

	orderService := NewOrderService(orderRepository, paymentGateway)
	_, err := orderService.PlaceOrder(domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "phone",
				UnitPrice:   14.7,
				Quantity:    1,
			},
		},
		CreatedAt: 0,
	})

	assert.EqualError(t, err, "connection error")
}

func Test_Should_Return_Error_When_Payment_Fail(t *testing.T) {
	paymentGateway := new(mockedPaymentGateway)
	orderRepository := new(mockedOrderRepository)
	paymentGateway.On("Charge", mock.Anything).Return(errors.New("insufficient balance"))
	orderRepository.On("Save", mock.Anything).Return(nil)

	orderService := NewOrderService(orderRepository, paymentGateway)
	_, err := orderService.PlaceOrder(domain.Order{
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
	assert.Equal(t, st.Code(), codes.InvalidArgument)
}
