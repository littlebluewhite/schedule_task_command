package task_server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"schedule_task_command/api"
	"schedule_task_command/app/command_server"
	"schedule_task_command/app/dbs"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
	"sync"
	"testing"
	"time"
)

func setUpServer() (ts *TaskServer[api.CommandServer]) {
	l := logFile.NewLogFile("test", "taskServer.log")
	DBS := dbs.NewDbs(l, true)
	cs := command_server.NewCommandServer(DBS)
	ts = NewTaskServer[api.CommandServer](DBS, cs)
	return
}

func TestExecuteReturnId(t *testing.T) {
	ts := setUpServer()
	ctx := context.Background()
	go func() { ts.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		task := e_task.Task{
			Token: "test1",
		}
		id, err := ts.ExecuteReturnId(ctx, task)
		require.NotEqual(t, id, "")
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		fmt.Println(ts.GetList())
	})
	t.Run("test2", func(t *testing.T) {
		e := util.MyErr("test err")
		task := e_task.Task{
			Token:   "test2",
			Message: &e,
		}
		id, err := ts.ExecuteReturnId(ctx, task)
		require.Error(t, err)
		require.Equal(t, id, "")
		time.Sleep(1 * time.Second)
		fmt.Println(ts.GetList())
	})
}

func TestReadTask(t *testing.T) {
	ts := setUpServer()
	ctx := context.Background()
	go func() { ts.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		tm := ts.ReadMap()
		tl := ts.GetList()

		fmt.Println(tl)
		fmt.Println(tm)
	})
}

