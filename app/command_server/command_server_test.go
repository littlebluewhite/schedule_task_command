package command_server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"schedule_task_command/api/group/command"
	"schedule_task_command/app/dbs"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
	"testing"
	"time"
)

func setUpServer() (cs *CommandServer, o *command.Operate) {
	l := logFile.NewLogFile("test", "commandServer.log")
	DBS := dbs.NewDbs(l, true)
	cs = NewCommandServer(DBS)
	o = command.NewOperate(cs)
	return
}

func TestExecuteReturnId(t *testing.T) {
	cs, _ := setUpServer()
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
		commandId, err := cs.ExecuteReturnId(ctx, com)
		require.Error(t, err)
		require.Equal(t, commandId, "")
		fmt.Println(commandId)
	})
}

func TestReadCommand(t *testing.T) {
	cs, o := setUpServer()
	ctx := context.Background()
	go func() { cs.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		sl, err := o.List()
		fmt.Println(sl)
		require.NoError(t, err)
	})
}
