package time_template

import (
	"context"
	"github.com/goccy/go-json"
	"schedule_task_command/util/logFile"
	"schedule_task_command/util/redis_stream"
)

func rdbSub(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start timeTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendTimeTemplate")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			l.Error().Println(err)
		}
		b := []byte(msg.Payload)
		var s SendTime
		err = json.Unmarshal(b, &s)
		if err != nil {
			l.Error().Println("send data is not correctly")
		}
		pt := o.generatePublishTime(s)
		_, err = o.timeS.Execute(pt)
		if err != nil {
			l.Error().Println("Error executing time from timeTemplate")
		}
	}
}

func receiveStream(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start timeTemplate receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "TimeTemplate", "server", l)
	rs.Start(ctx, o.getStreamComMap())
}
