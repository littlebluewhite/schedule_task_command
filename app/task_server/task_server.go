package task_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/util/logFile"
	"sync"
	"time"
)

type commandServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	Execute(ctx context.Context, templateId int, triggerFrom []string,
		triggerAccount string, token string) (com e_command.Command)
}

type TaskServer struct {
	dbs dbs.Dbs
	l   logFile.LogFile
	t   map[string]e_task.Task
	cs  commandServer
	chs chs
}

func NewTaskServer(dbs dbs.Dbs, cs commandServer) *TaskServer {
	l := logFile.NewLogFile("app", "task_server")
	t := make(map[string]e_task.Task)
	rec := make(chan e_task.Task)
	mu := new(sync.RWMutex)
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

func (t *TaskServer) Start(ctx context.Context, removeTime time.Duration) {
	t.l.Info().Println("Task server started")
	defer t.l.Error().Println("Task server stopped")
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func(wg *sync.WaitGroup) {
		t.removeFinishedTask(ctx, removeTime)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		t.rdbSub(ctx)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		t.cs.Start(ctx, removeTime)
		wg.Done()
	}(wg)
	wg.Wait()
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
			t.l.Error().Println("Error executing Task")
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

func (t *TaskServer) removeFinishedTask(ctx context.Context, s time.Duration) {
Loop1:
	for {
		select {
		case <-ctx.Done():
			break Loop1
		default:
			t.chs.mu.Lock()
			now := time.Now()
			for tId, item := range t.t {
				if item.Status.TStatus != e_task.Process && item.To.Add(s).After(now) {
					delete(t.t, tId)
				}
			}
			time.Sleep(s)
		}
	}
}

func (t *TaskServer) writeToHistory(task e_task.Task) {
	ctx := context.Background()
	p := influxdb2.NewPoint("task_history",
		map[string]string{"task_id": task.TaskId, "status": task.Status.TStatus.String()},
		map[string]interface{}{"data": task},
		task.From,
	)
	if err := t.dbs.GetIdb().Writer().WritePoint(ctx, p); err != nil {
		panic(err)
	}
}

func (t *TaskServer) ReadFromHistory(taskId, start, stop, status string) (ht []e_task.Task) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	statusValue := ""
	if status != "" {
		statusValue = fmt.Sprintf(`|> filter(fn: (r) => r.status == "%s"`, status)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule"
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "task_history"
|> filter(fn: (r) => r.task_id == "%s")
|> filter(fn: (r) => r."_field" == "data")
%s
`, start, stopValue, taskId, statusValue)
	result, err := t.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var task e_task.Task
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &task); e != nil {
				panic(e)
			}
			ht = append(ht, task)
		}
	} else {
		panic(err)
	}
	return
}

func (t *TaskServer) writeTask(task e_task.Task) {
	t.chs.mu.Lock()
	defer t.chs.mu.Unlock()
	t.t[task.TaskId] = task
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
