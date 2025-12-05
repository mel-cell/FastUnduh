package config

import (
	"context"
	"time"

	"github.com/alicebob/miniredis/v2"
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
	// 1. Coba connect ke Redis asli dulu
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(Ctx, 1*time.Second)
	defer cancel()

	_, err := Rdb.Ping(ctx).Result()
	if err == nil {
		// Berhasil connect ke Redis asli
		return
	}

	// 2. Jika gagal, gunakan Miniredis (In-Memory)
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	
	// Print info biar user tau
	// (Kita pakai fmt.Println karena log standard mungkin belum di-setup di main)
	println("⚠️  Redis asli tidak ditemukan. Menggunakan In-Memory Redis (Miniredis).")
	println("⚠️  Data akan hilang saat aplikasi dimatikan.")
}
