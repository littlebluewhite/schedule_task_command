package command_server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"schedule_task_command/api/group/command"
	"schedule_task_command/app/dbs"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
	"sync"
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
	cs, o := setUpServer()
	ctx := context.Background()
	go func() { cs.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		com := e_command.Command{
			Token: "test1",
		}
		commandId, err := cs.ExecuteReturnId(ctx, com)
		require.NotEqual(t, commandId, "")
		require.NoError(t, err)
		fmt.Printf("commandId: %s\n", commandId)
		time.Sleep(1 * time.Second)
		sl, _ := o.List()
		fmt.Printf("data: %+v\n", sl)
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
		sl, _ := o.List()
		fmt.Printf("data: %+v\n", sl)
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

func TestDoCommand(t *testing.T) {
	cs, _ := setUpServer()
	ctx := context.Background()
	go func() { cs.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		h1 := e_command_template.HTTPSCommand{
			Method: e_command_template.GET,
			URL:    "http://192.168.1.10:9330/api/object/value/?id_list=1",
			Header: nil,
		}
		m1 := e_command_template.Monitor{
			StatusCode: 200,
			Interval:   1000,
			MConditions: []e_command_template.MCondition{
				{
					Order:         0,
					CalculateType: ">=",
					SearchRule:    "root.[0]array.value",
					Value:         "2",
				},
			},
		}
		com1 := e_command.Command{
			Token: "test",
			Template: e_command_template.CommandTemplate{
				Name:     "object_test",
				Protocol: e_command_template.Http,
				Timeout:  10000,
				Http:     &h1,
				Monitor:  &m1,
			},
		}

		h2 := e_command_template.HTTPSCommand{
			Method:   e_command_template.PUT,
			URL:      "http://192.168.1.10:9330/api/object/insert_value/",
			Header:   []byte(`[{"key": "test","value": "123456","is_active": true,"data_type": "text"}]`),
			Body:     []byte(`[{"id": 1,"value": "{{value}}"}]`),
			BodyType: e_command_template.Json,
		}
		com2 := e_command.Command{
			Token: "test",
			Template: e_command_template.CommandTemplate{
				Name:     "object_test",
				Protocol: e_command_template.Http,
				Timeout:  10000,
				Http:     &h2,
			},
			Variables: map[string]string{"value": "3"},
		}
		com3 := e_command.Command{
			Token: "test",
			Template: e_command_template.CommandTemplate{
				Name:     "object_test",
				Protocol: e_command_template.Http,
				Timeout:  10000,
				Http:     &h2,
			},
			Variables: map[string]string{"value": "2"},
		}

		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			com2 = cs.ExecuteWait(ctx, com2)
			fmt.Printf("com2: %+v\n", com2)
			fmt.Printf("com2: %+v\n", string(com2.RespData))
			wg.Done()
		}()
		time.Sleep(1 * time.Second)
		comId, _ := cs.ExecuteReturnId(ctx, com1)
		time.Sleep(1 * time.Second)
		com1 = cs.ReadMap()[comId]
		go func() {
			time.Sleep(4 * time.Second)
			com3 = cs.ExecuteWait(ctx, com3)
			fmt.Printf("com2: %+v\n", com2)
			fmt.Printf("com2: %+v\n", string(com2.RespData))
			wg.Done()
		}()
		wg.Wait()
		time.Sleep(1 * time.Second)
		com1 = cs.ReadMap()[comId]
		fmt.Printf("%+v\n", com1)
		fmt.Printf("%+v\n", string(com1.RespData))
		time.Sleep(1 * time.Second)
		err := cs.CancelCommand(comId)
		require.Error(t, err)
	})
	t.Run("command cancel", func(t *testing.T) {
		h := e_command_template.HTTPSCommand{
			Method: e_command_template.GET,
			URL:    "http://192.168.1.10:9330/api/object/value/?id_list=1",
			Header: nil,
		}
		m := e_command_template.Monitor{
			StatusCode: 200,
			Interval:   1000,
			MConditions: []e_command_template.MCondition{
				{
					Order:         0,
					CalculateType: ">=",
					SearchRule:    "root.[0]array.value",
					Value:         "{{value}}",
				},
			},
		}
		com := e_command.Command{
			Token: "test",
			Template: e_command_template.CommandTemplate{
				Name:     "object_test",
				Protocol: e_command_template.Http,
				Timeout:  10000,
				Http:     &h,
				Monitor:  &m,
			},
			Variables: map[string]string{"value": "4"},
		}
		comId, _ := cs.ExecuteReturnId(ctx, com)
		time.Sleep(1 * time.Second)
		com = cs.ReadMap()[comId]
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			time.Sleep(3 * time.Second)
			e := cs.CancelCommand(comId)
			require.NoError(t, e)
			fmt.Println("---------------------------------------------------------")
			wg.Done()
		}()
		wg.Wait()
		time.Sleep(1 * time.Second)
		com = cs.ReadMap()[comId]
		fmt.Printf("%+v\n", com)
		fmt.Printf("%+v\n", string(com.RespData))
		time.Sleep(1 * time.Second)
		err := cs.CancelCommand(comId)
		require.Error(t, err)
	})
	t.Run("execute wait", func(t *testing.T) {
		h := e_command_template.HTTPSCommand{
			Method: e_command_template.GET,
			URL:    "http://192.168.1.10:9330/api/object/value/?id_list=1",
			Header: nil,
		}
		m := e_command_template.Monitor{
			StatusCode: 200,
			Interval:   1000,
			MConditions: []e_command_template.MCondition{
				{
					Order:         0,
					CalculateType: ">=",
					SearchRule:    "root.[0]array.value",
					Value:         "50",
				},
			},
		}
		com := e_command.Command{
			Token: "test",
			Template: e_command_template.CommandTemplate{
				Name:     "object_test",
				Protocol: e_command_template.Http,
				Timeout:  3000,
				Http:     &h,
				Monitor:  &m,
			},
		}
		com = cs.ExecuteWait(ctx, com)
		fmt.Printf("%+v\n", com)
		fmt.Printf("%+v\n", string(com.RespData))
	})
	t.Run("variable error", func(t *testing.T) {
		h := e_command_template.HTTPSCommand{
			Method: e_command_template.GET,
			URL:    "http://192.168.1.10:9330/api/object/value/?id_list=1",
			Header: []byte(`[{"key": "{{test}}","value": "123456","is_active": true,"data_type": "text"}]`),
		}
		m := e_command_template.Monitor{
			StatusCode: 200,
			Interval:   1000,
			MConditions: []e_command_template.MCondition{
				{
					Order:         0,
					CalculateType: ">=",
					SearchRule:    "root.[0]array.value",
					Value:         "50",
				},
			},
		}
		com := e_command.Command{
			Token: "test",
			Template: e_command_template.CommandTemplate{
				Name:     "object_test",
				Protocol: e_command_template.Http,
				Timeout:  3000,
				Http:     &h,
				Monitor:  &m,
			},
		}
		com = cs.ExecuteWait(ctx, com)
		fmt.Printf("%+v\n", com)
		fmt.Printf("%+v\n", string(com.RespData))
		require.ErrorIs(t, com.Message, &HeaderVariables)
	})
}

func TestReadHistory(t *testing.T) {
	cs, _ := setUpServer()
	ctx := context.Background()
	go func() { cs.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		hc, err := cs.ReadFromHistory("", "-50d", "", "Success")
		require.NoError(t, err)
		for _, com := range hc {
			require.Equal(t, e_command.Success, com.Status)
		}
	})
}
