package time_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
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
		st.TriggerFrom = append(st.TriggerFrom, "redis channel")
		t.Execute(st.TemplateId, st.TriggerFrom, st.TriggerAccount, st.Token)
		if err != nil {
			t.l.Error().Println("Error executing Command")
		}
	}
}

func (t *TimeServer) Execute(templateId int, triggerFrom []string,
	triggerAccount string, token string) {
	pt := publishTime{
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
}

func (t *TimeServer) checkTime(pt publishTime) publishTime {
	timeTemplate, ok := t.getTimeTemplate()[pt.TemplateId]
	nowTime := time.Now()
	pt.Time = nowTime
	if !ok {
		pt.Status = Failure
		pt.Message = CannotFindTemplate
		return pt
	}
	pt.Status = Success
	timeData := timeTemplate.GetTimeData()
	isTime := timeData.CheckTimeData(nowTime)
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

func (t *TimeServer) writeToHistory(pt publishTime) {
	ctx := context.Background()
	templateId := fmt.Sprintf("%d", pt.TemplateId)
	p := influxdb2.NewPoint("time_history",
		map[string]string{"template_id": templateId},
		map[string]interface{}{"data": pt},
		pt.Time,
	)
	if err := t.dbs.GetIdb().Writer().WritePoint(ctx, p); err != nil {
		panic(err)
	}
}

func (t *TimeServer) ReadFromHistory(templateId, start, stop string) (ht []publishTime) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	templateIdValue := ""
	if templateId != "" {
		templateIdValue = fmt.Sprintf(`|> filter(fn: (r) => r.template_id == "%s"`, templateId)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule"
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "time_history"
|> filter(fn: (r) => r."_field" == "data")
%s
`, start, stopValue, templateIdValue)
	result, err := t.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var pt publishTime
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &pt); e != nil {
				panic(e)
			}
			ht = append(ht, pt)
		}
	} else {
		panic(err)
	}
	return
}

func (t *TimeServer) rdbPub(pt publishTime) (e error) {
	ctx := context.Background()
	trb, _ := json.Marshal(pt)
	e = t.dbs.GetRdb().Publish(ctx, "timeRec", trb).Err()
	if e != nil {
		t.l.Error().Println("redis Publish error")
		return
	}
	return
}
