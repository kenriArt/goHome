package queue

import (
	"GOTRAKER/internal/orders"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
		GroupID: "orders-group",
	})
	fmt.Println("👂 Подписчик слушает Kafka topic 'orders'...")
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}
		var red orders.Order
		json.Unmarshal(msg.Value, &red)
		log.Printf("[Kafka] Новый заказ: %v", red)
	}
}
