package command_server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"schedule_task_command/app/dbs"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
	"testing"
	"time"
)

func setUpServer() (cs *CommandServer) {
	l := logFile.NewLogFile("test", "commandServer.log")
	DBS := dbs.NewDbs(l, true)
	cs = NewCommandServer(DBS)
	return
}

func TestCommandServer(t *testing.T) {
	cs := setUpServer()
	ctx := context.Background()
	go func() { cs.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		com := e_command.Command{
			Token: "test1",
		}
		commandId, err := cs.ExecuteReturnId(ctx, com)
		require.NotEqual(t, commandId, "")
		require.NoError(t, err)
		fmt.Printf("commandId: %s", commandId)
	})
	t.Run("test2", func(t *testing.T) {
		e := util.MyErr("test err")
		com := e_command.Command{
			Token:   "test2",
			Message: &e,
		}
		ctx := context.Background()
		commandId, err := cs.ExecuteReturnId(ctx, com)
		require.Error(t, err)
		require.Equal(t, commandId, "")
		fmt.Println(commandId)
	})
}
