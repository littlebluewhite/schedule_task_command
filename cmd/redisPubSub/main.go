package main

import (
	"context"
	"fmt"
	"schedule_task_command/app/command_server"
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
	c := command_server.SendCommand{
		TemplateId:     1,
		TriggerFrom:    []string{"execute"},
		TriggerAccount: "Wilson",
	}
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	err = r.Publish(ctx, "commandTest", b).Err()
	if err != nil {
		panic(err)
	}
}

func sub(ctx context.Context, r *redis.Client) {
	pubsub := r.Subscribe(ctx, "commandTest")
	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		panic(err)
	}
	m := msg.Payload
	fmt.Println(m)
	b := []byte(m)
	var s command_server.SendCommand
	json.Unmarshal(b, &s)

	fmt.Println(s)
	fmt.Printf("%T, %v\n", s, s)
	fmt.Println("msg.Pattern: ", msg.Pattern)
	fmt.Println(msg.PayloadSlice)
	fmt.Println(msg.Channel)
}
