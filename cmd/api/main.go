package main

import (
	apphttp "GOTRAKER/internal/http"
	"GOTRAKER/internal/orders"
	"GOTRAKER/internal/repository"
	"GOTRAKER/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("не удалось подключиться к БД:", err)
	}
	repo := repository.NewPostgresOrderRepo(db)
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
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
