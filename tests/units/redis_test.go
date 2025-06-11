package units

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/galiherlangga/go-attendance/config"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	redismock "github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	db, mock := redismock.NewClientMock()
	config.SetRedisClient(db)

	mock.ExpectSet("hello", "world", 10*time.Second).SetVal("OK")
	mock.ExpectGet("hello").SetVal("world")

	// Now your app can call Redis via config.RedisClient
	err := config.RedisClient.Set(context.Background(), "hello", "world", 10*time.Second).Err()
	assert.NoError(t, err)

	val, err := config.RedisClient.Get(context.Background(), "hello").Result()
	assert.NoError(t, err)
	assert.Equal(t, "world", val)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBuildKey(t *testing.T) {
	id := 42
	key := utils.BuildKey("user", &id, "profile")
	assert.Equal(t, "user:42:profile", key)

	var nilInt *int
	assert.Equal(t, "user:null", utils.BuildKey("user", nilInt))
}

func TestSetAndGetCache(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	mockRedis, mock := redismock.NewClientMock()
	ctx := context.Background()

	user := &User{ID: 1, Name: "Alice"}
	key := "user:1"

	// Set expectations
	jsonVal, _ := json.Marshal(user)
	mock.ExpectSet(key, jsonVal, 10*time.Second).SetVal("OK")
	mock.ExpectGet(key).SetVal(string(jsonVal))

	// Call utils.SetCache
	utils.SetCache(ctx, mockRedis, key, user, 10*time.Second)

	// Call utils.GetCache
	result, err := utils.GetCache[User](ctx, mockRedis, key)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}
