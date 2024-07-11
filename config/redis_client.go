package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var RDB *redis.Client

func initRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	// check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}
}
