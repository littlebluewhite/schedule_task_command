package e_command_template

import (
	"schedule_task_command/dal/model"
)

func Format(ct []model.CommandTemplate) []CommandTemplate {
	result := make([]CommandTemplate, 0, len(ct))
	for _, item := range ct {
		p := item.Protocol
		i := CommandTemplate{
			ID:          item.ID,
			Name:        item.Name,
			Protocol:    S2Protocol(&p),
			Timeout:     item.Timeout,
			Description: item.Description,
			Host:        item.Host,
			Port:        item.Port,
			UpdatedAt:   item.UpdatedAt,
			CreatedAt:   item.CreatedAt,
			Tags:        item.Tags,
		}
		if item.Http != nil {
			method := item.Http.Method
			i.Http = &HTTPSCommand{
				Method:            S2HTTPMethod(&method),
				URL:               item.Http.URL,
				AuthorizationType: item.Http.AuthorizationType,
				Params:            item.Http.Params,
				Header:            item.Http.Header,
				BodyType:          S2BodyType(item.Http.BodyType),
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
		i := model.CommandTemplate{
			Name:        item.Name,
			Protocol:    item.Protocol.String(),
			Timeout:     item.Timeout,
			Description: item.Description,
			Host:        item.Host,
			Port:        item.Port,
			Tags:        item.Tags,
		}
		if item.Http != nil {
			var bodyType *string
			if *bodyType == "" {
				bodyType = nil
			} else {
				*bodyType = item.Http.BodyType.String()
			}
			i.Http = &model.HTTPSCommand{
				Method:            item.Http.Method.String(),
				URL:               item.Http.URL,
				AuthorizationType: item.Http.AuthorizationType,
				Params:            item.Http.Params,
				Header:            item.Http.Header,
				BodyType:          bodyType,
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

func S2Protocol(s *string) Protocol {
	if s == nil {
		return Http
	}
	switch *s {
	case "http":
		return Http
	case "websocket":
		return Websocket
	case "mqtt":
		return Mqtt
	case "redis_topic":
		return RedisTopic
	default:
		return Http
	}
}

func S2HTTPMethod(s *string) HTTPMethod {
	if s == nil {
		return GET
	}
	switch *s {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PATCH":
		return PATCH
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	default:
		return GET
	}
}

func S2BodyType(s *string) BodyType {
	if s == nil {
		return BodyTypeNone
	}
	switch *s {
	case "text":
		return Text
	case "html":
		return HTML
	case "xml":
		return XML
	case "form_data":
		return FormData
	case "x_www_form_urlencoded":
		return XWWWFormUrlencoded
	case "json":
		return Json
	default:
		return BodyTypeNone
	}
}
