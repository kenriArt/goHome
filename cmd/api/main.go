package main

import (
	apphttp "GOTRAKER/internal/http"
	"GOTRAKER/internal/orders"
	"GOTRAKER/internal/repository"
	"GOTRAKER/internal/service"
	"log"
	"net/http"
)

func main() {
	repo := repository.NewInMemoryOrderRepository()
	svc := service.NewOrderService(repo)
	handler := apphttp.NewHandler(svc)
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
	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
