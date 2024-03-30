package domain

type PaymentGateway interface {
	Charge(*Order) error
}
