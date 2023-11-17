package task_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/dal/query"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/util/logFile"
	"schedule_task_command/util/redis_stream"
	"sync"
	"sync/atomic"
	"time"
)

type commandServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	ExecuteWait(ctx context.Context, com e_command.Command) e_command.Command
}

type TaskServer[T any] struct {
	dbs          dbs.Dbs
	wm           websocketManager
	l            logFile.LogFile
	t            map[uint64]e_task.Task
	cs           commandServer
	streamComMap map[string]func(rsc map[string]interface{}) (string, error)
	count        atomic.Uint64
	chs          chs
}

func NewTaskServer[T any](dbs dbs.Dbs, cs commandServer, wm websocketManager) *TaskServer[T] {
	l := logFile.NewLogFile("app", "task_server")
	t := make(map[uint64]e_task.Task)
	mu := new(sync.RWMutex)
	return &TaskServer[T]{
		dbs: dbs,
		wm:  wm,
		l:   l,
		t:   t,
		cs:  cs,
		chs: chs{
			mu: mu,
		},
	}
}

func (t *TaskServer[T]) Start(ctx context.Context, removeTime time.Duration) {
	t.initialCounter(ctx)
	// stream command initial
	t.initStreamComMap()
	t.l.Info().Println("Task server started")
	go func() {
		t.removeFinishedTask(ctx, removeTime)
	}()
	go func() {
		t.rdbSub(ctx)
	}()
	go func() {
		t.receiveStream(ctx)
	}()
	go func() {
		t.cs.Start(ctx, removeTime)
	}()
	go func() {
		_ = <-ctx.Done()
		t.stopCounter()
		t.l.Info().Println("task server stop gracefully")
	}()
}

func (t *TaskServer[T]) initialCounter(ctx context.Context) {
	qc := query.Use(t.dbs.GetSql()).Counter
	tc, err := qc.WithContext(ctx).Where(qc.Name.In("task")).First()
	if err != nil {
		tc = &model.Counter{Name: "task", Value: 0}
		e := qc.WithContext(ctx).Create(tc)
		if e != nil {
			t.l.Error().Println(e)
		}
	}
	t.count.Store(uint64(tc.Value))
}

func (t *TaskServer[T]) initStreamComMap() {
	t.streamComMap = map[string]func(rsc map[string]interface{}) (string, error){
		"cancel_task": t.streamCancelTask,
	}
}

func (t *TaskServer[T]) stopCounter() {
	ctx := context.Background()
	qc := query.Use(t.dbs.GetSql()).Counter
	_, err := qc.WithContext(ctx).Where(qc.Name.Eq("task")).Update(qc.Value, t.count.Load())
	if err != nil {
		t.l.Error().Println(err)
	}
}

func (t *TaskServer[T]) rdbSub(ctx context.Context) {
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
			t.l.Error().Println(SendToRedisErr)
		}
		s.TriggerFrom = append(s.TriggerFrom, "redis channel")
		_, err = t.ExecuteReturnId(ctx, s)
		if err != nil {
			t.l.Error().Println("Error executing Task")
		}
	}
}

func (t *TaskServer[T]) receiveStream(ctx context.Context) {
	t.l.Info().Println("----------------------------------- start task receiveStream --------------------------------")
	rs := redis_stream.NewStreamRead(t.dbs.GetRdb(), "Task", "server", t.l)
	rs.Start(ctx, t.streamComMap)
}

func (t *TaskServer[T]) streamCancelTask(rsc map[string]interface{}) (result string, err error) {
	var entry StreamCancel
	err = json.Unmarshal([]byte(rsc["data"].(string)), &entry)
	if err != nil {
		return
	}
	err = t.CancelTask(entry.ID, entry.Message)
	if err == nil {
		result = "ok"
	}
	return
}

func (t *TaskServer[T]) ReadMap() map[uint64]e_task.Task {
	t.chs.mu.RLock()
	defer t.chs.mu.RUnlock()
	return t.t
}

func (t *TaskServer[T]) GetList() []e_task.Task {
	tl := make([]e_task.Task, 0, len(t.t))
	m := t.ReadMap()
	for _, v := range m {
		tl = append(tl, v)
	}
	return tl
}

func (t *TaskServer[T]) ExecuteReturnId(ctx context.Context, task e_task.Task) (id uint64, err error) {
	task.Stages = map[int32]e_task.Stage{}
	// pass the variables
	task = t.getVariables(task)
	if task.Message != nil {
		err = task.Message
		t.l.Error().Println(err)
		return
	}
	from := time.Now()
	task.From = from
	task.ID = t.count.Add(1)
	id = task.ID
	go func() {
		t.doTask(ctx, task)
	}()
	return
}

