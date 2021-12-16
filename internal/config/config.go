package config

import (
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/rest"
)

type Auth struct {
	AccessSecret string
	AccessExpire int64
}

type Config struct {
	rest.RestConf
	Auth      Auth
	RedisConf redis.RedisConf
}
