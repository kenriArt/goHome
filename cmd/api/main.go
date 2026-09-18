package main

import (
	"GOTRAKER/internal/orders"
	"GOTRAKER/internal/repository"
	"GOTRAKER/internal/service"
	"fmt"
)

// ваваыаыавФФФФФФФ4444444
func main() {
	repo := repository.NewInMemoryOrderRepository()
	svc := service.NewOrderService(repo)
	repo.Add(orders.Order{
		ID:          1,
		Customer:    "Vasya",
		Address:     "Moscow",
		IsDelivered: false})
	repo.Add(orders.Order{
		ID:          2,
		Customer:    "Petya",
		Address:     "Moscow",
		IsDelivered: false,
	})
	repo.Add(orders.Order{
		ID:          3,
		Customer:    "Olya",
		Address:     "Kazan",
		IsDelivered: false,
	})
	fmt.Println("[Before]")
	svc.PrintAllOrders()

	svc.DeliverMany([]int{2, 3, 4})

	fmt.Println("\n[After]")
	svc.PrintAllOrders()
	fmt.Println("[Пошел в пень!]")
}
