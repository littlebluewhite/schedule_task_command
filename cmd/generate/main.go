package main

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
	"schedule_task_command/app/dbs/sql"
)

func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath: "./dal/query",
		/* Mode: gen.WithoutContext,*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable:  true,
		FieldCoverable: true,
	})

	db, err := sql.NewDB("mySQL", "gen_sql.log", "db")
	if err != nil {
		panic(err)
	}

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary, or it will panic
	g.UseDB(db)
	timeData := g.GenerateModel("time_data",
		gen.FieldType("t_condition", "json.RawMessage"),
		gen.FieldType("start_time", "[]byte"),
		gen.FieldType("end_time", "[]byte"))
	timeTemplate := g.GenerateModel("time_template", gen.FieldRelate(field.BelongsTo, "TimeData", timeData,
		&field.RelateConfig{
			GORMTag: map[string][]string{"foreignKey": {"time_data_id"}},
		}))
	headerTemplate := g.GenerateModel("header_template", gen.FieldType("data", "json.RawMessage"))
	httpsCommand := g.GenerateModel("https_command",
		gen.FieldType("header", "json.RawMessage"),
		gen.FieldType("params", "json.RawMessage"),
		gen.FieldType("body", "*json.RawMessage"))
	websocketCommand := g.GenerateModel("websocket_command",
		gen.FieldType("header", "json.RawMessage"))
	mqttCommand := g.GenerateModel("mqtt_command",
		gen.FieldType("header", "json.RawMessage"),
		gen.FieldType("message", "*json.RawMessage"))
	redisCommand := g.GenerateModel("redis_command",
		gen.FieldType("message", "*json.RawMessage"))
	mCondition := g.GenerateModel("m_condition")
	monitor := g.GenerateModel("monitor", gen.FieldRelate(field.HasMany, "MConditions",
		mCondition, &field.RelateConfig{
			GORMTag: map[string][]string{"foreignKey": {"monitor_id"}},
		}))
	commandTemplate := g.GenerateModel("command_template",
		gen.FieldRelate(field.HasOne, "Http", httpsCommand, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"command_template_id"}},
			RelatePointer: true,
		}),
		gen.FieldRelate(field.HasOne, "Mqtt", mqttCommand, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"command_template_id"}},
			RelatePointer: true,
		}),
		gen.FieldRelate(field.HasOne, "Websocket", websocketCommand, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"command_template_id"}},
			RelatePointer: true,
		}),
		gen.FieldRelate(field.HasOne, "Redis", redisCommand, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"command_template_id"}},
			RelatePointer: true,
		}),
		gen.FieldRelate(field.HasOne, "Monitor", monitor, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"command_template_id"}},
			RelatePointer: true,
		}),
		gen.FieldType("tags", "json.RawMessage"),
	)
	taskStage := g.GenerateModel("task_stage",
		gen.FieldRelate(field.BelongsTo, "CommandTemplate", commandTemplate, &field.RelateConfig{
			GORMTag:       map[string][]string{"foreignKey": {"command_template_id"}},
			RelatePointer: true,
		}),
		gen.FieldType("tags", "json.RawMessage"))
	taskTemplate := g.GenerateModel("task_template",
		gen.FieldRelate(field.Many2Many, "Stages", taskStage, &field.RelateConfig{
			GORMTag: map[string][]string{"many2many": {"task_template_stage"}},
		}),
		gen.FieldType("variable", "json.RawMessage"),
		gen.FieldType("tags", "json.RawMessage"),
	)

	taskTemplateStage := g.GenerateModel("task_template_stage")

	schedule := g.GenerateModel("schedule",
		gen.FieldRelate(field.BelongsTo, "TimeData", timeData,
			&field.RelateConfig{
				GORMTag: map[string][]string{"foreignKey": {"time_data_id"}},
			}),
		gen.FieldType("tags", "json.RawMessage"))

	g.ApplyBasic(timeData, timeTemplate, headerTemplate, httpsCommand, commandTemplate,
		redisCommand, mqttCommand, websocketCommand, monitor, mCondition, taskTemplateStage,
		taskStage, taskTemplate, schedule)

	// execute the action of code generation
	g.Execute()
}
