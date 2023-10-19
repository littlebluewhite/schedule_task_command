package task

import "schedule_task_command/util"

var NoStartTime = util.MyErr("No start time input")
var HistoryStatusErr = util.MyErr("History Status input error")

type CancelBody struct {
	ClientMessage string `json:"client_message"`
}
