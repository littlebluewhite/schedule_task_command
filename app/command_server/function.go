package command_server

import (
	"github.com/littlebluewhite/schedule_task_command/entry/e_command"
	"strconv"
)

func addDefaultVariables(com e_command.Command) e_command.Command {
	v := com.Variables
	v["__command_id__"] = strconv.FormatUint(com.ID, 10)
	com.Variables = v
	return com
}

func addDefaultParserReturn(parserReturn map[string]string, com e_command.Command) map[string]string {
	parserReturn["__command_id__"] = strconv.FormatUint(com.ID, 10)
	parserReturn["__status_code__"] = strconv.Itoa(com.StatusCode)
	parserReturn["__status__"] = com.Status.String()
	return parserReturn
}
