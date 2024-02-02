package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/util/config"
)

func newRedis(redisConfig config.RedisConfig) *redis.Client {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		redisConfig.User, redisConfig.Password, redisConfig.Host, redisConfig.Ports[0], redisConfig.DB)
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
	address := make([]string, 0, len(redisConfig.Ports))
	for _, port := range redisConfig.Ports {
		address = append(address, fmt.Sprintf("%s:%s", redisConfig.Host, port))
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

func NewClient(config config.RedisConfig) redis.UniversalClient {
	if config.IsCluster {
		return newClusterRedis(config)
	}
	return newRedis(config)
}
