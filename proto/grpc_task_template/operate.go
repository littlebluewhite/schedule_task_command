package grpc_task_template

import (
	"errors"
	"fmt"
	"schedule_task_command/api"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"time"
)

type Operate struct {
	d api.Dbs
}

func NewOperate(d api.Dbs) *Operate {
	return &Operate{
		d: d,
	}
}

func (o *Operate) GenerateTask(st *SendTaskTemplateRequest) (task e_task.Task) {
	t := time.Now()
	task = e_task.Task{
		TemplateId:     int(st.TemplateId),
		Source:         st.Source,
		TriggerAccount: st.TriggerAccount,
		Token:          fmt.Sprintf("task-%d-%d", st.TemplateId, t.Unix()),
	}

	// Convert TriggerFrom to the expected []map[string]string type
	triggerFrom := make([]map[string]string, 0, len(st.TriggerFrom))
	for _, tf := range st.TriggerFrom {
		triggerFrom = append(triggerFrom, tf.KeyValue)
	}
	task.TriggerFrom = triggerFrom

	// Convert Variables to the expected map[int]map[string]string type
	variables := make(map[int]map[string]string, len(st.Variables))
	for key, val := range st.Variables {
		variables[int(key)] = val.KeyValue
	}
	task.Variables = variables

	ttList, err := o.findCache([]int32{int32(st.TemplateId)})
	if err != nil {
		task.Status = e_task.Failure
		task.Message = &CannotFindTemplate
		return
	}
	trigger := map[string]string{"task_template": fmt.Sprintf("%d", st.TemplateId)}
	task.TriggerFrom = append(task.TriggerFrom, trigger)
	task.TaskData = e_task_template.Format(ttList)[0]
	return
}

func (o *Operate) findCache(ids []int32) ([]model.TaskTemplate, error) {
	tt := make([]model.TaskTemplate, 0, len(ids))
	var cacheMap map[int]model.TaskTemplate
	if x, found := o.d.GetCache().Get("taskTemplates"); found {
		cacheMap = x.(map[int]model.TaskTemplate)
	} else {
		return nil, errors.New("cache error")
	}
	for _, id := range ids {
		t, ok := cacheMap[int(id)]
		if !ok {
			return nil, fmt.Errorf("id: %v not found", id)
		}
		tt = append(tt, t)
	}
	return tt, nil
}
