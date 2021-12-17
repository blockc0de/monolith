package storage

import (
	"encoding/json"
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

func (m *GraphsManager) GetLogs(hash string) ([]types.Log, error) {
	val, err := m.RedisClient.Lrange(m.graphLogsKey(hash), 0, -1)
	if err != nil {
		return nil, err
	}

	result := make([]types.Log, 0, len(val))
	for _, v := range val {
		var log types.Log
		err = json.Unmarshal([]byte(v), &log)
		if err != nil {
			return nil, err
		}
		result = append(result, log)
	}
	return result, nil
}

func (m *GraphsManager) AppendLog(hash string, message types.Log) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = m.RedisClient.Pipelined(func(pipeliner redis.Pipeliner) error {
		pipeliner.RPush(m.graphLogsKey(hash), data)

		pipeliner.LTrim(m.graphLogsKey(hash), -50, -1)

		_, err := pipeliner.Exec()
		return err
	})
	return err
}

func (m *GraphsManager) ClearLogs(hash string) error {
	_, err := m.RedisClient.Del(m.graphLogsKey(hash))
	return err
}
