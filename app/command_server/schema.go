package command_server

import (
	"schedule_task_command/entry/e_command"
	"schedule_task_command/util"
	"sync"
)

type Protocol int

const (
	https Protocol = iota
	websocket
	mqtt
	redisTopic
)

func (p Protocol) String() string {
	return [...]string{"http", "websocket", "mqtt", "redis_topic"}[p]
}

type chs struct {
	rec chan e_command.Command
	mu  *sync.RWMutex
}

type httpHeader struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsActive bool   `json:"is_active"`
	DataType string `json:"data_type"`
}

type analyzeResult struct {
	getSuccess  bool
	valueResult any
	arrayResult []any
}

type assertResult struct {
	order         int32
	assertSuccess bool
	preLogicType  *string
}

var (
	valueCalculate = []string{"=", "<", "<=", ">", ">=", "!="}
	sliceCalculate = []string{"include", "exclude"}
)

type SendCommand struct {
	TemplateId     int      `json:"template_id"`
	TriggerFrom    []string `json:"trigger_from"`
	TriggerAccount string   `json:"trigger_account"`
	Token          string   `json:"token"`
}

var CannotFindTemplate = util.MyErr("can not find Command template")
var CommandCanceled = util.MyErr("Command has been canceled")
var CommandTimeout = util.MyErr("Command not match monitor and timeout")
var HttpTimeout = util.MyErr("http request timeout")
var HttpCodeErr = util.MyErr("http request status code error")
var ConditionFailed = util.MyErr("monitor condition is not suitable now")
var SendToRedisErr = util.MyErr("send task to redis cannot format")
var CommandNotFind = util.MyErr("can not find command")
var CommandCannotCancel = util.MyErr("command cannot be canceled")
