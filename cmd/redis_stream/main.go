package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       5,  // use default DB
	})
	groupSetID(rdb)
	go groupR(rdb)
	//ctx := context.Background()
	//groupR(rdb)
	infoGroup(rdb)
	time.Sleep(2 * time.Second)
	add(rdb)
	time.Sleep(2 * time.Second)
	add(rdb)
	time.Sleep(10 * time.Second)
}

func add(rdb *redis.Client) {
	ctx := context.Background()
	err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "my-stream",
		Values: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
	}).Err()
	if err != nil {
		panic(err)
	}
}

func trim(rdb *redis.Client) {
	ctx := context.Background()
	err := rdb.XTrimMaxLen(ctx, "my-stream", 2).Err()
	if err != nil {
		panic(err)
	}
}

func createGroup(rdb *redis.Client) {
	ctx := context.Background()
	err := rdb.XGroupCreate(ctx, "my-stream", "dd", "0").Err()
	if err != nil {
		panic(err)
	}
}

func groupR(rdb *redis.Client) {
	ctx := context.Background()
	for {
		re, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "schedule",
			Consumer: "consumer",
			Streams:  []string{"my-stream", ">"},
			Count:    1,
			Block:    0,
			NoAck:    true,
		}).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(re)
	}
}

func pend(rdb *redis.Client) {
	ctx := context.Background()
	r, err := rdb.XPending(ctx, "my-stream", "schedule").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}

func groupSetID(rdb *redis.Client) {
	ctx := context.Background()
	r, err := rdb.XGroupSetID(ctx, "my-stream", "schedule", "0").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}

func infoGroup(rdb *redis.Client) {
	ctx := context.Background()
	r, err := rdb.XInfoGroups(ctx, "my-stream").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
