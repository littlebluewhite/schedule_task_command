package task_template

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task_template"
	"github.com/littlebluewhite/schedule_task_command/util/redis_stream"
	"github.com/redis/go-redis/v9"
)

func rdbSub(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start taskTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendTaskTemplate")

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
		go func(payload string, ctx context.Context) {
			var s e_task_template.SendTaskTemplate
			if e := json.Unmarshal([]byte(payload), &s); e != nil {
				l.Errorln("send data is not correctly")
				return
			}
			_, e := o.Execute(ctx, s)
			if e != nil {
				l.Errorln("Error executing Task from taskTemplate")
			}
		}(msg.Payload, ctx)
	}
}

func receiveStream(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start taskTemplate receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "TaskTemplate", "server", l)
	rs.Start(ctx, o.getStreamComMap())
}
