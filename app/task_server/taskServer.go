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
		triggerAccount string, token string) (com e_command.Command)
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

func (t *TaskServer) rdbSub(ctx context.Context) {
	pubsub := t.dbs.GetRdb().Subscribe(ctx, "sendTask")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		b := []byte(msg.Payload)
		var s SendTask
		err = json.Unmarshal(b, &s)
		s.TriggerFrom = append(s.TriggerFrom, "redis channel")
		_, err = t.execute(s)
		if err != nil {
			t.l.Error().Println("Error executing Command")
		}
	}
}

func (t *TaskServer) execute(sc SendTask) (taskId string, err error) {
	ctx := context.Background()
	task := t.generateTask(sc)
	// publish to redis
	_ = t.rdbPub(task)
	if task.Message != nil {
		t.l.Error().Println(task.Message)
		return
	}
	go func() {
		t.doTask(ctx, task)
	}()
	taskId = task.TaskId
	return
}

func (t *TaskServer) Execute(ctx context.Context, templateId int, triggerFrom []string,
	triggerAccount string, token string) (taskId string) {
	sc := SendTask{
		TemplateId:     templateId,
		TriggerFrom:    triggerFrom,
		TriggerAccount: triggerAccount,
		Token:          token,
	}
	task := t.generateTask(sc)
	// publish to redis
	_ = t.rdbPub(task)
	if task.Message != nil {
		t.l.Error().Println(task.Message)
		return
	}
	ch := make(chan e_task.Task)
	go func() {
		t.doTask(ctx, task)
	}()
	task = <-ch
	return
}

func (t *TaskServer) generateTask(sc SendTask) (task e_task.Task) {
	cache := t.dbs.GetCache()
	var cacheMap map[int]model.TaskTemplate
	if x, found := cache.Get("taskTemplates"); found {
		cacheMap = x.(map[int]model.TaskTemplate)
	}
	tt, ok := cacheMap[sc.TemplateId]
	if !ok {
		task = e_task.Task{Token: sc.Token, Message: CannotFindTemplate,
			Status: e_task.Status{TStatus: e_task.Failure}}
		return
	}
	from := time.Now()
	taskId := fmt.Sprintf("%v_%v_%v", sc.TemplateId, tt.Name, from.UnixMicro())
	task = e_task.Task{
		TaskId:         taskId,
		Token:          sc.Token,
		From:           from,
		TriggerFrom:    sc.TriggerFrom,
		TriggerAccount: sc.TriggerAccount,
		TemplateID:     sc.TemplateId,
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
