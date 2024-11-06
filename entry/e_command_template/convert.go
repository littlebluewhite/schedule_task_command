package e_command_template

import (
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/dal/model"
	"github.com/littlebluewhite/schedule_task_command/util"
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
	parserReturn := make([]ParserReturn, 0, len(mct.ParserReturn))
	for _, pr := range mct.ParserReturn {
		parserReturn = append(parserReturn, ParserReturn{
			ID:         pr.ID,
			Name:       pr.Name,
			Key:        pr.Key,
			SearchRule: pr.SearchRule,
		})
	}
	ct = CommandTemplate{
		ID:           mct.ID,
		Name:         mct.Name,
		Visible:      mct.Visible,
		Protocol:     S2Protocol(&p),
		Timeout:      mct.Timeout,
		Description:  mct.Description,
		Host:         mct.Host,
		Port:         mct.Port,
		UpdatedAt:    mct.UpdatedAt,
		CreatedAt:    mct.CreatedAt,
		Tags:         mct.Tags,
		Variable:     mct.Variable,
		VariableKey:  mct.VariableKey,
		ParserReturn: parserReturn,
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
				ID:            m.ID,
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
		variableKey := make([]string, 0, 8)
		parserReturn := make([]model.ParserReturn, 0, len(item.ParserReturn))
		for _, pr := range item.ParserReturn {
			parserReturn = append(parserReturn, model.ParserReturn{
				ID:         pr.ID,
				Name:       pr.Name,
				Key:        pr.Key,
				SearchRule: pr.SearchRule,
			})
		}
		i := model.CommandTemplate{
			Name:         item.Name,
			Visible:      item.Visible,
			Protocol:     item.Protocol.String(),
			Timeout:      item.Timeout,
			Description:  item.Description,
			Host:         item.Host,
			Port:         item.Port,
			Tags:         item.Tags,
			Variable:     item.Variable,
			ParserReturn: parserReturn,
		}
		if item.Http != nil {
			var bodyType *string
			if item.Http.BodyType.String() != "" {
				bt := item.Http.BodyType.String()
				bodyType = &bt
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
			variableKey = append(variableKey, util.GetStringVariables(item.Http.URL)...)
			variableKey = append(variableKey, util.GetByteVariables(item.Http.Header)...)
			variableKey = append(variableKey, util.GetByteVariables(item.Http.Body)...)
		}
		if item.Mqtt != nil {
			i.Mqtt = &model.MqttCommand{
				Topic:   item.Mqtt.Topic,
				Header:  item.Mqtt.Header,
				Message: item.Mqtt.Message,
				Type:    item.Mqtt.Type,
			}
			variableKey = append(variableKey, util.GetStringVariables(item.Mqtt.Topic)...)
			variableKey = append(variableKey, util.GetByteVariables(item.Mqtt.Header)...)
			variableKey = append(variableKey, util.GetByteVariables(item.Mqtt.Message)...)
		}
		if item.Websocket != nil {
			i.Websocket = &model.WebsocketCommand{
				URL:     item.Websocket.URL,
				Header:  item.Websocket.Header,
				Message: item.Websocket.Message,
			}
			variableKey = append(variableKey, util.GetStringVariables(item.Websocket.URL)...)
			variableKey = append(variableKey, util.GetByteVariables(item.Websocket.Header)...)
			variableKey = append(variableKey, util.GetStringVariables(*item.Websocket.Message)...)
		}
		if item.Redis != nil {
			i.Redis = &model.RedisCommand{
				Password: item.Redis.Password,
				Db:       item.Redis.Db,
				Topic:    item.Redis.Topic,
				Message:  item.Redis.Message,
				Type:     item.Redis.Type,
			}
			variableKey = append(variableKey, util.GetStringVariables(*item.Redis.Topic)...)
			variableKey = append(variableKey, util.GetByteVariables(item.Redis.Message)...)
		}
		if item.Monitor != nil {
			mResult := make([]model.MCondition, 0, len(item.Monitor.MConditions))
			for _, m := range item.Monitor.MConditions {
				mc := model.MCondition{
					ID:            m.ID,
					Order:         m.Order,
					CalculateType: m.CalculateType,
					PreLogicType:  m.PreLogicType,
					Value:         m.Value,
					SearchRule:    m.SearchRule,
				}
				mResult = append(mResult, mc)
			}
			i.Monitor = &model.Monitor{
				StatusCode:  item.Monitor.StatusCode,
				Interval:    item.Monitor.Interval,
				MConditions: mResult,
			}
		}
		vkb, _ := json.Marshal(variableKey)
		i.VariableKey = vkb
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
		if u.Visible != nil {
			ct.Visible = *u.Visible
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
			if u.Http.BodyType.String() != "" {
				bt := u.Http.BodyType.String()
				bodyType = &bt
			}
			ct.Http.Method = u.Http.Method.String()
			ct.Http.URL = u.Http.URL
			ct.Http.AuthorizationType = u.Http.AuthorizationType
			ct.Http.Params = u.Http.Params
			ct.Http.Header = u.Http.Header
			ct.Http.BodyType = bodyType
			ct.Http.Body = u.Http.Body
		}
		if u.Mqtt != nil && ct.Mqtt != nil {
			ct.Mqtt.Topic = u.Mqtt.Topic
			ct.Mqtt.Header = u.Mqtt.Header
			ct.Mqtt.Message = u.Mqtt.Message
			ct.Mqtt.Type = u.Mqtt.Type
		}
		if u.Websocket != nil && ct.Websocket != nil {
			ct.Websocket.URL = u.Websocket.URL
			ct.Websocket.Header = u.Websocket.Header
			ct.Websocket.Message = u.Websocket.Message
		}
		if u.Redis != nil && ct.Redis != nil {
			ct.Redis.Password = u.Redis.Password
			ct.Redis.Db = u.Redis.Db
			ct.Redis.Topic = u.Redis.Topic
			ct.Redis.Message = u.Redis.Message
			ct.Redis.Type = u.Redis.Type
		}
		parserReturn := make([]model.ParserReturn, 0, len(u.ParserReturn))
		for _, pr := range u.ParserReturn {
			parserReturn = append(parserReturn, model.ParserReturn{
				ID:         pr.ID,
				Name:       pr.Name,
				Key:        pr.Key,
				SearchRule: pr.SearchRule,
			})
		}
		ct.ParserReturn = parserReturn
		if u.Monitor != nil && ct.Monitor != nil {
			mResult := make([]model.MCondition, 0, len(u.Monitor.MConditions))
			for _, m := range u.Monitor.MConditions {
				i := model.MCondition{
					ID:            m.ID,
					Order:         m.Order,
					CalculateType: m.CalculateType,
					PreLogicType:  m.PreLogicType,
					Value:         m.Value,
					SearchRule:    m.SearchRule,
				}
				mResult = append(mResult, i)
			}
			ct.Monitor.MConditions = mResult
			if u.Monitor.StatusCode != nil {
				ct.Monitor.StatusCode = *u.Monitor.StatusCode
			}
			if u.Monitor.Interval != nil {
				ct.Monitor.Interval = *u.Monitor.Interval
			}
		}
		variableKey := make([]string, 0, 8)
		if ct.Http != nil {
			variableKey = append(variableKey, util.GetStringVariables(ct.Http.URL)...)
			variableKey = append(variableKey, util.GetByteVariables(ct.Http.Header)...)
			variableKey = append(variableKey, util.GetByteVariables(ct.Http.Body)...)
		}
		if ct.Mqtt != nil {
			variableKey = append(variableKey, util.GetStringVariables(ct.Mqtt.Topic)...)
			variableKey = append(variableKey, util.GetByteVariables(ct.Mqtt.Header)...)
			variableKey = append(variableKey, util.GetByteVariables(ct.Mqtt.Message)...)
		}
		if ct.Websocket != nil {
			variableKey = append(variableKey, util.GetStringVariables(ct.Websocket.URL)...)
			variableKey = append(variableKey, util.GetByteVariables(ct.Websocket.Header)...)
			variableKey = append(variableKey, util.GetStringVariables(*ct.Websocket.Message)...)
		}
		if ct.Redis != nil {
			variableKey = append(variableKey, util.GetStringVariables(*ct.Redis.Topic)...)
			variableKey = append(variableKey, util.GetByteVariables(ct.Redis.Message)...)
		}
		vkb, _ := json.Marshal(variableKey)
		ct.VariableKey = vkb
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
