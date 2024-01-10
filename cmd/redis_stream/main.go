package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/util/redis_stream"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       15, // use default DB
	})
	//groupSetID(rdb)
	//destroyGroup(rdb)
	//createGroup(rdb, "0")
	//groupR(rdb)
	//i := add(rdb)
	//fmt.Println(i)
	values := redis_stream.CreateRedisStreamCom()
	data := map[string]interface{}{
		"id":      1,
		"stages":  "success",
		"status":  2,
		"message": nil,
	}
	jd, _ := json.Marshal(data)
	values["command"] = "track_task"
	values["timestamp"] = time.Now().Unix()
	values["data"] = string(jd)
	values["is_wait_call_back"] = 0
	values["callback_token"] = ""
	values["send_pattern"] = "1"
	e := redis_stream.StreamAdd(context.Background(), rdb, "AlarmAPIModuleReceiver", values)
	fmt.Println(e)
	//createGroup(rdb, i)
	//infoGroup(rdb)
	//infoConsumer(rdb)
	//del(rdb, i)
	//groupR(rdb)
	//go groupR(rdb)

	//time.Sleep(2 * time.Second)
	//add(rdb)
	//time.Sleep(2 * time.Second)
	//add(rdb)
	//time.Sleep(2 * time.Second)
	//add(rdb)
	//time.Sleep(2 * time.Second)
	//add(rdb)
	//time.Sleep(2 * time.Second)
	//add(rdb)
	//time.Sleep(10 * time.Second)
}

func add(rdb redis.UniversalClient) string {
	ctx := context.Background()
	r, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "AlarmAPIModuleReceiver",
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
	}).Result()
	fmt.Printf("%+v\n", r)
	if err != nil {
		panic(err)
	}
	return r
}

func del(rdb redis.UniversalClient, id string) {
	err := rdb.XDel(context.Background(), "my-stream", id).Err()
	if err != nil {
		panic(err)
	}
}

func trim(rdb redis.UniversalClient) {
	ctx := context.Background()
	err := rdb.XTrimMaxLen(ctx, "my-stream", 2).Err()
	if err != nil {
		panic(err)
	}
}

func createGroup(rdb redis.UniversalClient, start string) {
	ctx := context.Background()
	err := rdb.XGroupCreate(ctx, "my-stream", "schedule", start).Err()
	if err != nil {
		panic(err)
	}
}

func groupR(rdb redis.UniversalClient) {
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
		e := json.Unmarshal([]byte(s), &en)
		if e != nil {
			panic(e)
		}
		fmt.Println(en)
	}
}

func destroyGroup(rdb redis.UniversalClient) {
	ctx := context.Background()
	r, err := rdb.XGroupDestroy(ctx, "my-stream", "schedule").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}

func pend(rdb redis.UniversalClient) {
	ctx := context.Background()
	r, err := rdb.XPending(ctx, "my-stream", "schedule").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}

func groupSetID(rdb redis.UniversalClient) {
	ctx := context.Background()
	r, err := rdb.XGroupSetID(ctx, "my-stream", "schedule", "1698649259746-0").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}

func infoGroup(rdb redis.UniversalClient) {
	ctx := context.Background()
	r, err := rdb.XInfoGroups(ctx, "my-stream").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)
}

func infoConsumer(rdb redis.UniversalClient) {
	ctx := context.Background()
	r, err := rdb.XInfoConsumers(ctx, "my-stream", "schedule").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)
}

type Entry struct {
	SourceID    int     `json:"source_id"`
	TableID     int     `json:"table_id"`
	SourceValue string  `json:"source_value"`
	Timestamp   float64 `json:"timestamp"`
}
