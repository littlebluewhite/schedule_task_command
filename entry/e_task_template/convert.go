package e_task_template

import (
	"fmt"
	"math"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_command_template"
)

func Format(ct []model.TaskTemplate) []TaskTemplate {
	result := make([]TaskTemplate, 0, len(ct))
	for _, item := range ct {
		fmt.Printf("%+v\n", item)
		sResult := make([]TaskStage, 0, len(item.Stages))
		for _, s := range item.Stages {
			var cTemplate *e_command_template.CommandTemplate
			if s.CommandTemplate != nil {
				cTemplate = &e_command_template.Format([]model.CommandTemplate{*s.CommandTemplate})[0]
			} else {
				cTemplate = nil
			}
			i := TaskStage{
				ID:                s.ID,
				Name:              s.Name,
				StageNumber:       s.StageNumber,
				Mode:              s.Mode,
				CommandTemplateID: s.CommandTemplateID,
				Tags:              s.Tags,
				CommandTemplate:   cTemplate,
			}
			sResult = append(sResult, i)
		}
		i := TaskTemplate{
			ID:        item.ID,
			Name:      item.Name,
			Variable:  item.Variable,
			UpdatedAt: item.UpdatedAt,
			CreatedAt: item.CreatedAt,
			Stages:    sResult,
			Tags:      item.Tags,
		}
		result = append(result, i)
	}
	return result
}

func CreateConvert(c []*TaskTemplateCreate) []*model.TaskTemplate {
	result := make([]*model.TaskTemplate, 0, len(c))
	for _, item := range c {
		sResult := make([]model.TaskStage, 0, len(item.Stages))
		for _, s := range item.Stages {
			i := model.TaskStage{
				Name:              s.Name,
				StageNumber:       s.StageNumber,
				Mode:              s.Mode,
				CommandTemplateID: s.CommandTemplateID,
				Tags:              s.Tags,
			}
			sResult = append(sResult, i)
		}
		i := model.TaskTemplate{
			Name:     item.Name,
			Variable: item.Variable,
			Stages:   sResult,
			Tags:     item.Tags,
		}
		result = append(result, &i)
	}
	return result
}

func UpdateConvert(ttMap map[int]model.TaskTemplate, utt []*TaskTemplateUpdate) (result []*model.TaskTemplate, err error) {
	for _, u := range utt {
		tt, ok := ttMap[int(u.ID)]
		if !ok {
			err = TaskTemplateNotFound(int(u.ID))
			return
		}
		if u.Name != nil {
			tt.Name = *u.Name
		}
		if u.Variable != nil {
			tt.Variable = *u.Variable
		}
		if u.Tags != nil {
			tt.Tags = *u.Tags
		}
		sId := make(map[int32]struct{})
		for _, s := range tt.Stages {
			sId[s.ID] = struct{}{}
		}
		if u.Stages != nil {
			sResult := make([]model.TaskStage, 0, len(u.Stages))
			for _, s := range u.Stages {
				_, ok := sId[int32(math.Abs(float64(s.ID)))]
				if !ok && s.ID != 0 {
					continue
				}
				ts := model.TaskStage{
					ID:                s.ID,
					Name:              s.Name,
					StageNumber:       s.StageNumber,
					Mode:              s.Mode,
					CommandTemplateID: s.CommandTemplateID,
					Tags:              s.Tags,
				}
				sResult = append(sResult, ts)
			}
			tt.Stages = sResult
		}
		result = append(result, &tt)
	}
	return
}
