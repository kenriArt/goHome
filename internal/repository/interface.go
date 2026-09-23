package repository

import "GOTRAKER/internal/orders"

type OrderRepository interface {
	Add(order orders.Order) error
	GetByID(id int) (orders.Order, error)
	Update(order orders.Order) error
	GetAll() []orders.Order
}
