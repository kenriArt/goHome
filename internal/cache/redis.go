package cache

import (
	"GOTRAKER/internal/orders"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, password string, ttl time.Duration) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &RedisCache{client: rdb, ttl: ttl}
}
func (c *RedisCache) Get(ctx context.Context, id int) (*orders.Order, error) {
	key := fmt.Sprintf("order:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("данные в кэше не найдены")
	}
	fmt.Println("из редис")
	var order orders.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, fmt.Errorf("ошибка десериализации: %w", err)
	}
	return &order, nil
}
func (c *RedisCache) Set(ctx context.Context, order *orders.Order) error {
	key := fmt.Sprintf("order:%d", order.ID)
	val, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("ошибка сериализации заказа: %w", err)
	}
	return c.client.Set(ctx, key, val, c.ttl).Err()
}
func (c *RedisCache) Delete(ctx context.Context, id int) error {
	key := fmt.Sprintf("order:%d", id)
	return c.client.Del(ctx, key).Err()
}
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
