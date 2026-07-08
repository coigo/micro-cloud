package infra

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func NewRedisConn (ctx context.Context) {
	Redis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0,
	})

	err := Redis.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Conectado ao redis")
}