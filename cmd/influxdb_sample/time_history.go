package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"schedule_task_command/app/dbs/influxdb"
	"schedule_task_command/entry/e_time"
)

func main() {
	ctx := context.Background()
	idb := influxdb.NewInfluxdb("influxdb")
	defer idb.Close()
	ht := make([]e_time.PublishTime, 0, 20)
	start := "-8d"
	stopValue := ""
	//if stop != "" {
	//	stopValue = fmt.Sprintf(", stop: %s", stop)
	//}
	//templateIdValue := ""
	templateIdValue := fmt.Sprintf(`|> filter(fn: (r) => r.template_id == "%s")`, "1")
	//}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "time_history")
|> filter(fn: (r) => r._field == "data")
%s
`, start, stopValue, templateIdValue)
	result, err := idb.Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var pt e_time.PublishTime
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &pt); e != nil {
				panic(e)
			}
			ht = append(ht, pt)
		}
	} else {
		panic(err)
	}
	fmt.Println(ht)
}
