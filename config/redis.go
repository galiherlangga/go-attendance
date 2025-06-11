package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_HOST", "localhost") + ":" + GetEnv("REDIS_PORT", "6379"),
		Password: GetEnv("REDIS_PASSWORD", ""), // no password for local dev
		DB:       0,                            // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Println("Redis not connected:", err)
	} else {
		log.Println("Redis connected")
	}
}

func SetRedisClient(client *redis.Client) {
	RedisClient = client
}