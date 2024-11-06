package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/littlebluewhite/schedule_task_command/app/dbs/influxdb"
	"github.com/littlebluewhite/schedule_task_command/dal/model"
	"github.com/littlebluewhite/schedule_task_command/util/config"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
	"time"
)

func main() {
	d := model.CommandTemplate{ID: 1, Name: "aaa"}
	j, _ := json.Marshal(d)
	ctx := context.Background()
	influxdbConfig := config.InfluxdbConfig{
		Host:   "127.0.0.1",
		Port:   "8086",
		Org:    "my-org",
		Token:  "my-super-influxdb-auth-token",
		Bucket: "schedule",
	}
	idb := influxdb.NewInfluxdb(influxdbConfig, my_log.NewLog("test/influxdb_sample"))
	defer idb.Close()
	p := influxdb2.NewPoint("schedule_history",
		map[string]string{"id": "2", "name": "alarm SOP", "user": "wilson"},
		map[string]interface{}{"data": j},
		time.Now())
	p2 := influxdb2.NewPoint("schedule_history",
		map[string]string{"id": "1", "name": "alarm SOP", "user": "wilson"},
		map[string]interface{}{"complete": 1, "duration": 2},
		time.Now())
	idb.Writer().WritePoint(p)

	idb.Writer().WritePoint(p2)

	result, err := idb.Querier().Query(ctx, `from(bucket:"schedule")
|> range(start: -2h)
|> filter(fn: (r) => r._measurement == "schedule_history")`)
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			v := result.Record().Value()
			var c model.CommandTemplate
			if result.Record().Field() == "data" {
				_ = json.Unmarshal([]byte(v.(string)), &c)
				fmt.Println(c)
			}
			fmt.Printf("value: %v\ntype: %T\n", v, v)
			fmt.Printf("values: %v\n", result.Record().Values())
			fmt.Printf("result: %v\n", result.Record().Result())
			fmt.Printf("measurement: %v\n", result.Record().Measurement())
			fmt.Printf("field: %v\n", result.Record().Field())
			fmt.Printf("table: %v\n", result.Record().Table())
			fmt.Printf("start: %v\n", result.Record().Start())
			fmt.Printf("stop: %v\n", result.Record().Stop())
			fmt.Printf("time: %v\n", result.Record().Time())
			fmt.Printf("value by key(id): %v\n", result.Record().ValueByKey("id"))
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
}
