package task_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/api"
	"schedule_task_command/dal/model"
	"schedule_task_command/dal/query"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_module"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/util/my_log"
	"schedule_task_command/util/redis_stream"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type commandServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	ExecuteWait(ctx context.Context, com e_command.Command) e_command.Command
	Close()
}

type TaskServer[T any] struct {
	dbs          api.Dbs
	hm           hubManager
	l            api.Logger
	t            sync.Map
	cs           commandServer
	streamComMap map[string]func(rsc map[string]interface{}) (string, error)
	count        atomic.Uint64
	chs          chs
}

func NewTaskServer[T any](dbs api.Dbs, cs commandServer, wm hubManager) *TaskServer[T] {
	l := my_log.NewLog("app/task_server")
	return &TaskServer[T]{
		dbs: dbs,
		hm:  wm,
		l:   l,
		t:   sync.Map{},
		cs:  cs,
		chs: chs{
			mu: new(sync.RWMutex),
			wg: new(sync.WaitGroup),
		},
	}
}

func (t *TaskServer[T]) Start(ctx context.Context, removeTime time.Duration) {
	t.initialCounter(ctx)
	// stream command initial
	t.initStreamComMap()
	t.l.Infoln("Task server started")

	t.chs.wg.Add(1)
	go func() {
		defer t.chs.wg.Done()
		t.removeFinishedTask(ctx, removeTime)
	}()

	t.chs.wg.Add(1)
	go func() {
		defer t.chs.wg.Done()
		t.rdbSub(ctx)
	}()
	t.chs.wg.Add(1)
	go func() {
		defer t.chs.wg.Done()
		t.receiveStream(ctx)
	}()

	t.chs.wg.Add(1)
	go func() {
		defer t.chs.wg.Done()
		t.cs.Start(ctx, removeTime)
	}()

	t.chs.wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		defer t.chs.wg.Done()
		for {
			select {
			case <-ctx.Done():
				t.counterWrite()
				return
			case <-ticker.C:
				t.counterWrite()
			}
		}
	}()
}

func (t *TaskServer[T]) Close() {
	t.cs.Close()
	t.chs.wg.Wait()
	t.l.Infoln("task server stop gracefully")
}

func (t *TaskServer[T]) initialCounter(ctx context.Context) {
	qc := query.Use(t.dbs.GetSql()).Counter
	tc, err := qc.WithContext(ctx).Where(qc.Name.In("task")).First()
	if err != nil {
		tc = &model.Counter{Name: "task", Value: 0}
		e := qc.WithContext(ctx).Create(tc)
		if e != nil {
			t.l.Errorln(e)
		}
	}
	t.count.Store(uint64(tc.Value))
}

func (t *TaskServer[T]) initStreamComMap() {
	t.streamComMap = map[string]func(rsc map[string]interface{}) (string, error){
		"cancel_task": t.streamCancelTask,
	}
}

func (t *TaskServer[T]) counterWrite() {
	ctx := context.Background()
	qc := query.Use(t.dbs.GetSql()).Counter
	_, err := qc.WithContext(ctx).Where(qc.Name.Eq("task")).Update(qc.Value, t.count.Load())
	if err != nil {
		t.l.Errorln(err)
	}
}

func (t *TaskServer[T]) rdbSub(ctx context.Context) {
	pubsub := t.dbs.GetRdb().Subscribe(ctx, "sendTask")
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			t.l.Errorln(err)
		}
	}(pubsub)
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				t.l.Infoln("rdbSub receive canceled")
				return
			}
			t.l.Errorln(err)
			if errors.Is(err, redis.ErrClosed) {
				return
			}
			continue
		}

		// deal with message
		go func(payload string) {
			var s e_task.Task
			if e := json.Unmarshal([]byte(payload), &s); e != nil {
				t.l.Errorln(SendToRedisErr)
				return
			}
			if _, e := t.ExecuteReturnId(ctx, s); e != nil {
				t.l.Errorln("Error executing Task:", e)
			}
		}(msg.Payload)
	}
}

func (t *TaskServer[T]) receiveStream(ctx context.Context) {
	t.l.Infoln("----------------------------------- start task receiveStream --------------------------------")
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
	result := make(map[uint64]e_task.Task)
	t.t.Range(func(key, value interface{}) bool {
		id, ok1 := key.(uint64)
		task, ok2 := value.(e_task.Task)
		if ok1 && ok2 {
			result[id] = task
		}
		return true
	})
	return result
}

func (t *TaskServer[T]) ReadOne(id uint64) (e_task.Task, error) {
	if task, ok := t.t.Load(id); ok {
		return task.(e_task.Task), nil
	} else {
		return e_task.Task{}, errors.New(fmt.Sprintf("cannot find id: %d", id))
	}
}

func (t *TaskServer[T]) DeleteOne(id uint64) {
	t.t.Delete(id)
}

