package e_header_template

import "schedule_task_command/dal/model"

func CreateConvert(c []*HeaderTemplateCreate) []*model.HeaderTemplate {
	result := make([]*model.HeaderTemplate, 0, len(c))
	for _, item := range c {
		i := model.HeaderTemplate{
			Name: item.Name,
			Data: item.Data,
		}
		result = append(result, &i)
	}
	return result
}

func UpdateConvert(htMap map[int]model.HeaderTemplate, uht []*HeaderTemplateUpdate) (result []*model.HeaderTemplate, err error) {
	for _, u := range uht {
		ht, ok := htMap[int(u.ID)]
		if !ok {
			err = HeaderNotFound(int(u.ID))
			return
		}
		if u.Name != nil {
			ht.Name = *u.Name
		}
		if u.Data != nil {
			ht.Data = *u.Data
		}
		result = append(result, &ht)
	}
	return
}
