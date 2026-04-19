package utils

import (
	"context"
	"fmt"
	"miniorder-order-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

func InitConnectREDIS(redisConfig domain.RedisConfig) *redis.Client {
	fmt.Printf("Connecting to Redis : '%v'\n", redisConfig.Host)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host,
		Password: redisConfig.Password,
		DB:       redisConfig.Database,
	})
	result := redisClient.Ping(context.Background())
	if result.Err() != nil {
		panic(result.Err())
	}
	fmt.Printf("Connected to Redis : '%v'\n", redisConfig.Host)
	return redisClient
}
