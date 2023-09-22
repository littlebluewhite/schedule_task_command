package time_template

import (
	"context"
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_time_template"
	"schedule_task_command/util/logFile"
)

func rdbSub(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start timeTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendTimeTemplate")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		b := []byte(msg.Payload)
		var s e_time_template.SendTimeTemplate
		err = json.Unmarshal(b, &s)
		if err != nil {
			l.Error().Println("send data is not correctly")
		}
		s.TriggerFrom = append(s.TriggerFrom, "redis channel")
		pt := o.generatePublishTime(s)
		_, err = o.timeS.Execute(pt)
		if err != nil {
			l.Error().Println("Error executing time from timeTemplate")
		}
	}
}
