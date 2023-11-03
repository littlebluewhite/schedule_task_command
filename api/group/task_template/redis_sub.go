package task_template

import (
	"context"
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util/logFile"
	"schedule_task_command/util/redis_stream"
)

func rdbSub(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start taskTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendTaskTemplate")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		b := []byte(msg.Payload)
		var s e_task_template.SendTaskTemplate
		err = json.Unmarshal(b, &s)
		if err != nil {
			l.Error().Println("send data is not correctly")
		}
		s.TriggerFrom = append(s.TriggerFrom, "task_template redis channel")
		task := o.generateTask(s)
		_, err = o.taskS.ExecuteReturnId(ctx, task)
		if err != nil {
			l.Error().Println("Error executing Task from taskTemplate")
		}
	}
}

func receiveStream(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start taskTemplate receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "TaskTemplate", "server", l)
	rs.Start(ctx, o.getStreamComMap())
}