func TestDoTask(t *testing.T) {
	ts := setUpServer()
	ctx := context.Background()
	go func() { ts.Start(ctx, 2*time.Minute) }()
	t.Run("test success", func(t *testing.T) {
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
		h2 := e_command_template.HTTPSCommand{
			Method:   e_command_template.PUT,
			URL:      "http://192.168.1.10:9330/api/object/{{iv}}/",
			Header:   []byte(`[{"key": "test","value": "123456","is_active": true,"data_type": "text"}]`),
			Body:     []byte(`[{"id": 1,"value": "{{value}}"}]`),
			BodyType: e_command_template.Json,
		}
		h3 := e_command_template.HTTPSCommand{
			Method:   e_command_template.PUT,
			URL:      "http://192.168.1.10:9330/api/object/{{iv}}/",
			Header:   []byte(`[{"key": "test","value": "123456","is_active": true,"data_type": "text"}]`),
			Body:     []byte(`[{"id": 1,"value": "{{value}}"}]`),
			BodyType: e_command_template.Json,
		}
		t1 := e_task.Task{
			Token: "test1",
			Variables: map[int]map[string]string{
				1: {"value": "1", "iv": "insert_value"},
				2: {"value": "2", "iv": "insert_value"},
			},
			TaskData: e_task_template.TaskTemplate{
				Name: "test1_name",
				StageItems: []e_task_template.StageItem{
					{
						ID:          3,
						Name:        "c3",
						StageNumber: 2,
						Mode:        e_task_template.Monitor,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_get_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h1,
							Monitor:  &m1,
						},
					},
					{
						ID:          2,
						Name:        "c2",
						StageNumber: 2,
						Mode:        e_task_template.Execute,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_put_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h2,
						},
					},
					{
						ID:          1,
						Name:        "c1",
						StageNumber: 1,
						Mode:        e_task_template.Execute,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_put_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h3,
						},
					},
				},
			},
		}
		id, _ := ts.ExecuteReturnId(ctx, t1)
		time.Sleep(3 * time.Second)
		task := ts.ReadMap()[id]
		fmt.Printf("task: %+v\n", task)
		cm := ts.GetCommandServer().ReadMap()
		fmt.Printf("command: %+v\n", cm)
	})
	t.Run("task failure", func(t *testing.T) {
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
		h3 := e_command_template.HTTPSCommand{
			Method:   e_command_template.PUT,
			URL:      "http://192.168.1.10:9330/api/object/insert_value/",
			Header:   []byte(`[{"key": "test","value": "123456","is_active": true,"data_type": "text"}]`),
			Body:     []byte(`[{"id": 1,"value": "1"}]`),
			BodyType: e_command_template.Json,
		}
		t1 := e_task.Task{
			Token: "test2",
			TaskData: e_task_template.TaskTemplate{
				Name: "test2_name",
				StageItems: []e_task_template.StageItem{
					{
						Name:        "c3",
						StageNumber: 2,
						Mode:        e_task_template.Monitor,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_get_test",
							Protocol: e_command_template.Http,
							Timeout:  5000,
							Http:     &h1,
							Monitor:  &m1,
						},
					},
					{
						Name:        "c1",
						StageNumber: 1,
						Mode:        e_task_template.Execute,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_put_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h3,
						},
					},
				},
			},
		}
		task := ts.ExecuteWait(ctx, t1)
		fmt.Printf("task: %+v\n", task)
		cm := ts.GetCommandServer().ReadMap()
		fmt.Printf("command: %+v\n", cm)
	})
	t.Run("task cancel", func(t *testing.T) {
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
		h3 := e_command_template.HTTPSCommand{
			Method:   e_command_template.PUT,
			URL:      "http://192.168.1.10:9330/api/object/insert_value/",
			Header:   []byte(`[{"key": "test","value": "123456","is_active": true,"data_type": "text"}]`),
			Body:     []byte(`[{"id": 1,"value": "1"}]`),
			BodyType: e_command_template.Json,
		}
		t1 := e_task.Task{
			Token: "test2",
			TaskData: e_task_template.TaskTemplate{
				Name: "test2_name",
				StageItems: []e_task_template.StageItem{
					{
						Name:        "c3",
						StageNumber: 2,
						Mode:        e_task_template.Monitor,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_get_test",
							Protocol: e_command_template.Http,
							Timeout:  5000,
							Http:     &h1,
							Monitor:  &m1,
						},
					},
					{
						Name:        "c1",
						StageNumber: 1,
						Mode:        e_task_template.Execute,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_put_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h3,
						},
					},
				},
			},
		}
		id, _ := ts.ExecuteReturnId(ctx, t1)
		time.Sleep(1 * time.Second)
		task := ts.ReadMap()[id]
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			time.Sleep(2 * time.Second)
			e := ts.CancelTask(id, "test")
			require.NoError(t, e)
			fmt.Println("---------------------------------------------------------")
			wg.Done()
		}()
		wg.Wait()
		time.Sleep(1 * time.Second)
		task = ts.ReadMap()[id]
		comM := ts.GetCommandServer().ReadMap()
		fmt.Printf("tasks: %+v\n", task)
		fmt.Printf("coms: %+v\n", comM)
		require.Equal(t, e_task.Cancel, task.Status)
	})
	t.Run("test variable error", func(t *testing.T) {
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
		h2 := e_command_template.HTTPSCommand{
			Method:   e_command_template.PUT,
			URL:      "http://192.168.1.10:9330/api/object/{{iv}}/",
			Header:   []byte(`[{"key": "test","value": "123456","is_active": true,"data_type": "text"}]`),
			Body:     []byte(`[{"id": 1,"value": "{{value}}"}]`),
			BodyType: e_command_template.Json,
		}
		h3 := e_command_template.HTTPSCommand{
			Method:   e_command_template.PUT,
			URL:      "http://192.168.1.10:9330/api/object/{{iv}}/",
			Header:   []byte(`[{"key": "test","value": "123456","is_active": true,"data_type": "text"}]`),
			Body:     []byte(`[{"id": 1,"value": "{{value}}"}]`),
			BodyType: e_command_template.Json,
		}
		t1 := e_task.Task{
			Token: "test1",
			Variables: map[int]map[string]string{
				1: {"value": "1", "iv": "insert_value"},
				2: {"value": "2"},
			},
			TaskData: e_task_template.TaskTemplate{
				Name: "test1_name",
				StageItems: []e_task_template.StageItem{
					{
						Name:        "c3",
						StageNumber: 2,
						Mode:        e_task_template.Monitor,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_get_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h1,
							Monitor:  &m1,
						},
					},
					{
						ID:          2,
						Name:        "c2",
						StageNumber: 2,
						Mode:        e_task_template.Execute,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_put_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h2,
						},
					},
					{
						ID:          1,
						Name:        "c1",
						StageNumber: 1,
						Mode:        e_task_template.Execute,
						CommandTemplate: e_command_template.CommandTemplate{
							Name:     "object_put_test",
							Protocol: e_command_template.Http,
							Timeout:  10000,
							Http:     &h3,
						},
					},
				},
			},
		}
		id, _ := ts.ExecuteReturnId(ctx, t1)
		time.Sleep(3 * time.Second)
		task := ts.ReadMap()[id]
		fmt.Printf("task: %+v\n", task)
		cm := ts.GetCommandServer().ReadMap()
		fmt.Printf("command: %+v\n", cm)
	})
}

func TestReadHistory(t *testing.T) {
	ts := setUpServer()
	ctx := context.Background()
	go func() { ts.Start(ctx, 2*time.Minute) }()
	t.Run("test1", func(t *testing.T) {
		hc, err := ts.ReadFromHistory("", "-50d", "", "Failure")
		require.NoError(t, err)
		for _, task := range hc {
			require.Equal(t, e_task.Failure, task.Status)
		}
		fmt.Printf("%+v\n", hc)
	})
}
