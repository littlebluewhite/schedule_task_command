package command_template

import (
	"context"
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/util/logFile"
)

func rdbSub(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start commandTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendCommandTemplate")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		b := []byte(msg.Payload)
		var sc e_command_template.SendCommandTemplate
		err = json.Unmarshal(b, &sc)
		if err != nil {
			l.Error().Println("send data is not correctly")
		}
		sc.TriggerFrom = append(sc.TriggerFrom, "redis channel")
		c := o.generateCommand(sc)
		_, err = o.commandS.ExecuteReturnId(ctx, c)
		if err != nil {
			l.Error().Println("Error executing command from commandTemplate")
		}
	}
}
