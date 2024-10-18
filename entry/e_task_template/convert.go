package e_task_template

import (
	"github.com/goccy/go-json"
	"math"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_command_template"
)

func Format(ct []model.TaskTemplate) []TaskTemplate {
	result := make([]TaskTemplate, 0, len(ct))
	for _, item := range ct {
		//fmt.Printf("%+v\n", item)
		sResult := make([]StageItem, 0, len(item.StageItems))
		for _, s := range item.StageItems {
			m := s.Mode
			var parser []ParserItem
			_ = json.Unmarshal(s.Parser, &parser)
			i := StageItem{
				ID:                s.ID,
				Name:              s.Name,
				StageNumber:       s.StageNumber,
				Mode:              S2Mode(&m),
				CommandTemplateID: s.CommandTemplateID,
				Tags:              s.Tags,
				Variable:          s.Variable,
				Parser:            parser,
				CommandTemplate:   e_command_template.Format([]model.CommandTemplate{s.CommandTemplate})[0],
			}
			sResult = append(sResult, i)
		}
		i := TaskTemplate{
			ID:         item.ID,
			Name:       item.Name,
			Visible:    item.Visible,
			UpdatedAt:  item.UpdatedAt,
			CreatedAt:  item.CreatedAt,
			StageItems: sResult,
			Tags:       item.Tags,
		}
		result = append(result, i)
	}
	return result
}

func CreateConvert(c []*TaskTemplateCreate) []*model.TaskTemplate {
	result := make([]*model.TaskTemplate, 0, len(c))
	for _, item := range c {
		sResult := make([]model.StageItem, 0, len(item.StageItems))
		for _, s := range item.StageItems {
			i := model.StageItem{
				Name:              s.Name,
				StageNumber:       s.StageNumber,
				Mode:              s.Mode.String(),
				CommandTemplateID: s.CommandTemplateID,
				Tags:              s.Tags,
				Variable:          s.Variable,
			}
			sResult = append(sResult, i)
		}
		i := model.TaskTemplate{
			Name:       item.Name,
			Visible:    item.Visible,
			StageItems: sResult,
			Tags:       item.Tags,
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
		if u.Visible != nil {
			tt.Visible = *u.Visible
		}
		if u.Tags != nil {
			tt.Tags = u.Tags
		}
		sId := make(map[int32]model.StageItem)
		for _, s := range tt.StageItems {
			sId[s.ID] = s
		}
		sResult := make([]model.StageItem, 0, len(u.StageItems))
		if u.StageItems != nil {
			for _, s := range u.StageItems {
				ts, ok := sId[int32(math.Abs(float64(s.ID)))]
				if !ok && s.ID != 0 {
					continue
				}
				ts.ID = s.ID
				if s.Name != nil {
					ts.Name = *s.Name
				}
				if s.StageNumber != nil {
					ts.StageNumber = *s.StageNumber
				}
				if s.Mode != NoneMode {
					ts.Mode = s.Mode.String()
				}
				if s.CommandTemplateID != nil {
					ts.CommandTemplateID = *s.CommandTemplateID
				}
				if s.Tags != nil {
					ts.Tags = s.Tags
				}
				if s.Variable != nil {
					ts.Variable = s.Variable
				}
				if s.Parser != nil {
					ts.Parser = s.Parser
				}
				sResult = append(sResult, ts)
			}
		}
		tt.StageItems = sResult
		result = append(result, &tt)
	}
	return
}

func S2Mode(s *string) Mode {
	if s == nil {
		return NoneMode
	}
	switch *s {
	case "monitor":
		return Monitor
	case "execute":
		return Execute
	default:
		return NoneMode
	}
}
