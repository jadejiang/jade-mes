package util

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"jade-mes/app/infrastructure/persistence/database"
)

// CacheClient ...
type CacheClient struct {
	*cache.Codec
	*redis.Client
}

// Cache ...
var Cache *CacheClient

// RedisNil ...
var RedisNil = database.RedisNil

// CacheServerDefault 缓存服务器 - 默认
var CacheServerDefault = "default"

func newCacheClient(serverName string) *CacheClient {
	redisClient := getRedisClient(serverName)

	cacheClient := CacheClient{
		Client: redisClient,
		Codec: &cache.Codec{
			Redis: redisClient,
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
	}

	return &cacheClient
}

func getRedisClient(serverName string) *redis.Client {
	var redisClient *redis.Client
	redisClient = database.GetRedis().Default
	return redisClient
}

// Use ...
func (c *CacheClient) Use(serverName string) *CacheClient {
	return newCacheClient(serverName)
}
