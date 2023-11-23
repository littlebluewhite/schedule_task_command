package command_server

import (
	"schedule_task_command/util"
	"sync"
)

type chs struct {
	mu *sync.RWMutex
}

type websocketManager interface {
	Broadcast(d int, message []byte)
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

var CommandCanceled = util.MyErr("Command has been canceled")
var CommandTimeout = util.MyErr("Command not match monitor and timeout")
var HttpTimeout = util.MyErr("http request timeout")
var HttpCodeErr = util.MyErr("http request status code error")
var ConditionFailed = util.MyErr("monitor condition is not suitable now")
var SendToRedisErr = util.MyErr("send command to redis cannot format")
var CommandNotFind = util.MyErr("can not find command")
var CommandCannotCancel = util.MyErr("command cannot be canceled")
var HeaderVariables = util.MyErr("header use variables failed")
var RequestErr = util.MyErr("http request failed")
var URLVariables = util.MyErr("url use variables failed")
var BodyVariables = util.MyErr("body use variables failed")
var CommandTemplateVariable = util.MyErr("command template variable failed to format")
