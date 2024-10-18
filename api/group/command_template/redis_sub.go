package command_template

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_command_template"
)

func rdbSub(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start commandTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendCommandTemplate")

	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			l.Errorln(err)
		}
	}(pubsub)
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				l.Infoln("rdbSub receive canceled")
				return
			}
			l.Errorln(err)
			if errors.Is(err, redis.ErrClosed) {
				return
			}
			continue
		}

		// deal with message
		go func(payload string) {
			var sc e_command_template.SendCommandTemplate
			if e := json.Unmarshal([]byte(payload), &sc); e != nil {
				l.Errorln("send data is not correctly")
				return
			}
			c := o.generateCommand(sc)
			_, err = o.commandS.ExecuteReturnId(ctx, c)
			if err != nil {
				l.Errorln("Error executing command from commandTemplate")
			}
		}(msg.Payload)
	}
}
