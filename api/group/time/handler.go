package time

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_time"
	"schedule_task_command/util"
)

type hOperate interface {
	GetHistory(id, templateId, start, stop, isTime string) ([]e_time.PublishTime, error)
}

type Handler struct {
	o hOperate
	l api.Logger
}

func NewHandler(o hOperate, l api.Logger) *Handler {
	return &Handler{
		o: o,
		l: l,
	}
}

// GetHistory swagger
// @Summary get time history
// @Tags    time
// @Accept  json
// @Produce json
// @Param       id  query     int false "time id"
// @Param       template_id  query     int false "time template id"
// @Param       is_time  query     string false "is_time" Enums(false, true)
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Success 200 {array} e_time.PublishTime
// @Router  /api/time/history [get]
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	templateId := c.Query("template_id")
	id := c.Query("id")
	isTime := c.Query("is_time")
	start := c.Query("start")
	if start == "" {
		return util.Err(c, NoStartTime, 0)
	}
	stop := c.Query("stop")
	data, err := h.o.GetHistory(id, templateId, start, stop, isTime)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(data)
}
