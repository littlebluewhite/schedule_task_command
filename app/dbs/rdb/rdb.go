package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/util/config"
	"strings"
)

func newSingleRedis(redisConfig config.RedisConfig, hostPort [][]string) *redis.Client {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		redisConfig.User, redisConfig.Password, hostPort[0][0], hostPort[0][1], redisConfig.DB)
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

func newClusterRedis(redisConfig config.RedisConfig, hostPort [][]string) *redis.ClusterClient {
	address := make([]string, 0, len(hostPort))
	for _, v := range hostPort {
		address = append(address, fmt.Sprintf("%s:%s", v[0], v[1]))
	}
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    address,
		Username: redisConfig.User,
		Password: redisConfig.Password,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis cluster connect success")
	return rdb
}

func NewClient(config config.RedisConfig) redis.UniversalClient {
	host := strings.Split(config.Host, ",")
	hostPort := make([][]string, 0, len(host))
	for _, v := range host {
		hostPort = append(hostPort, strings.Split(strings.Trim(v, " "), ":"))
	}
	if len(host) == 1 {
		return newSingleRedis(config, hostPort)
	}
	return newClusterRedis(config, hostPort)
}
