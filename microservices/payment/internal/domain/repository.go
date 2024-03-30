package domain

import "context"

type PaymentRepository interface {
	Get(ctx context.Context, id string) (Payment, error)
	Save(ctx context.Context, payment *Payment) error
}
