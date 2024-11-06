package rdb

import (
	"context"
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/util/config"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync"
	"time"
)

var (
	clientInstance redis.UniversalClient
	once           sync.Once
)

func newSingleRedis(redisConfig config.RedisConfig, hostPort [][]string) *redis.Client {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		redisConfig.User, redisConfig.Password, hostPort[0][0], hostPort[0][1], redisConfig.DB)
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		panic(err)
	}

	// 設置連接池選項
	opt.PoolSize = 5000     // 設置最大連接數，根據需求調整
	opt.MinIdleConns = 5    // 設置最小閒置連接數
	opt.MaxIdleConns = 1000 // 設置最大閒置連接數
	opt.DialTimeout = 5 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second

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
		Addrs:        address,
		Username:     redisConfig.User,
		Password:     redisConfig.Password,
		PoolSize:     5000,
		MinIdleConns: 5,
		MaxIdleConns: 1000,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis cluster connect success")
	return rdb
}

func newClient(config config.RedisConfig) redis.UniversalClient {
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

func NewClient(config config.RedisConfig) redis.UniversalClient {
	once.Do(func() {
		clientInstance = newClient(config)
	})
	return clientInstance
}
