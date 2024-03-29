package application

import (
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
)

type OrderService struct {
	repository domain.OrderRepository
}

func NewOrderService(repository domain.OrderRepository) *OrderService {
	return &OrderService{
		repository: repository,
	}
}

func (s OrderService) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := s.repository.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
