package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func BuildKey(prefix string, parts ...interface{}) string {
	key := prefix
	for _, part := range parts {
		switch v := part.(type) {
		case *uint:
			if v != nil {
				key += fmt.Sprintf(":%d", *v)
			} else {
				key += ":null"
			}
		case *int:
			if v != nil {
				key += fmt.Sprintf(":%d", *v)
			} else {
				key += ":null"
			}
		case *float64:
			if v != nil {
				key += fmt.Sprintf(":%.2f", *v)
			} else {
				key += ":null"
			}
		default:
			key += fmt.Sprintf(":%v", v)
		}
	}
	return key
}

func GetCache[T any](ctx context.Context, rdb *redis.Client, key string) (*T, error) {
	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("⚠️ Redis Get error:", err)
		}
		return nil, err
	}

	var val T
	if err := json.Unmarshal([]byte(data), &val); err != nil {
		log.Println("⚠️ Unmarshal error:", err)
		return nil, err
	}

	return &val, nil
}

func SetCache[T any](ctx context.Context, rdb *redis.Client, key string, value *T, ttl time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("⚠️ Marshal error:", err)
		return
	}

	if err := rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Println("⚠️ Redis Set error:", err)
	}
}

func DeleteCache(ctx context.Context, rdb *redis.Client, key string) error {
	if err := rdb.Del(ctx, key).Err(); err != nil {
		log.Println("⚠️ Redis Delete error:", err)
		return err
	}
	return nil
}

func DeleteCacheByPattern(ctx context.Context, rdb *redis.Client, pattern string) error {
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		log.Println("⚠️ Redis Keys error:", err)
		return err
	}

	if len(keys) == 0 {
		return nil // No keys to delete
	}

	if err := rdb.Del(ctx, keys...).Err(); err != nil {
		log.Println("⚠️ Redis Del error:", err)
		return err
	}

	return nil
}