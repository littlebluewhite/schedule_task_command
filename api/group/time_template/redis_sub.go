package time_template

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/api"
	"schedule_task_command/util/redis_stream"
)

func rdbSub(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start timeTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendTimeTemplate")

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
			var s SendTime
			if e := json.Unmarshal([]byte(payload), &s); e != nil {
				l.Errorln("send data is not correctly")
				return
			}
			pt := o.generatePublishTime(s)
			_, e := o.timeS.Execute(pt)
			if e != nil {
				l.Errorln("Error executing Task from taskTemplate")
			}
		}(msg.Payload, ctx)
	}
}

func receiveStream(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start timeTemplate receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "TimeTemplate", "server", l)
	rs.Start(ctx, o.getStreamComMap())
}
