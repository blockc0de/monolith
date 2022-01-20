package storage

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/blockc0de/monolith/internal/types"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Redis Keys
// "graph::{hash}" => map
//     "data" => string
//     "state" => string
//     "owner" => string
// "graph::logs::{hash}" => list

type Graph struct {
	Hash  string
	Data  string
	State string
	Owner string
}

type GraphsManager struct {
	RedisClient *redis.Redis
}

func (m *GraphsManager) graphKey(hash string) string {
	return fmt.Sprintf("graph::{%s}", hash)
}

func (m *GraphsManager) graphLogsKey(hash string) string {
	return fmt.Sprintf("graph::logs::{%s}", hash)
}

func (m *GraphsManager) Get(hash string) (string, error) {
	return m.RedisClient.Hget(m.graphKey(hash), "data")
}

func (m *GraphsManager) Save(hash, owner, data string) error {
	key := m.graphKey(hash)
	return m.RedisClient.Pipelined(func(pipeliner redis.Pipeliner) error {
		pipeliner.HSet(key, "data", data)
		pipeliner.HSet(key, "owner", owner)

		_, err := pipeliner.Exec()
		return err
	})
}

func (m *GraphsManager) GetState(hash string) (string, error) {
	return m.RedisClient.Hget(m.graphKey(hash), "state")
}

func (m *GraphsManager) SetState(hash, state string) error {
	return m.RedisClient.Hset(m.graphKey(hash), "state", state)
}

func (m *GraphsManager) Scan(cursor uint64, count int64) ([]Graph, uint64, error) {
	keys, cur, err := m.RedisClient.Scan(cursor, "graph::{*", count)
	if err != nil {
		return nil, 0, err
	}

	result := make([]Graph, 0, len(keys))
	for _, key := range keys {
		hash := strings.SplitN(strings.SplitN(key, "{", 2)[1], "}", 2)[0]
		mapper, err := m.RedisClient.Hgetall(m.graphKey(hash))
		if err != nil {
			return nil, 0, err
		}

		graph := Graph{Hash: hash}
		for k, v := range mapper {
			switch k {
			case "data":
				graph.Data = v
			case "owner":
				graph.Owner = v
			case "state":
				graph.State = v
			}
		}
		result = append(result, graph)
	}

	return result, cur, nil
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
