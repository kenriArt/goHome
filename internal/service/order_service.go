// Package service реализует бизнес-логику управления заказами.
package service

import (
	"GOTRAKER/internal/repository"
	"fmt"
	"sync"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) PrintAllOrders() {
	for _, order := range s.repo.GetAll() {
		fmt.Printf("ID: %d, Name: %s, Is delivery: %v\n", order.ID, order.Customer, order.IsDelivered)
	}
}
func (s *OrderService) DeliverMany(ids []int) {
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			order, err := s.repo.GetByID(n)
			if err != nil {
				fmt.Printf("Ошибка: %v (ID: %d)\n", err, n)
				return
			}
			order.MarkDelivered()
			if err := s.repo.Update(order); err != nil {
				fmt.Printf("Ошибка при обновлении ID %d: %v\n", n, err)
				return
			}
			fmt.Printf("Order %d delivered\n", n)
		}(id)
	}
	wg.Wait()
}
