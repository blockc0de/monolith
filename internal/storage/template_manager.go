package storage

import (
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Redis Keys
// "template::{id}" => string
// "template::next" => integer

type Template struct {
	Title       string `json:"title"`
	Key         string `json:"key"`
	Bytes       string `json:"bytes"`
	Description string `json:"description"`
	CustomImg   string `json:"customImg"`
}

type TemplateManager struct {
	RedisClient *redis.Redis
}

func (m *TemplateManager) redisNextIdKey() string {
	return "template::next"
}

func (m *TemplateManager) redisTemplateKey(id int64) string {
	return fmt.Sprintf("template::{%d}", id)
}

func (m *TemplateManager) Get(id int64) (*Template, error) {
	val, err := m.RedisClient.Get(m.redisTemplateKey(id))
	if err != nil {
		return nil, err
	}

	var template Template
	err = json.Unmarshal([]byte(val), &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func (m *TemplateManager) Add(template Template) error {
	data, err := json.Marshal(template)
	if err != nil {
		return err
	}

	id, err := m.RedisClient.Incr(m.redisNextIdKey())
	if err != nil {
		return err
	}

	return m.RedisClient.Set(m.redisTemplateKey(id), string(data))
}

func (m *TemplateManager) Remove(id int64) error {
	_, err := m.RedisClient.Del(m.redisTemplateKey(id))
	return err
}
