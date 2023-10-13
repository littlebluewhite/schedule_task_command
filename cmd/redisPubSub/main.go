package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/api/group/time_template"
	"schedule_task_command/app/dbs/rdb"
	"sync"
	"time"
)

func main() {
	w := &sync.WaitGroup{}
	r := rdb.NewRedis("redis")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Add(2)
	go func() {
		sub(ctx, r)
		w.Done()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		pub(ctx, r)
		w.Done()
	}()
	w.Wait()
	fmt.Println("ok")
}

func pub(ctx context.Context, r *redis.Client) {
	c := time_template.SendTime{
		TemplateId: 2,
		Token:      "test",
	}
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	err = r.Publish(ctx, "sendTimeTemplate", b).Err()
	if err != nil {
		panic(err)
	}
}

func sub(ctx context.Context, r *redis.Client) {
	pubsub := r.Subscribe(ctx, "timeRec")
	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		panic(err)
	}
	m := msg.Payload
	fmt.Println(m)
}
