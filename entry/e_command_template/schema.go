package e_command_template

import (
	"fmt"
	"github.com/goccy/go-json"
	"schedule_task_command/util"
	"time"
)

type CommandTemplate struct {
	ID           int32             `json:"id"`
	Name         string            `json:"name"`
	Protocol     Protocol          `json:"protocol"`
	Timeout      int32             `json:"timeout"`
	Description  *string           `json:"description"`
	Host         string            `json:"host"`
	Port         string            `json:"port"`
	UpdatedAt    *time.Time        `json:"updated_at"`
	CreatedAt    *time.Time        `json:"created_at"`
	Http         *HTTPSCommand     `json:"http,omitempty"`
	Mqtt         *MqttCommand      `json:"mqtt,omitempty"`
	Websocket    *WebsocketCommand `json:"websocket,omitempty"`
	Redis        *RedisCommand     `json:"redis,omitempty"`
	Monitor      *Monitor          `json:"monitor"`
	Tags         json.RawMessage   `json:"tags"`
	Variable     json.RawMessage   `json:"variable"`
	ParserReturn []ParserReturn    `json:"parser_return"`
}

type HTTPSCommand struct {
	Method            HTTPMethod      `json:"method"`
	URL               string          `json:"url"`
	AuthorizationType *string         `json:"authorization_type"`
	Params            json.RawMessage `json:"params"`
	Header            json.RawMessage `json:"header"`
	BodyType          BodyType        `json:"body_type"`
	Body              json.RawMessage `json:"body"`
}

type MqttCommand struct {
	Topic   string          `json:"topic" binding:"required"`
	Header  json.RawMessage `json:"header"`
	Message json.RawMessage `json:"message"`
	Type    string          `json:"type" binding:"required"`
}

type WebsocketCommand struct {
	URL     string          `json:"url" binding:"required"`
	Header  json.RawMessage `json:"header"`
	Message *string         `json:"message"`
}

type RedisCommand struct {
	Password *string         `json:"password"`
	Db       *int32          `json:"db"`
	Topic    *string         `json:"topic"`
	Message  json.RawMessage `json:"message"`
	Type     string          `json:"type" binding:"required"`
}

type Monitor struct {
	StatusCode  int32        `json:"status_code" binding:"required"`
	Interval    int32        `json:"interval" binding:"required"`
	MConditions []MCondition `json:"m_conditions"`
}

type MCondition struct {
	ID            int32   `json:"id"`
	Order         int32   `json:"order"`
	CalculateType string  `json:"calculate_type"`
	PreLogicType  *string `json:"pre_logic_type"`
	Value         string  `json:"value"`
	SearchRule    string  `json:"search_rule"`
}

type CommandTemplateCreate struct {
	Name         string            `json:"name" binding:"required"`
	Protocol     Protocol          `json:"protocol" binding:"required"`
	Timeout      int32             `json:"timeout" binding:"required"`
	Description  *string           `json:"description"`
	Host         string            `json:"host" binding:"required"`
	Port         string            `json:"port" binding:"required"`
	Http         *HTTPSCommand     `json:"http"`
	Mqtt         *MqttCommand      `json:"mqtt"`
	Websocket    *WebsocketCommand `json:"websocket"`
	Redis        *RedisCommand     `json:"redis"`
	Monitor      *Monitor          `json:"monitor"`
	Tags         json.RawMessage   `json:"tags"`
	Variable     json.RawMessage   `json:"variable"`
	ParserReturn []ParserReturn    `json:"parser_return"`
}

type CommandTemplateUpdate struct {
	ID           int32             `json:"id"`
	Name         *string           `json:"name"`
	Protocol     Protocol          `json:"protocol"`
	Timeout      *int32            `json:"timeout"`
	Description  *string           `json:"description"`
	Host         *string           `json:"host"`
	Port         *string           `json:"port"`
	Http         *HTTPSCommand     `json:"http"`
	Mqtt         *MqttCommand      `json:"mqtt"`
	Websocket    *WebsocketCommand `json:"websocket"`
	Redis        *RedisCommand     `json:"redis"`
	Monitor      *MonitorUpdate    `json:"monitor"`
	Tags         json.RawMessage   `json:"tags"`
	Variable     json.RawMessage   `json:"variable"`
	ParserReturn []ParserReturn    `json:"parser_return"`
}

type MonitorUpdate struct {
	StatusCode  *int32       `json:"status_code" binding:"required"`
	Interval    *int32       `json:"interval" binding:"required"`
	MConditions []MCondition `json:"m_conditions"`
}

type SendCommandTemplate struct {
	TemplateId     int32             `json:"template_id"`
	TriggerFrom    []string          `json:"trigger_from"`
	TriggerAccount string            `json:"trigger_account"`
	Token          string            `json:"token"`
	Variables      map[string]string `json:"variables"`
}

type ParserReturn struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key"`
	SearchRule string `json:"search_rule"`
}

var CannotFindTemplate = util.MyErr("can not find Command template")

func CommandTemplateNotFound(id int) util.MyErr {
	e := fmt.Sprintf("command template id: %d not found", id)
	return util.MyErr(e)
}
