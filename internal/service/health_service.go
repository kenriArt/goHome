package service

import (
	"GOTRAKER/internal/cache"
	"context"

	"github.com/jmoiron/sqlx"
)

type HealthService struct {
	db    *sqlx.DB
	redis *cache.RedisCache
}

func NewHealthService(db *sqlx.DB, redisClient *cache.RedisCache) *HealthService {
	return &HealthService{
		db:    db,
		redis: redisClient,
	}
}

func (s *HealthService) Check(ctx context.Context) map[string]string {
	result := map[string]string{
		"status":   "ok",
		"postgres": "ok",
		"redis":    "ok",
	}

	if err := s.db.PingContext(ctx); err != nil {
		result["postgres"] = "unavailable"
		result["status"] = "unhealthy"
	}

	if err := s.redis.Ping(ctx); err != nil {
		result["redis"] = "unavailable"
		result["status"] = "unhealthy"
	}

	return result
}
