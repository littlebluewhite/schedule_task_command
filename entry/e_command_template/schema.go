package e_command_template

import (
	"github.com/goccy/go-json"
	"time"
)

type CommandTemplate struct {
	ID          int32             `json:"id"`
	Name        string            `json:"name"`
	Protocol    string            `json:"protocol"`
	Timeout     int32             `json:"timeout"`
	Description *string           `json:"description"`
	Host        string            `json:"host"`
	Port        string            `json:"port"`
	UpdatedAt   *time.Time        `json:"updated_at"`
	CreatedAt   *time.Time        `json:"created_at"`
	Http        *HTTPSCommand     `json:"http,omitempty"`
	Mqtt        *MqttCommand      `json:"mqtt,omitempty"`
	Websocket   *WebsocketCommand `json:"websocket,omitempty"`
	Redis       *RedisCommand     `json:"redis,omitempty"`
	Monitor     *Monitor          `json:"monitor"`
	Tags        json.RawMessage   `json:"tags"`
}

type HTTPSCommand struct {
	Method            string           `json:"method"`
	URL               string           `json:"url"`
	AuthorizationType *string          `json:"authorization_type"`
	Params            json.RawMessage  `json:"params"`
	Header            json.RawMessage  `json:"header"`
	BodyType          *string          `json:"body_type"`
	Body              *json.RawMessage `json:"body"`
}

type MqttCommand struct {
	Topic   string           `json:"topic" binding:"required"`
	Header  json.RawMessage  `json:"header"`
	Message *json.RawMessage `json:"message"`
	Type    string           `json:"type" binding:"required"`
}

type WebsocketCommand struct {
	URL     string          `json:"url" binding:"required"`
	Header  json.RawMessage `json:"header"`
	Message *string         `json:"message"`
}

type RedisCommand struct {
	Password *string          `json:"password"`
	Db       *int32           `json:"db"`
	Topic    *string          `json:"topic"`
	Message  *json.RawMessage `json:"message"`
	Type     string           `json:"type" binding:"required"`
}

type Monitor struct {
	StatusCode  int32        `json:"status_code" binding:"required"`
	Interval    int32        `json:"interval" binding:"required"`
	MConditions []MCondition `json:"m_conditions"`
}

type MCondition struct {
	Order         int32   `json:"order"`
	CalculateType string  `json:"calculate_type"`
	PreLogicType  *string `json:"pre_logic_type"`
	Value         string  `json:"value"`
	SearchRule    string  `json:"search_rule"`
}

type CommandTemplateCreate struct {
	Name        string            `json:"name" binding:"required"`
	Protocol    string            `json:"protocol" binding:"required"`
	Timeout     int32             `json:"timeout" binding:"required"`
	Description *string           `json:"description"`
	Host        string            `json:"host" binding:"required"`
	Port        string            `json:"port" binding:"required"`
	Http        *HTTPSCommand     `json:"http"`
	Mqtt        *MqttCommand      `json:"mqtt"`
	Websocket   *WebsocketCommand `json:"websocket"`
	Redis       *RedisCommand     `json:"redis"`
	Monitor     *Monitor          `json:"monitor"`
	Tags        json.RawMessage   `json:"tags"`
}
