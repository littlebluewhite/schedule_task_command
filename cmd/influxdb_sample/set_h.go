package main

import (
	"context"
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/app/dbs/influxdb"
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
		Bucket: "node_object",
	}
	idb := influxdb.NewInfluxdb(influxdbConfig, my_log.NewLog("test/influxdb"))
	defer idb.Close()
	//p := influxdb2.NewPoint("test",
	//	map[string]string{"id": "1", "name": "alarm SOP", "user": "wilson"},
	//	map[string]interface{}{"data": 5},
	//	time.Now().Add(time.Minute*-5))
	//p2 := influxdb2.NewPoint("test",
	//	map[string]string{"id": "1", "name": "alarm SOP", "user": "wilson"},
	//	map[string]interface{}{"data": 20},
	//	time.Now().Add(time.Minute*-3))
	//p3 := influxdb2.NewPoint("test",
	//	map[string]string{"id": "1", "name": "alarm SOP", "user": "wilson"},
	//	map[string]interface{}{"data": 10},
	//	time.Now().Add(time.Minute*-2))
	//p4 := influxdb2.NewPoint("test",
	//	map[string]string{"id": "1", "name": "alarm SOP", "user": "wilson"},
	//	map[string]interface{}{"data": 100},
	//	time.Now().Add(time.Minute*-1))
	//if err := idb.Writer().WritePoint(ctx, p); err != nil {
	//	panic(err)
	//}
	//if err := idb.Writer().WritePoint(ctx, p2); err != nil {
	//	panic(err)
	//}
	//if err := idb.Writer().WritePoint(ctx, p3); err != nil {
	//	panic(err)
	//}
	//if err := idb.Writer().WritePoint(ctx, p4); err != nil {
	//	panic(err)
	//}

	//	result, err := idb.Querier().Query(ctx, `from(bucket:"history")
	//|> range(start: -20m)
	//|> filter(fn: (r) => r._measurement == "test")
	//|> aggregateWindow(every: 2m, fn: max)
	//|> fill(usePrevious: true)
	//`)
	result, err := idb.Querier().Query(ctx, `from(bucket:"node_object")
|> range(start: -6h)
|> filter(fn:(r) => r._measurement == "object_value")
|> filter(fn:(r) => r.id == "1")
|> timeWeightedAvg(unit: 6h)
`)
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			v := result.Record().Value()
			if result.Record().Field() == "data" {
				fmt.Println(v)
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
			fmt.Println("--------------------------------------------------------")
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

}
