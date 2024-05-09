package ping

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"time"
)

type hOperate interface {
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

// GetPing swagger
// @Summary    test ping
// @Description test ping
// @Tags        ping
// @Produce     json
// @Success     200 {object} ping.SwaggerPing
// @Router      /api/ping/test [get]
func (h *Handler) GetPing(c *fiber.Ctx) error {
	h.l.Infoln("get ping: example")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "test",
	})
}

// GetListPing swagger
// @Summary     return list ping
// @Description test list ping
// @Tags        ping
// @Produce     json
// @Success     200 {array} ping.SwaggerListPing
// @Router      /api/ping/list [get]
func (h *Handler) GetListPing(c *fiber.Ctx) error {
	data := []map[string]interface{}{
		{
			"name": "wilson",
			"age":  5,
			"time": time.Now(),
		},
		{
			"name": "phoebe",
			"age":  4,
		},
	}
	h.l.Infoln("get ping list: data: ", data)
	return c.Status(fiber.StatusOK).JSON(data)
}
