package svc

import (
	"github.com/blockc0de/monolith/internal/config"
	"github.com/blockc0de/monolith/internal/graphs"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
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

	container := graphs.NewContainer(redisClient)
	count, err := container.LoadGraphs()
	if err != nil {
		logx.Errorf("Failed to load active graphs, reason: %s", err.Error())
	} else {
		logx.Infof("Load active graphs from storage done, count: %d", count)
	}

	return &ServiceContext{
		Config:         c,
		RedisClient:    redisClient,
		GraphContainer: container,
	}
}