func (t *TaskServer[T]) ExecuteWait(ctx context.Context, task e_task.Task) e_task.Task {
	task.Stages = map[int32]e_task.Stage{}
	// pass the variables
	task = t.getVariables(task)
	if task.Message != nil {
		t.l.Error().Println(task.Message)
		return task
	}
	from := time.Now()
	task.From = from
	task.ID = t.count.Add(1)
	ch := make(chan e_task.Task)
	go func() {
		ch <- t.doTask(ctx, task)
	}()
	task = <-ch
	return task
}

func (t *TaskServer[T]) CancelTask(id uint64, message string) error {
	m := t.ReadMap()
	task, ok := m[id]
	if !ok {
		return TaskNotFind
	}
	if task.Status != e_task.Process {
		return TaskCannotCancel
	}
	task.ClientMessage = message
	t.writeTask(task)
	task.CancelFunc()
	return nil
}

func (t *TaskServer[T]) removeFinishedTask(ctx context.Context, s time.Duration) {
Loop1:
	for {
		select {
		case <-ctx.Done():
			break Loop1
		default:
			t.chs.mu.Lock()
			now := time.Now()
			for tId, item := range t.t {
				if item.Status != e_task.Process && item.To.Add(s).After(now) {
					delete(t.t, tId)
				}
			}
			t.chs.mu.Unlock()
			time.Sleep(s)
		}
	}
}

func (t *TaskServer[T]) writeToHistory(task e_task.Task) {
	ctx := context.Background()
	jTask, err := json.Marshal(e_task.ToPub(task))
	if err != nil {
		panic(err)
	}
	templateId := fmt.Sprintf("%d", task.TemplateId)
	p := influxdb2.NewPoint("task_history",
		map[string]string{"task_template_id": templateId, "status": task.Status.String()},
		map[string]interface{}{"data": jTask},
		task.From,
	)
	if err := t.dbs.GetIdb().Writer().WritePoint(ctx, p); err != nil {
		panic(err)
	}
}

func (t *TaskServer[T]) ReadFromHistory(taskTemplateId, start, stop, status string) (ht []e_task.TaskPub, err error) {
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
			var task e_task.TaskPub
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &task); e != nil {
				panic(e)
			}
			ht = append(ht, task)
		}
	} else {
		return
	}
	return
}

func (t *TaskServer[T]) writeTask(task e_task.Task) {
	t.chs.mu.Lock()
	defer t.chs.mu.Unlock()
	t.t[task.ID] = task
}

func (t *TaskServer[T]) publishContainer(ctx context.Context, task e_task.Task) {
	go func() {
		_ = t.rdbPub(ctx, task)
	}()
	go func() {
		_ = t.StreamPub(ctx, task)
	}()
	go func() {
		t.sendWebsocket(task)
	}()
}

func (t *TaskServer[T]) rdbPub(ctx context.Context, task e_task.Task) (err error) {
	trb, _ := json.Marshal(e_task.ToPub(task))
	err = t.dbs.GetRdb().Publish(ctx, "taskRec", trb).Err()
	if err != nil {
		t.l.Error().Println("redis publish error: ", err)
		return
	}
	return
}

func (t *TaskServer[T]) StreamPub(ctx context.Context, task e_task.Task) (err error) {
	data := map[string]interface{}{
		"id":      task.ID,
		"stages":  task.StageNumber,
		"status":  task.Status,
		"message": task.Message,
	}
	jd, _ := json.Marshal(data)
	values := redis_stream.CreateRedisStreamCom()
	values["command"] = "track_task"
	values["timestamp"] = time.Now().Unix()
	values["data"] = jd
	values["is_wait_call_back"] = 0
	values["callback_token"] = task.Token
	values["send_pattern"] = "1"
	values["callback_timeout"] = 5
	values["status_code"] = 1
	values["callback_until_feed_back"] = 0

	err = redis_stream.StreamAdd(ctx, t.dbs.GetRdb(), "AlarmAPIModuleReceiver", values)
	if err != nil {
		t.l.Error().Println(err)
		return
	}
	t.l.Info().Println("stream publish success")
	return
}

func (t *TaskServer[T]) sendWebsocket(task e_task.Task) {
	tb, _ := json.Marshal(e_task.ToPub(task))
	t.wm.Broadcast(2, tb)
}

func (t *TaskServer[T]) GetCommandServer() T {
	return t.cs.(T)
}

func (t *TaskServer[T]) getVariables(task e_task.Task) e_task.Task {
	if task.Variables == nil {
		v := make(map[int]map[string]string)
		task.Variables = v
		// template have variables
	}
	return task
}
