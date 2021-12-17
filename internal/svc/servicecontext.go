package svc

import (
	"github.com/blockc0de/monolith/internal/config"
	"github.com/blockc0de/monolith/internal/graphs"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config         config.Config
	RedisClient    *redis.Redis
	GraphContainer *graphs.Container
}

func NewServiceContext(c config.Config) *ServiceContext {
	options := []redis.Option{
		redis.WithPass(c.RedisConf.Pass),
	}
	if c.RedisConf.Type == redis.ClusterType {
		options = append(options, redis.Cluster())
	}
	redisClient := redis.New(c.RedisConf.Host, options...)

	return &ServiceContext{
		Config:         c,
		RedisClient:    redisClient,
		GraphContainer: graphs.NewContainer(redisClient),
	}
}
