package command_template

import (
	"context"
	"github.com/goccy/go-json"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_command_template"
)

func rdbSub(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start commandTemplate rdbSub --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pubsub := o.rdb.Subscribe(ctx, "sendCommandTemplate")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			l.Errorln(err)
		}
		b := []byte(msg.Payload)
		var sc e_command_template.SendCommandTemplate
		err = json.Unmarshal(b, &sc)
		if err != nil {
			l.Errorln("send data is not correctly")
		}
		c := o.generateCommand(sc)
		_, err = o.commandS.ExecuteReturnId(ctx, c)
		if err != nil {
			l.Errorln("Error executing command from commandTemplate")
		}
	}
}
