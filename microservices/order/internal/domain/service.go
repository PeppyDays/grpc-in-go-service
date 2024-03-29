package domain

type OrderService interface {
	PlaceOrder(order Order) (Order, error)
}
