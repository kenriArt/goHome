package queue

import (
	"GOTRAKER/internal/orders"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 1 * time.Millisecond,
		},
	}
}
func (p *Producer) SendOrderCreated(ctx context.Context, order *orders.Order) error {
	r, err := json.Marshal(order)
	if err != nil {
		return err
	}
	fmt.Println("Отправили публикацию в Kafka topic 'orders'...")
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: r,
	})
}
