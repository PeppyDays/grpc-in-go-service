package application

import (
	"context"

	"github.com/peppydays/grpc-in-go-service/microservices/payment/internal/domain"
)

type PaymentService struct {
	repository domain.PaymentRepository
}

func NewPaymentService(repository domain.PaymentRepository) *PaymentService {
	return &PaymentService{
		repository: repository,
	}
}

func (s PaymentService) Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error) {
	err := s.repository.Save(ctx, &payment)
	if err != nil {
		return domain.Payment{}, err
	}

	return payment, nil
}
