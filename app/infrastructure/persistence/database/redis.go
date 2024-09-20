package database

import (
	"fmt"
	"log"

	goRedis "github.com/go-redis/redis"
	"github.com/spf13/viper"
)

type redisConfig struct {
	Host     string
	Port     uint
	Password string
	DB       int
	DBs      *struct {
		Default    int
		DataReport int
	}
}

type redisConfigs struct {
	Default *redisConfig
}

// RedisClient ...
type RedisClient struct {
	// Default 默认
	Default *goRedis.Client
}

const (
	// RedisNil ...
	RedisNil = goRedis.Nil
)

var redisClient *RedisClient

// Close ...
func (client *RedisClient) Close() {
	if client.Default != nil {
		if err := client.Default.Close(); err != nil {
			log.Printf("close connection of default client exception: %s", err)
		}
	}
}

// GetRedis returns redis instance
func GetRedis() *RedisClient {
	if redisClient == nil {
		log.Fatalln(errNotInited)
	}

	return redisClient
}

func loadRedisConf(config *redisConfig, db int) *goRedis.Options {
	return &goRedis.Options{
		Addr:     fmt.Sprintf("%s:%v", config.Host, config.Port),
		Password: config.Password,
		DB:       db,
	}
}

// TODO: 需要找个更合理的方式管理redis的配置，尤其是在有多个服务器的情况下
func initRedis(settings *viper.Viper) {
	var cfg redisConfigs

	if err := settings.UnmarshalKey("redis", &cfg); err != nil {
		log.Fatalf("load redis config error: %v", err)
	}

	defaultRedis := goRedis.NewClient(loadRedisConf(cfg.Default, cfg.Default.DBs.Default))

	if err := defaultRedis.Ping().Err(); err != nil {
		log.Fatalf("base redis ping error: %v", err)
	}

	redisClient = &RedisClient{
		Default: defaultRedis,
	}
}