func (t *TaskServer[T]) GetList() []e_task.Task {
	result := make([]e_task.Task, 0, 1000)
	t.t.Range(func(key, value interface{}) bool {
		_, ok1 := key.(uint64)
		task, ok2 := value.(e_task.Task)
		if ok1 && ok2 {
			result = append(result, task)
		}
		return true
	})
	return result
}

func (t *TaskServer[T]) ExecuteReturnId(ctx context.Context, task e_task.Task) (id uint64, err error) {
	task.Stages = map[int32]e_task.Stage{}
	// pass the variables
	task = t.getVariables(task)
	if task.Message != nil {
		err = task.Message
		t.l.Errorln(err)
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
		t.l.Errorln(task.Message)
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
	task, err := t.ReadOne(id)
	if err != nil {
		return err
	}
	if task.Status != e_task.Process {
		return TaskCannotCancel
	}
	task.ClientMessage = message
	t.writeTask(task)
	task.CancelFunc()
	return nil
}

func (t *TaskServer[T]) removeFinishedTask(ctx context.Context, removeTime time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			t.l.Infoln("removeFinishedTask goroutine exiting")
			return
		case <-ticker.C:
			t.cleanupTasks(removeTime)
		}
	}
}

func (t *TaskServer[T]) cleanupTasks(removeTime time.Duration) {
	now := time.Now()
	t.t.Range(func(key, value interface{}) bool {
		id, ok1 := key.(uint64)
		task, ok2 := value.(e_task.Task)
		if ok1 && ok2 {
			if task.To != nil && task.Status != e_task.Process && task.To.Add(removeTime).Before(now) {
				t.t.Delete(id)
			}
		}
		return true
	})
}

func (t *TaskServer[T]) writeToHistory(task e_task.Task) {
	tp := e_task.ToPub(task)
	jTask, err := json.Marshal(tp)
	if err != nil {
		t.l.Errorf("err: %v, ToPub: %+v", err, tp)
	}
	templateId := fmt.Sprintf("%d", task.TemplateId)
	p := influxdb2.NewPoint("task_history",
		map[string]string{
			"id":               strconv.FormatUint(task.ID, 10),
			"task_template_id": templateId,
			"status":           task.Status.String()},
		map[string]interface{}{"data": jTask},
		task.From,
	)

	// write to influxdb
	t.dbs.GetIdb().Writer().WritePoint(p)
}

func (t *TaskServer[T]) ReadFromHistory(id, taskTemplateId, start, stop, status string) (ht []e_task.TaskPub, err error) {
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
	taskIDValue := ""
	if id != "" {
		taskIDValue = fmt.Sprintf(`|> filter(fn: (r) => r.id == "%s")`, id)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "task_history")
|> filter(fn: (r) => r._field == "data")
%s
%s
%s
`, start, stopValue, taskTemplateValue, statusValue, taskIDValue)
	result, err := t.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var task e_task.TaskPub
			v := result.Record().Value()
			vString, ok := v.(string)
			if !ok {
				t.l.Errorf("value: %v is not string", v)
				continue
			}
			if err = json.Unmarshal([]byte(vString), &task); err != nil {
				t.l.Errorln(err)
				continue
			}
			ht = append(ht, task)
		}
	} else {
		return
	}
	// send empty []
	if ht == nil {
		ht = make([]e_task.TaskPub, 0)
	}
	return
}

func (t *TaskServer[T]) writeTask(task e_task.Task) {
	t.t.Store(task.ID, task)
}

func (t *TaskServer[T]) publishContainer(ctx context.Context, task e_task.Task) {
	//go func() {
	//	_ = t.rdbPub(ctx, task)
	//}()
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
		t.l.Errorln("redis publish error: ", err)
		return
	}
	return
}

func (t *TaskServer[T]) StreamPub(ctx context.Context, task e_task.Task) (err error) {
	data := map[string]interface{}{
		"id":           task.ID,
		"stages":       task.StageNumber,
		"status":       task.Status,
		"message":      task.Message,
		"trigger_from": task.TriggerFrom,
	}
	timeNow := time.Now()
	jd, _ := json.Marshal(data)
	values := redis_stream.CreateRedisStreamCom()
	values["command"] = "track_task"
	values["timestamp"] = timeNow.Unix()
	values["data"] = jd
	values["is_wait_call_back"] = 0
	values["callback_token"] = fmt.Sprintf("%d_%d", task.ID, timeNow.UnixNano())
	values["send_pattern"] = "1"
	values["callback_timeout"] = 5
	values["status_code"] = 1
	values["callback_until_feed_back"] = 0

	err = redis_stream.StreamAdd(ctx, t.dbs.GetRdb(), "AlarmAPIModuleReceiver", values)
	if err != nil {
		t.l.Errorln(err)
		return
	}
	t.l.Infoln("stream publish success")
	return
}

func (t *TaskServer[T]) sendWebsocket(task e_task.Task) {
	tb, _ := json.Marshal(e_task.ToPub(task))
	if t.hm != nil {
		t.hm.Broadcast(e_module.Task, tb)
	}
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
