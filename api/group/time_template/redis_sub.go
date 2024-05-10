package time_template

import (
	"context"
	"github.com/goccy/go-json"
	"schedule_task_command/api"
	"schedule_task_command/util/redis_stream"
)

func rdbSub(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start timeTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendTimeTemplate")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			l.Errorln(err)
		}
		b := []byte(msg.Payload)
		var s SendTime
		err = json.Unmarshal(b, &s)
		if err != nil {
			l.Errorln("send data is not correctly")
		}
		pt := o.generatePublishTime(s)
		_, err = o.timeS.Execute(pt)
		if err != nil {
			l.Errorln("Error executing time from timeTemplate")
		}
	}
}

func receiveStream(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start timeTemplate receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "TimeTemplate", "server", l)
	rs.Start(ctx, o.getStreamComMap())
}
