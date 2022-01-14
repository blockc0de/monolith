package graphs

import (
	"crypto/tls"

	red "github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	defaultDatabase = 1
	maxRetries      = 3
	idleConns       = 8
)

func newRedisClient(
	addr, pass, nodeType string, enableTLS bool) red.Cmdable {

	if nodeType == redis.ClusterType {
		return newRedisClusterClient(addr, pass, enableTLS)
	}
	return newRedisNodeClient(addr, pass, enableTLS)
}

func newRedisNodeClient(addr, pass string, enableTLS bool) *red.Client {
	var tlsConfig *tls.Config
	if enableTLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return red.NewClient(&red.Options{
		Addr:         addr,
		Password:     pass,
		DB:           defaultDatabase,
		MaxRetries:   maxRetries,
		MinIdleConns: idleConns,
		TLSConfig:    tlsConfig,
	})
}

func newRedisClusterClient(addr, pass string, enableTLS bool) *red.ClusterClient {
	var tlsConfig *tls.Config
	if enableTLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return red.NewClusterClient(&red.ClusterOptions{
		Addrs:        []string{addr},
		Password:     pass,
		MaxRetries:   maxRetries,
		MinIdleConns: idleConns,
		TLSConfig:    tlsConfig,
	})
}
