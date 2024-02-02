package group

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/entry/e_log"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
)

type hOperate interface {
	ReadLog(start, stop, account, ip, method, module, statusCode string) ([]e_log.Log, error)
}

type Handler struct {
	o hOperate
	l logFile.LogFile
}

func NewHandler(o hOperate, l logFile.LogFile) *Handler {
	return &Handler{
		o: o,
		l: l,
	}
}

// GetHistory swagger
// @Summary get logs history
// @Tags    Log
// @Accept  json
// @Produce json
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Param       account  query     string false "account"
// @Param       ip  query     string false "ip"
// @Param       method  query     string false "method" Enums(GET, POST, PATCH, PUT, DELETE)
// @Param       module  query     string false "module"
// @Param       status_code  query     string false "status_code"
// @Success 200 {array} e_log.Log
// @Router  /api/logs [get]
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	start := c.Query("start")
	account := c.Query("account")
	ip := c.Query("ip")
	method := c.Query("method")
	module := c.Query("module")
	statusCode := c.Query("statusCode")
	if start == "" {
		return util.Err(c, util.MyErr("No start time input"), 0)
	}
	stop := c.Query("stop")
	data, err := h.o.ReadLog(start, stop, account, ip, method, module, statusCode)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(data)
}
