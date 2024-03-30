package domain

import "context"

type PaymentService interface {
	Charge(ctx context.Context, payment Payment) (Payment, error)
}
