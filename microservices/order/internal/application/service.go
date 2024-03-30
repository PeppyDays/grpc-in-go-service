package application

import (
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
)

type OrderService struct {
	repository     domain.OrderRepository
	paymentGateway domain.PaymentGateway
}

func NewOrderService(repository domain.OrderRepository, paymentGateway domain.PaymentGateway) *OrderService {
	return &OrderService{
		repository:     repository,
		paymentGateway: paymentGateway,
	}
}

func (s OrderService) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := s.repository.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	err = s.paymentGateway.Charge(&order)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
