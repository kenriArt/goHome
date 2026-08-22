// Package repository отвечает за хранение и получение данных о заказах.
package repository

import (
	"GOTRAKER/internal/orders"
	"sync"
)

type InMemoryOrderRepository struct {
	mu     sync.Mutex
	orders map[int]orders.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{orders: make(map[int]orders.Order)}
}
func (i *InMemoryOrderRepository) Add(order orders.Order) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.orders[order.ID] = order
}
func (i *InMemoryOrderRepository) GetByID(id int) (orders.Order, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	o, ok := i.orders[id]
	if !ok {
		return orders.Order{}, orders.ErrOrderNotFound
	}
	return o, nil
}
func (i *InMemoryOrderRepository) Update(o orders.Order) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.orders[o.ID]
	if !ok {
		return orders.ErrOrderNotFound
	}
	i.orders[o.ID] = o
	return nil
}
func (i *InMemoryOrderRepository) GetAll() []orders.Order {
	i.mu.Lock()
	defer i.mu.Unlock()
	orders := []orders.Order{}
	for _, o := range i.orders {
		orders = append(orders, o)
	}
	return orders
}
