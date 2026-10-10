// Package service реализует бизнес-логику управления заказами.
package service

import (
	"GOTRAKER/internal/cache"
	"GOTRAKER/internal/orders"
	"GOTRAKER/internal/queue"
	"GOTRAKER/internal/repository"
	"context"
	"fmt"
	"sync"
)

type OrderService struct {
	repo     repository.OrderRepository
	postgres *repository.PostgresOrderRepo
	rds      *cache.RedisCache
	producer *queue.Producer
}

func NewOrderService(repo repository.OrderRepository, postgres *repository.PostgresOrderRepo, rds *cache.RedisCache, producer *queue.Producer) *OrderService {
	return &OrderService{repo: repo, postgres: postgres, rds: rds, producer: producer}
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
func (s *OrderService) AddOrder(ctx context.Context, o *orders.Order) error {
	s.repo.Add(o)
	err := s.postgres.Add(o)
	if err != nil {
		return err
	}
	s.rds.Delete(ctx, o.ID)
	s.producer.SendOrderCreated(ctx, o)
	return nil
}
func (s *OrderService) GetAll() []orders.Order {
	return s.repo.GetAll()
}
func (s *OrderService) GetByID(ctx context.Context, id int) (*orders.Order, error) {
	rds, err := s.rds.Get(ctx, id)
	if err != nil {
		res, err := s.postgres.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if err := s.rds.Set(ctx, res); err != nil {
			return nil, err
		}
		return res, nil
	}
	return rds, nil
}
func (s *OrderService) UpdateOrder(ctx context.Context, o orders.Order) error {
	s.repo.Update(o)
	s.postgres.Update(o)
	s.rds.Delete(ctx, o.ID)
	return nil
}
