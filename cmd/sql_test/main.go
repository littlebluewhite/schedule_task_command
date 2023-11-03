package main

import (
	"context"
	"fmt"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/query"
	"schedule_task_command/util/logFile"
)

func main() {
	ctx := context.Background()
	DBS := dbs.NewDbs(logFile.NewLogFile("test", "sql_test"), true)
	qtt := query.Use(DBS.GetSql()).TaskTemplate
	qts := query.Use(DBS.GetSql()).TaskStage
	qtts := query.Use(DBS.GetSql()).TaskTemplateStage
	ts, err := qts.WithContext(ctx).Select(qts.ID).Where(qts.CommandTemplateID.Eq(18)).Find()
	if err != nil {
		panic(err)
	}
	IDs := make([]int32, 0, 20)
	for _, item := range ts {
		fmt.Println(item)
		IDs = append(IDs, item.ID)
	}
	tts, err := qtts.WithContext(ctx).Where(qtts.TaskStageID.In(IDs...)).Find()
	if err != nil {
		panic(err)
	}
	IDs = make([]int32, 0, 20)
	for _, item := range tts {
		fmt.Println(item)
		IDs = append(IDs, item.TaskTemplateID)
	}
	tt, err := qtt.WithContext(ctx).Where(qtt.ID.In(IDs...)).Find()
	if err != nil {
		panic(err)
	}
	fmt.Println(tt)
}
