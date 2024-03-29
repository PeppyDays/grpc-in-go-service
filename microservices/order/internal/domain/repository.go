package domain

type OrderRepository interface {
	Get(id string) (Order, error)
	Save(*Order) error
}
