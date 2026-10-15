package main

import (
	"GOTRAKER/internal/cache"
	apphttp "GOTRAKER/internal/http"
	"GOTRAKER/internal/queue"
	"GOTRAKER/internal/repository"
	"GOTRAKER/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// ваваыаыавФФФФФФФ4444444
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	ttl1 := os.Getenv("REDIS_TTL")
	numTtl, err := strconv.Atoi(ttl1)
	ttl := time.Duration(numTtl) * time.Minute
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("не удалось подключиться к БД:", err)
	}
	producer := queue.NewProducer(os.Getenv("KAFKA_BROKER"), "orders")
	repoPostgres := repository.NewPostgresOrderRepo(db)
	repoInMemory := repository.NewInMemoryOrderRepository()
	rds := cache.NewRedisCache(os.Getenv("REDIS_ADDR"), "password", ttl*time.Minute)

	svc := service.NewOrderService(repoInMemory, repoPostgres, rds, producer)
	handler := apphttp.NewHandler(svc)
	ctx, cnlkafka := context.WithCancel(context.Background())
	defer cnlkafka()
	go queue.StartConsumer(ctx)
	log.Println("🚀 Консьюмер Kafka успешно запущен в фоновой горутине...")

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
