package storage

import (
	"fmt"

	"github.com/tal-tech/go-zero/core/stores/redis"

	"github.com/blockc0de/monolith/internal/types"
)

// Redis Keys
// "graph::{hash}" => map
//     "data" => string
//     "state" => string
// "graph::{hash}::logs" => list

type GraphsManager struct {
	RedisClient *redis.Redis
}

func (m *GraphsManager) graphKey(hash string) string {
	return fmt.Sprintf("graph::{%s}", hash)
}

func (m *GraphsManager) graphLogsKey(hash string) string {
	return fmt.Sprintf("graph::{%s}::logs", hash)
}

func (m *GraphsManager) Get(hash string) (string, error) {
	return m.RedisClient.Hget(m.graphKey(hash), "data")
}

func (m *GraphsManager) Save(hash, data string) error {
	return m.RedisClient.Hset(m.graphKey(hash), "data", data)
}

func (m *GraphsManager) GetState(hash string) (string, error) {
	return m.RedisClient.Hget(m.graphKey(hash), "state")
}

func (m *GraphsManager) SetState(hash, state string) error {
	return m.RedisClient.Hset(m.graphKey(hash), "state", state)
}

func (m *GraphsManager) AppendLog(hash string, message types.Log) (string, error) {
	return "", nil
}
