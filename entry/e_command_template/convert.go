package e_command_template

import (
	"schedule_task_command/dal/model"
)

func Format(ct []model.CommandTemplate) []CommandTemplate {
	result := make([]CommandTemplate, 0, len(ct))
	for _, item := range ct {
		i := CommandTemplate{
			ID:          item.ID,
			Name:        item.Name,
			Protocol:    item.Protocol,
			Timeout:     item.Timeout,
			Description: item.Description,
			Host:        item.Host,
			Port:        item.Port,
			UpdatedAt:   item.UpdatedAt,
			CreatedAt:   item.CreatedAt,
			Tags:        item.Tags,
		}
		if item.Http != nil {
			i.Http = &HTTPSCommand{
				Method:            item.Http.Method,
				URL:               item.Http.URL,
				AuthorizationType: item.Http.AuthorizationType,
				Params:            item.Http.Params,
				Header:            item.Http.Header,
				BodyType:          item.Http.BodyType,
				Body:              item.Http.Body,
			}
		}
		if item.Mqtt != nil {
			i.Mqtt = &MqttCommand{
				Topic:   item.Mqtt.Topic,
				Header:  item.Mqtt.Header,
				Message: item.Mqtt.Message,
				Type:    item.Mqtt.Type,
			}
		}
		if item.Websocket != nil {
			i.Websocket = &WebsocketCommand{
				URL:     item.Websocket.URL,
				Header:  item.Websocket.Header,
				Message: item.Websocket.Message,
			}
		}
		if item.Redis != nil {
			i.Redis = &RedisCommand{
				Password: item.Redis.Password,
				Db:       item.Redis.Db,
				Topic:    item.Redis.Topic,
				Message:  item.Redis.Message,
				Type:     item.Redis.Type,
			}
		}
		if item.Monitor != nil {
			mResult := make([]MCondition, 0, len(item.Monitor.MConditions))
			for _, m := range item.Monitor.MConditions {
				i := MCondition{
					Order:         m.Order,
					CalculateType: m.CalculateType,
					PreLogicType:  m.PreLogicType,
					Value:         m.Value,
					SearchRule:    m.SearchRule,
				}
				mResult = append(mResult, i)
			}
			i.Monitor = &Monitor{
				StatusCode:  item.Monitor.StatusCode,
				Interval:    item.Monitor.Interval,
				MConditions: mResult,
			}
		}
		result = append(result, i)
	}
	return result
}

func CreateConvert(c []*CommandTemplateCreate) []*model.CommandTemplate {
	result := make([]*model.CommandTemplate, 0, len(c))
	for _, item := range c {
		mResult := make([]model.MCondition, 0, len(item.Monitor.MConditions))
		for _, m := range item.Monitor.MConditions {
			i := model.MCondition{
				Order:         m.Order,
				CalculateType: m.CalculateType,
				PreLogicType:  m.PreLogicType,
				Value:         m.Value,
				SearchRule:    m.SearchRule,
			}
			mResult = append(mResult, i)
		}
		i := model.CommandTemplate{
			Name:        item.Name,
			Protocol:    item.Protocol,
			Timeout:     item.Timeout,
			Description: item.Description,
			Host:        item.Host,
			Port:        item.Port,
			Tags:        item.Tags,
		}
		if item.Http != nil {
			i.Http = &model.HTTPSCommand{
				Method:            item.Http.Method,
				URL:               item.Http.URL,
				AuthorizationType: item.Http.AuthorizationType,
				Params:            item.Http.Params,
				Header:            item.Http.Header,
				BodyType:          item.Http.BodyType,
				Body:              item.Http.Body,
			}
		}
		if item.Mqtt != nil {
			i.Mqtt = &model.MqttCommand{
				Topic:   item.Mqtt.Topic,
				Header:  item.Mqtt.Header,
				Message: item.Mqtt.Message,
				Type:    item.Mqtt.Type,
			}
		}
		if item.Websocket != nil {
			i.Websocket = &model.WebsocketCommand{
				URL:     item.Websocket.URL,
				Header:  item.Websocket.Header,
				Message: item.Websocket.Message,
			}
		}
		if item.Redis != nil {
			i.Redis = &model.RedisCommand{
				Password: item.Redis.Password,
				Db:       item.Redis.Db,
				Topic:    item.Redis.Topic,
				Message:  item.Redis.Message,
				Type:     item.Redis.Type,
			}
		}
		if item.Monitor != nil {
			i.Monitor = &model.Monitor{
				StatusCode:  item.Monitor.StatusCode,
				Interval:    item.Monitor.Interval,
				MConditions: mResult,
			}
		}
		result = append(result, &i)
	}
	return result
}
