package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/app/dbs/rdb"
	"schedule_task_command/util"
	"schedule_task_command/util/config"
	"sync"
	"time"
)

func main() {
	w := &sync.WaitGroup{}
	redisConfig := config.RedisConfig{
		Host:      "127.0.0.1",
		Ports:     []string{"6379"},
		DB:        "0",
		IsCluster: false,
	}
	r := rdb.NewClient(redisConfig)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Add(2)
	go func() {
		subT(ctx, r)
		w.Done()
	}()
	w.Wait()
	fmt.Println("ok")
}

func subT(ctx context.Context, r redis.UniversalClient) {
	pubsub := r.Subscribe(ctx, "timeRec")
	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		panic(err)
	}
	m := msg.Payload
	fmt.Println(m)
	b := []byte(m)
	var p publishTime
	_ = json.Unmarshal(b, &p)

	fmt.Println(p)
	fmt.Printf("%T, %v\n", p, p)
	fmt.Println("msg.Pattern: ", msg.Pattern)
	fmt.Println(msg.PayloadSlice)
	fmt.Println(msg.Channel)
}

type publishTime struct {
	TemplateId     int         `json:"template_id"`
	TriggerFrom    []string    `json:"trigger_from"`
	TriggerAccount string      `json:"trigger_account"`
	Token          string      `json:"token"`
	Time           time.Time   `json:"time"`
	IsTime         bool        `json:"is_time"`
	Status         Status      `json:"status"`
	Message        *util.MyErr `json:"message"`
}

type Status int

const (
	Prepared Status = iota
	Success
	Failure
)

func (s Status) String() string {
	return [...]string{
		"Prepared", "Success", "Failure"}[s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
