package e_command_template

import (
	"schedule_task_command/dal/model"
)

func Format(ct []model.CommandTemplate) []CommandTemplate {
	result := make([]CommandTemplate, 0, len(ct))
	for _, item := range ct {
		i := M2Entry(item)
		result = append(result, i)
	}
	return result
}

func M2Entry(mct model.CommandTemplate) (ct CommandTemplate) {
	p := mct.Protocol
	ct = CommandTemplate{
		ID:          mct.ID,
		Name:        mct.Name,
		Protocol:    S2Protocol(&p),
		Timeout:     mct.Timeout,
		Description: mct.Description,
		Host:        mct.Host,
		Port:        mct.Port,
		UpdatedAt:   mct.UpdatedAt,
		CreatedAt:   mct.CreatedAt,
		Tags:        mct.Tags,
		Variable:    mct.Variable,
	}
	if mct.Http != nil {
		method := mct.Http.Method
		ct.Http = &HTTPSCommand{
			Method:            S2HTTPMethod(&method),
			URL:               mct.Http.URL,
			AuthorizationType: mct.Http.AuthorizationType,
			Params:            mct.Http.Params,
			Header:            mct.Http.Header,
			BodyType:          S2BodyType(mct.Http.BodyType),
			Body:              mct.Http.Body,
		}
	}
	if mct.Mqtt != nil {
		ct.Mqtt = &MqttCommand{
			Topic:   mct.Mqtt.Topic,
			Header:  mct.Mqtt.Header,
			Message: mct.Mqtt.Message,
			Type:    mct.Mqtt.Type,
		}
	}
	if mct.Websocket != nil {
		ct.Websocket = &WebsocketCommand{
			URL:     mct.Websocket.URL,
			Header:  mct.Websocket.Header,
			Message: mct.Websocket.Message,
		}
	}
	if mct.Redis != nil {
		ct.Redis = &RedisCommand{
			Password: mct.Redis.Password,
			Db:       mct.Redis.Db,
			Topic:    mct.Redis.Topic,
			Message:  mct.Redis.Message,
			Type:     mct.Redis.Type,
		}
	}
	if mct.Monitor != nil {
		mResult := make([]MCondition, 0, len(mct.Monitor.MConditions))
		for _, m := range mct.Monitor.MConditions {
			i := MCondition{
				Order:         m.Order,
				CalculateType: m.CalculateType,
				PreLogicType:  m.PreLogicType,
				Value:         m.Value,
				SearchRule:    m.SearchRule,
			}
			mResult = append(mResult, i)
		}
		ct.Monitor = &Monitor{
			StatusCode:  mct.Monitor.StatusCode,
			Interval:    mct.Monitor.Interval,
			MConditions: mResult,
		}
	}
	return
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
			Variable:    item.Variable,
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

func UpdateConvert(ctMap map[int]model.CommandTemplate, uct []*CommandTemplateUpdate) (result []*model.CommandTemplate, err error) {
	for _, u := range uct {
		ct, ok := ctMap[int(u.ID)]
		if !ok {
			err = CommandTemplateNotFound(int(u.ID))
		}
		if u.Name != nil {
			ct.Name = *u.Name
		}
		if u.Protocol != ProtocolNone {
			ct.Protocol = u.Protocol.String()
		}
		if u.Timeout != nil {
			ct.Timeout = *u.Timeout
		}
		if u.Description != nil {
			ct.Description = u.Description
		}
		if u.Host != nil {
			ct.Host = *u.Host
		}
		if u.Port != nil {
			ct.Port = *u.Port
		}
		if u.Tags != nil {
			ct.Tags = u.Tags
		}
		if u.Variable != nil {
			ct.Variable = u.Variable
		}
		if u.Http != nil && ct.Http != nil {
			var bodyType *string
			if *bodyType == "" {
				bodyType = nil
			} else {
				*bodyType = u.Http.BodyType.String()
			}
			ct.Http = &model.HTTPSCommand{
				Method:            u.Http.Method.String(),
				URL:               u.Http.URL,
				AuthorizationType: u.Http.AuthorizationType,
				Params:            u.Http.Params,
				Header:            u.Http.Header,
				BodyType:          bodyType,
				Body:              u.Http.Body,
			}
		}
		if u.Mqtt != nil && ct.Mqtt != nil {
			ct.Mqtt = &model.MqttCommand{
				Topic:   u.Mqtt.Topic,
				Header:  u.Mqtt.Header,
				Message: u.Mqtt.Message,
				Type:    u.Mqtt.Type,
			}
		}
		if u.Websocket != nil && ct.Websocket != nil {
			ct.Websocket = &model.WebsocketCommand{
				URL:     u.Websocket.URL,
				Header:  u.Websocket.Header,
				Message: u.Websocket.Message,
			}
		}
		if u.Redis != nil && ct.Redis != nil {
			ct.Redis = &model.RedisCommand{
				Password: u.Redis.Password,
				Db:       u.Redis.Db,
				Topic:    u.Redis.Topic,
				Message:  u.Redis.Message,
				Type:     u.Redis.Type,
			}
		}
		if u.Monitor != nil && ct.Monitor != nil {
			mResult := make([]model.MCondition, 0, len(u.Monitor.MConditions))
			for _, m := range u.Monitor.MConditions {
				i := model.MCondition{
					Order:         m.Order,
					CalculateType: m.CalculateType,
					PreLogicType:  m.PreLogicType,
					Value:         m.Value,
					SearchRule:    m.SearchRule,
				}
				mResult = append(mResult, i)
			}
			ct.Monitor = &model.Monitor{
				StatusCode:  u.Monitor.StatusCode,
				Interval:    u.Monitor.Interval,
				MConditions: mResult,
			}
		}
		result = append(result, &ct)
	}
	return
}

func S2Protocol(s *string) Protocol {
	if s == nil {
		return ProtocolNone
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
		return ProtocolNone
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