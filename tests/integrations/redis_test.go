package integrations

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

type DummyData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestSetAndGetCache_RedisIntegration(t *testing.T) {
	ctx := context.Background()

	rdb, mock := redismock.NewClientMock()

	data := &DummyData{
		ID:   123,
		Name: "Redis Test",
	}

	key := utils.BuildKey("dummy", data.ID)

	// Setup expectation for Set
	expectedJSON, err := json.Marshal(data)
	assert.NoError(t, err)

	mock.ExpectSet(key, expectedJSON, 5*time.Minute).SetVal("OK")
	utils.SetCache(ctx, rdb, key, data, 5*time.Minute)

	// Setup expectation for Get
	mock.ExpectGet(key).SetVal(string(expectedJSON))
	result, err := utils.GetCache[DummyData](ctx, rdb, key)

	assert.NoError(t, err)
	assert.Equal(t, data, result)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
