package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Connect ...
func Connect() *redis.Client {
	address := os.Getenv("REDIS")
	if address == "" {
		address = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}
