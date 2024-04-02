package domain

type OrderRepository interface {
	Get(id int64) (Order, error)
	Save(*Order) error
}
