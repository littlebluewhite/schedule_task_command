package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/app/dbs/influxdb"
	"github.com/littlebluewhite/schedule_task_command/entry/e_time"
	"github.com/littlebluewhite/schedule_task_command/util/config"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
)

func main() {
	ctx := context.Background()
	influxdbConfig := config.InfluxdbConfig{
		Host:   "127.0.0.1",
		Port:   "8086",
		Org:    "my-org",
		Token:  "my-super-influxdb-auth-token",
		Bucket: "schedule",
	}
	idb := influxdb.NewInfluxdb(influxdbConfig, my_log.NewLog("test/influxdb"))
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
