package task_server

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"strconv"
)

func commandReturn2Variables(variables map[int]map[string]string,
	comB comBuilder) map[int]map[string]string {
	for _, parserItem := range comB.parser {
		for _, to := range parserItem.To {
			variables[to.ID][to.Key] = comB.com.Return[parserItem.FromKey]
		}
	}
	return variables
}

func getStageVariables(stage e_task_template.StageItem, task e_task.Task) (v map[string]string, err error) {
	v = make(map[string]string)
	// get stage variables
	if stage.Variable != nil && string(stage.Variable) != "null" {
		if err = json.Unmarshal(stage.Variable, &v); err != nil {
			err = util.MyErr("stage variables error")
		}
	}
	// get task variables and cover
	if tv, ok := task.Variables[int(stage.ID)]; ok {
		for k, value := range tv {
			v[k] = value
		}
	}

	// get global variables
	if tv, ok := task.Variables[-1]; ok {
		for k, value := range tv {
			v[k] = value
		}
	}

	// set global default variables
	v["__task_id__"] = strconv.FormatUint(task.ID, 10)
	return
}
