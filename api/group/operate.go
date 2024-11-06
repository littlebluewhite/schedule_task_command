package group

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	api2 "github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/entry/e_log"
	"strconv"
	"time"
)

type Operate struct {
	idb logDB
}

type logDB interface {
	Writer() api.WriteAPI
	Querier() api.QueryAPI
}

func NewOperate(dbs api2.Dbs) *Operate {
	o := &Operate{
		idb: dbs.GetIdb(),
	}
	return o
}

func (o *Operate) WriteLog(c *fiber.Ctx) {
	// write my_log
	now := time.Now()
	response := c.Response()

	var module string
	switch c.Locals("Module").(type) {
	case string:
		module = c.Locals("Module").(string)
	case nil:
		module = ""
	}

	l := e_log.Log{
		Timestamp:     float64(now.UnixMilli()) / 1000000,
		Account:       c.Get("Account"),
		ContentLength: len(response.Body()),
		Datetime:      now,
		IP:            c.IP(),
		Referer:       c.Get("Referer"),
		ApiUrl:        c.OriginalURL(),
		Method:        c.Method(),
		Module:        module,
		StatusCode:    response.StatusCode(),
		Token:         c.Get("Authorization"),
		UserAgent:     c.Get("User-Agent"),
		WebPath:       c.Get("Web-Path"),
	}
	jL, err := json.Marshal(l)
	if err != nil {
		return
	}
	p := influxdb2.NewPoint("my_log",
		map[string]string{
			"account":     l.Account,
			"ip":          l.IP,
			"method":      l.Method,
			"module":      l.Module,
			"status_code": strconv.FormatInt(int64(l.StatusCode), 10),
		},
		map[string]interface{}{"data": jL},
		now,
	)

	go func() {
		o.idb.Writer().WritePoint(p)
	}()
}

func (o *Operate) ReadLog(start, stop, account, ip, method, module, statusCode string) (logs []e_log.Log, err error) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	accountValue := ""
	if account != "" {
		accountValue = fmt.Sprintf(`|> filter(fn: (r) => r.account == "%s")`, account)
	}
	ipValue := ""
	if ip != "" {
		ipValue = fmt.Sprintf(`|> filter(fn: (r) => r.ip == "%s")`, ip)
	}
	methodValue := ""
	if method != "" {
		methodValue = fmt.Sprintf(`|> filter(fn: (r) => r.method == "%s")`, method)
	}
	moduleValue := ""
	if module != "" {
		moduleValue = fmt.Sprintf(`|> filter(fn: (r) => r.module == "%s")`, module)
	}
	statusCodeValue := ""
	if statusCode != "" {
		statusCodeValue = fmt.Sprintf(`|> filter(fn: (r) => r.status_code == "%s")`, statusCode)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "my_log")
|> filter(fn: (r) => r._field == "data")
%s
%s
%s
%s
%s
`, start, stopValue, accountValue, ipValue, methodValue, moduleValue, statusCodeValue)
	result, err := o.idb.Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var l e_log.Log
			v := result.Record().Value()
			vString, ok := v.(string)
			if !ok {
				fmt.Printf("value: %v is not string", v)
				continue
			}
			if e := json.Unmarshal([]byte(vString), &l); e != nil {
				fmt.Printf(e.Error())
				continue
			}
			logs = append(logs, l)
		}
	} else {
		return
	}
	// send empty []
	if logs == nil {
		logs = make([]e_log.Log, 0)
	}
	return
}
