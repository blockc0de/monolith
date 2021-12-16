package storage

import (
	"fmt"

	"github.com/tal-tech/go-zero/core/stores/redis"
)

// Redis Keys
// "wallet::{address}::graphs" => set

type WalletsManager struct {
	RedisClient *redis.Redis
}

func (m *WalletsManager) walletGraphsKey(address string) string {
	return fmt.Sprintf("wallet::{%s}::graphs", address)
}

func (m *WalletsManager) AddGraph(address, hash string) error {
	_, err := m.RedisClient.Sadd(m.walletGraphsKey(address), hash)
	return err
}

func (m *WalletsManager) RemoveGraph(address, hash string) error {
	_, err := m.RedisClient.Srem(m.walletGraphsKey(address), hash)
	return err
}
