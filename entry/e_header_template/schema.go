package e_header_template

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/util"
)

type HeaderTemplate struct {
	ID   int32           `json:"id"`
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type HeaderTemplateCreate struct {
	Name string          `json:"name" binding:"required"`
	Data json.RawMessage `json:"data" binding:"required"`
}

type HeaderTemplateUpdate struct {
	ID   int32            `json:"id" binding:"required"`
	Name *string          `json:"name"`
	Data *json.RawMessage `json:"data"`
}

func HeaderNotFound(id int) util.MyErr {
	e := fmt.Sprintf("header id: %d not found", id)
	return util.MyErr(e)
}
