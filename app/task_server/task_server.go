package task_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/util/logFile"
	"sync"
	"time"
)

type commandServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	Execute(ctx context.Context, sc e_command.SendCommand) (com e_command.Command)
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
		var s e_task.Task
		err = json.Unmarshal(b, &s)
		if err != nil {
			t.l.Error().Println("send data is not correctly")
		}
		s.TriggerFrom = append(s.TriggerFrom, "redis channel")
		_, err = t.ExecuteReturnId(ctx, s)
		if err != nil {
			t.l.Error().Println("Error executing Task")
		}
	}
}

func (t *TaskServer) ReadMap() map[string]e_task.Task {
	t.chs.mu.RLocker()
	defer t.chs.mu.RUnlock()
	return t.t
}

func (t *TaskServer) GetList() []e_task.Task {
	tl := make([]e_task.Task, 0, len(t.t))
	m := t.ReadMap()
	for _, v := range m {
		tl = append(tl, v)
	}
	return tl
}

func (t *TaskServer) ExecuteReturnId(ctx context.Context, task e_task.Task) (taskId string, err error) {
	// publish to redis
	_ = t.rdbPub(task)
	if task.Message != nil {
		err = task.Message
		t.l.Error().Println(task.Message)
		return
	}
	go func() {
		t.doTask(ctx, task)
	}()
	taskId = task.TaskId
	return
}

func (t *TaskServer) ExecuteWaitTask(ctx context.Context, task e_task.Task) e_task.Task {
	// publish to redis
	_ = t.rdbPub(task)
	if task.Message != nil {
		t.l.Error().Println(task.Message)
		return task
	}
	ch := make(chan e_task.Task)
	go func() {
		t.doTask(ctx, task)
	}()
	task = <-ch
	return task
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
	jTask, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	templateId := fmt.Sprintf("%d", task.TemplateID)
	p := influxdb2.NewPoint("task_history",
		map[string]string{"task_template_id": templateId, "status": task.Status.TStatus.String()},
		map[string]interface{}{"data": jTask},
		task.From,
	)
	if err := t.dbs.GetIdb().Writer().WritePoint(ctx, p); err != nil {
		panic(err)
	}
}

func (t *TaskServer) ReadFromHistory(taskTemplateId, status, start, stop string) (ht []e_task.Task) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	statusValue := ""
	if status != "" {
		statusValue = fmt.Sprintf(`|> filter(fn: (r) => r.status == "%s")`, status)
	}
	taskTemplateValue := ""
	if taskTemplateId != "" {
		taskTemplateValue = fmt.Sprintf(`|> filter(fn: (r) => r.task_template_id == "%s")`, taskTemplateId)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "task_history")
|> filter(fn: (r) => r._field == "data")
%s
%s
`, start, stopValue, taskTemplateValue, statusValue)
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
