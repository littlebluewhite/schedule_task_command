package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       5,  // use default DB
	})
	//groupSetID(rdb)
	groupSetID(rdb)
	groupR(rdb)
	//go groupR(rdb)
	//infoGroup(rdb)
	//time.Sleep(2 * time.Second)
	//add(rdb)
	//time.Sleep(2 * time.Second)
	//add(rdb)
	//time.Sleep(10 * time.Second)
}

func add(rdb *redis.Client) {
	ctx := context.Background()
	err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "my-stream",
		Values: map[string]interface{}{
			"command":                  "check_time",
			"timestamp":                "1695805466",
			"data":                     "{\"source_id\": 3, \"table_id\": 10618, \"source_value\": \"701\", \"timestamp\": 1695805644.1069221}",
			"callback_command":         "null",
			"callback_channel":         "AlarmAPIModuleReceiver",
			"is_wait_call_back":        "0",
			"callback_token":           "null",
			"callback_timeout":         "5",
			"callback_until_feed_back": "null",
			"command_sk":               "",
			"status_code":              "null",
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
	err := rdb.XGroupCreate(ctx, "my-stream", "schedule", "0").Err()
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
		s := re[0].Messages[0].Values["data"].(string)
		var en Entry
		fmt.Println(s)
		e := json.Unmarshal([]byte(s), &en)
		if e != nil {
			panic(e)
		}
		fmt.Println(en)
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

type Entry struct {
	SourceID    int    `json:"source_id"`
	TableID     int    `json:"table_id"`
	SourceValue string `json:"source_value"`
	Timestamp   int64  `json:"timestamp"`
}
