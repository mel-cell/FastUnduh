package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

const (
	MaxWorkers     = 5
	QueueKey       = "queue:downloads"
	JobPrefix      = "job:"
	ResultPrefix   = "result:"
	FileTTL        = 15 * time.Minute
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // In Docker: "redis:6379"
		DB:   0,
	})
}
