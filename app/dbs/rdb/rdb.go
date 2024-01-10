package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"path/filepath"
	"runtime"
	"schedule_task_command/util/config"
)

var (
	rootPath string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b))))
}

func newRedis(redisConfig config.RedisConfig) *redis.Client {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		redisConfig.User, redisConfig.Password, redisConfig.Host, redisConfig.Port, redisConfig.DB)
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis connect success")
	return rdb
}

func newClusterRedis(redisConfig config.RedisConfig) *redis.ClusterClient {
	address := make([]string, 0, len(redisConfig.ClusterPorts))
	for _, port := range redisConfig.ClusterPorts {
		address = append(address, fmt.Sprintf("%s:%s", redisConfig.ClusterHost, port))
	}
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: address,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis cluster connect success")
	return rdb
}

func NewClient(yamlName string) redis.UniversalClient {
	redisConfig := config.NewConfig[config.RedisConfig](rootPath, "env", yamlName)
	if redisConfig.IsCluster {
		return newClusterRedis(redisConfig)
	}
	return newRedis(redisConfig)
}
