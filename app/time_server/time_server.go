package time_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_time_data"
	"schedule_task_command/entry/e_time_template"
	"schedule_task_command/util/logFile"
	"sync"
	"time"
)

type TimeServer struct {
	dbs dbs.Dbs
	l   logFile.LogFile
	mu  *sync.RWMutex
}

func NewTimeServer(dbs dbs.Dbs) *TimeServer {
	l := logFile.NewLogFile("app", "time_server")
	mu := new(sync.RWMutex)
	return &TimeServer{
		dbs: dbs,
		l:   l,
		mu:  mu,
	}
}

func (t *TimeServer) Start(ctx context.Context) {
	t.l.Info().Println("Time server started")
	defer t.l.Error().Println("Time server stopped")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		t.rdbSub(ctx)
		wg.Done()
	}(wg)
	wg.Wait()
}

func (t *TimeServer) rdbSub(ctx context.Context) {
	pubsub := t.dbs.GetRdb().Subscribe(ctx, "sendTask")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		b := []byte(msg.Payload)
		var st SendTime
		err = json.Unmarshal(b, &st)
		if err != nil {
			t.l.Error().Println("Error executing Command")
		}
		st.TriggerFrom = append(st.TriggerFrom, "redis channel")
		_, _ = t.Execute(st.TemplateId, st.TriggerFrom, st.TriggerAccount, st.Token)
	}
}

func (t *TimeServer) Execute(templateId int, triggerFrom []string,
	triggerAccount string, token string) (bool, error) {
	pt := e_time_data.PublishTime{
		TemplateId:     templateId,
		TriggerFrom:    triggerFrom,
		TriggerAccount: triggerAccount,
		Token:          token,
	}
	// check time
	pt = t.checkTime(pt)

	// write to history
	t.writeToHistory(pt)

	// send to redis channel
	_ = t.rdbPub(pt)

	if pt.Message != nil {
		return false, pt.Message
	}

	return pt.IsTime, nil
}

func (t *TimeServer) checkTime(pt e_time_data.PublishTime) e_time_data.PublishTime {
	timeTemplate, ok := t.getTimeTemplate()[pt.TemplateId]
	nowTime := time.Now()
	pt.Time = nowTime
	if !ok {
		pt.Status = e_time_data.Failure
		pt.Message = &CannotFindTemplate
		return pt
	}
	pt.Status = e_time_data.Success
	isTime := timeTemplate.CheckTimeData(nowTime)
	pt.IsTime = isTime
	return pt
}

func (t *TimeServer) getTimeTemplate() map[int]e_time_template.TimeTemplate {
	cacheMap := make(map[int]e_time_template.TimeTemplate)
	if x, found := t.dbs.GetCache().Get("timeTemplates"); found {
		c := x.(map[int]model.TimeTemplate)
		for key, value := range c {
			cacheMap[key] = e_time_template.Model2Entry(value)
		}
	}
	return cacheMap
}

func (t *TimeServer) writeToHistory(pt e_time_data.PublishTime) {
	ctx := context.Background()
	templateId := fmt.Sprintf("%d", pt.TemplateId)
	jsonPt, err := json.Marshal(pt)
	if err != nil {
		panic(err)
	}
	p := influxdb2.NewPoint("time_history",
		map[string]string{"template_id": templateId},
		map[string]interface{}{"data": jsonPt},
		pt.Time,
	)
	if err = t.dbs.GetIdb().Writer().WritePoint(ctx, p); err != nil {
		panic(err)
	}
}

func (t *TimeServer) ReadFromHistory(templateId, start, stop string) (ht []e_time_data.PublishTime, err error) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	templateIdValue := ""
	if templateId != "" {
		templateIdValue = fmt.Sprintf(`|> filter(fn: (r) => r.template_id == "%s")`, templateId)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "time_history")
|> filter(fn: (r) => r._field == "data")
%s
`, start, stopValue, templateIdValue)
	result, err := t.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var pt e_time_data.PublishTime
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &pt); e != nil {
				panic(e)
			}
			ht = append(ht, pt)
		}
	} else {
		return nil, err
	}
	return
}

func (t *TimeServer) rdbPub(pt e_time_data.PublishTime) (e error) {
	ctx := context.Background()
	trb, _ := json.Marshal(pt)
	e = t.dbs.GetRdb().Publish(ctx, "timeRec", trb).Err()
	if e != nil {
		t.l.Error().Println("redis Publish error")
		return
	}
	return
}
