package task_template

import (
	"context"
	"github.com/goccy/go-json"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util/redis_stream"
)

func rdbSub(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start taskTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendTaskTemplate")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			l.Errorln(err)
		}
		b := []byte(msg.Payload)
		var s e_task_template.SendTaskTemplate
		err = json.Unmarshal(b, &s)
		if err != nil {
			l.Errorln("send data is not correctly")
		}
		_, err = o.Execute(ctx, s)
		if err != nil {
			l.Errorln("Error executing Task from taskTemplate")
		}
	}
}

func receiveStream(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start taskTemplate receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "TaskTemplate", "server", l)
	rs.Start(ctx, o.getStreamComMap())
}
