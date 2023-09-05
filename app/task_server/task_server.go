package task_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"schedule_task_command/app/command_server"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/util/logFile"
	"sync"
	"time"
)

type CommandServer interface {
	Start(removeTime time.Duration)
	Execute(ctx context.Context, templateId int, triggerFrom []string,
		triggerAccount string, token string) (com e_command.Command, err error)
}

type TaskServer struct {
	dbs dbs.Dbs
	l   logFile.LogFile
	t   map[string]e_task.Task
	cs  CommandServer
	chs chs
}

func NewTaskServer(dbs dbs.Dbs) *TaskServer {
	l := logFile.NewLogFile("app", "task_server")
	t := make(map[string]e_task.Task)
	rec := make(chan e_task.Task)
	mu := new(sync.RWMutex)
	cs := command_server.NewCommandServer(dbs)
	return &TaskServer{
		dbs: dbs,
		l:   l,
		t:   t,
		cs:  cs,
		chs: chs{
			rec: rec,
			mu:  mu,
		},
	}
}

func (t *TaskServer) Start(removeTime time.Duration) {
	t.cs.Start(removeTime)
}

func (t *TaskServer) Execute(ep executeParams) (taskId string, err error) {
	task, err := t.generateTask(ep)
	// publish to redis
	_ = t.rdbPub(task)
	if err != nil {
		t.l.Error().Println(err)
		return
	}
	return
}

func (t *TaskServer) generateTask(ep executeParams) (task e_task.Task, err error) {
	cache := t.dbs.GetCache()
	var cacheMap map[int]model.TaskTemplate
	if x, found := cache.Get("taskTemplates"); found {
		cacheMap = x.(map[int]model.TaskTemplate)
	}
	tt, ok := cacheMap[ep.templateId]
	if !ok {
		err = cannotFindTemplate
		task = e_task.Task{Token: ep.token, Message: "can not find task template",
			Status: e_task.Status{TStatus: e_task.Failure}}
		return
	}
	from := time.Now()
	taskId := fmt.Sprintf("%v_%v_%v", ep.templateId, tt.Name, from.UnixMicro())
	task = e_task.Task{
		TaskId:         taskId,
		Token:          ep.token,
		From:           from,
		TriggerFrom:    ep.triggerFrom,
		TriggerAccount: ep.triggerAccount,
		TemplateID:     ep.templateId,
	}
	return
}

func (t *TaskServer) rdbPub(task e_task.Task) (e error) {
	ctx := context.Background()
	trb, _ := json.Marshal(e_task.ToPub(task))
	e = t.dbs.GetRdb().Publish(ctx, "taskRec", trb).Err()
	if e != nil {
		t.l.Error().Println("redis publish error")
		return
	}
	return
}

func (t *TaskServer) aaa() {
	fmt.Println("aaa")
}
