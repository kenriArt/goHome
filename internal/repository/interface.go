package repository

import "GOTRAKER/internal/orders"

type OrderRepository interface {
	Add(orders.Order)
	GetByID(id int) (orders.Order, error)
	Update(order orders.Order) error
	GetAll() []orders.Order
}
