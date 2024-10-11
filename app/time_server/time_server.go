package time_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/api"
	"schedule_task_command/app/dbs"
	"schedule_task_command/entry/e_time"
	"schedule_task_command/util/my_log"
	"sync"
)

type TimeServer struct {
	dbs dbs.Dbs
	l   api.Logger
	mu  *sync.RWMutex
}

func NewTimeServer(dbs dbs.Dbs) *TimeServer {
	l := my_log.NewLog("app/time_server")
	mu := new(sync.RWMutex)
	return &TimeServer{
		dbs: dbs,
		l:   l,
		mu:  mu,
	}
}

func (t *TimeServer) Start(ctx context.Context) {
	t.l.Infoln("Time server started")
	defer t.l.Errorln("Time server stopped")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		t.rdbSub(ctx)
		wg.Done()
	}(wg)
	wg.Wait()
}

func (t *TimeServer) rdbSub(ctx context.Context) {
	pubsub := t.dbs.GetRdb().Subscribe(ctx, "sendTime")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			t.l.Errorln(err)
		}
		b := []byte(msg.Payload)
		var pt e_time.PublishTime
		err = json.Unmarshal(b, &pt)
		if err != nil {
			t.l.Errorln("Error executing Command")
		}
		_, _ = t.Execute(pt)
	}
}

func (t *TimeServer) Execute(pt e_time.PublishTime) (bool, error) {
	// check time
	pt = t.checkTime(pt)

	// write to history
	go func() {
		t.writeToHistory(pt)
	}()

	// send to redis channel
	//_ = t.rdbPub(pt)

	if pt.Message != nil {
		return false, pt.Message
	}

	return pt.IsTime, nil
}

func (t *TimeServer) checkTime(pt e_time.PublishTime) e_time.PublishTime {
	if pt.Message != nil {
		return pt
	}
	pt.Status = e_time.Success
	pt.IsTime = pt.TimeData.CheckTimeData(pt.Time)
	return pt
}

func (t *TimeServer) writeToHistory(pt e_time.PublishTime) {
	templateId := fmt.Sprintf("%d", pt.TemplateId)
	isTime := fmt.Sprintf("%t", pt.IsTime)
	jsonPt, err := json.Marshal(pt)
	if err != nil {
		t.l.Errorln(err)
	}
	p := influxdb2.NewPoint("time_history",
		map[string]string{"template_id": templateId, "is_time": isTime},
		map[string]interface{}{"data": jsonPt},
		pt.Time,
	)
	t.dbs.GetIdb().Writer().WritePoint(p)
}

func (t *TimeServer) ReadFromHistory(id, templateId, start, stop, isTime string) (ht []e_time.PublishTime, err error) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	templateIdValue := ""
	if templateId != "" {
		templateIdValue = fmt.Sprintf(`|> filter(fn: (r) => r.template_id == "%s")`, templateId)
	}
	isTimeValue := ""
	if isTime != "" {
		isTimeValue = fmt.Sprintf(`|> filter(fn: (r) => r.is_time == "%s")`, isTime)
	}
	timeIDValue := ""
	if id != "" {
		timeIDValue = fmt.Sprintf(`|> filter(fn: (r) => r.id == "%s")`, id)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "time_history")
|> filter(fn: (r) => r._field == "data")
%s
%s
%s
`, start, stopValue, templateIdValue, isTimeValue, timeIDValue)
	result, err := t.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var pt e_time.PublishTime
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &pt); e != nil {
				t.l.Errorln(e)
			}
			ht = append(ht, pt)
		}
	} else {
		return nil, err
	}
	return
}

func (t *TimeServer) rdbPub(pt e_time.PublishTime) (e error) {
	ctx := context.Background()
	trb, _ := json.Marshal(pt)
	e = t.dbs.GetRdb().Publish(ctx, "timeRec", trb).Err()
	if e != nil {
		t.l.Errorln("redis Publish error")
		return
	}
	return
}
