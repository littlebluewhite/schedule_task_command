package main

import (
	"context"
	"fmt"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/dal/query"
	"schedule_task_command/util/logFile"
)

func main() {
	ctx := context.Background()
	DBS := dbs.NewDbs(logFile.NewLogFile("test", "sql_test"), true)
	qc := query.Use(DBS.GetSql()).Counter
	counter, err := qc.WithContext(ctx).Where(qc.Name.Eq("task")).First()
	if err != nil {
		e := qc.WithContext(ctx).Create(&model.Counter{Name: "task", Value: 0})
		if e != nil {
			panic(e)
		}
	}
	fmt.Println(counter)
	fmt.Println(err)
	_, err = qc.WithContext(ctx).Where(qc.Name.Eq("task")).Update(qc.Value, 20)
	if err != nil {
		panic(err)
	}
}
