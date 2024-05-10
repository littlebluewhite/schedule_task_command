package schedule

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_schedule"
	"schedule_task_command/util"
)

type hOperate interface {
	List() ([]model.Schedule, error)
	Find(ids []int32) ([]model.Schedule, error)
	Create([]*e_schedule.ScheduleCreate) ([]model.Schedule, error)
	Update([]*e_schedule.ScheduleUpdate) error
	Delete([]int32) error
	ReloadCache() error
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

// GetSchedules swagger
// @Summary     Show all schedules
// @Description Get all schedules
// @Tags        schedule
// @Produce     json
// @Success     200 {array} e_schedule.Schedule
// @Router      /api/schedule/ [get]
func (h *Handler) GetSchedules(c *fiber.Ctx) error {
	s, err := h.o.List()
	if err != nil {
		h.l.Errorln("GetheaderTemplates: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Infoln("GetheaderTemplates: success")
	return c.Status(200).JSON(e_schedule.Format(s))
}

// GetScheduleById swagger
// @Summary     Show schedules
// @Description Get schedules by id
// @Tags        schedule
// @Produce     json
// @Param       id  path     int true "schedule id"
// @Success     200 {object} e_schedule.Schedule
// @Router      /api/schedule/{id} [get]
func (h *Handler) GetScheduleById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Errorln("GetScheduleById: ", err)
		return util.Err(c, err, 1)
	}
	s, err := h.o.Find([]int32{int32(id)})
	if err != nil {
		h.l.Errorln("GetScheduleById: ", err)
		return util.Err(c, err, 2)
	}
	h.l.Infoln("GetScheduleById: success")
	return c.Status(200).JSON(e_schedule.Format(s)[0])
}

// AddSchedule swagger
// @Summary Create schedules
// @Tags    schedule
// @Accept  json
// @Produce json
// @Param   schedule body     []e_schedule.ScheduleCreate true "schedule body"
// @Success 200           {array} e_schedule.Schedule
// @Router  /api/schedule/ [post]
func (h *Handler) AddSchedule(c *fiber.Ctx) error {
	entry := []*e_schedule.ScheduleCreate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("AddSchedule: ", err)
		return util.Err(c, err, 0)
	}
	s, err := h.o.Create(entry)
	if err != nil {
		h.l.Errorln("AddSchedule: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_schedule.Format(s))
}

// UpdateSchedule swagger
// @Summary Update schedules
// @Tags    schedule
// @Accept  json
// @Produce json
// @Param   schedule body     []e_schedule.ScheduleUpdate true "modify schedule body"
// @Success 200           {string} string "update successfully"
// @Router  /api/schedule/ [patch]
func (h *Handler) UpdateSchedule(c *fiber.Ctx) error {
	entry := []*e_schedule.ScheduleUpdate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("UpdateSchedule: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Update(entry)
	if err != nil {
		h.l.Errorln("UpdateSchedule: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("update successfully")
}

// DeleteSchedule swagger
// @Summary Delete schedules
// @Tags    schedule
// @Produce json
// @Param ids body []int true "schedule id"
// @Success 200 {string} string "delete successfully"
// @Router  /api/schedule/ [delete]
func (h *Handler) DeleteSchedule(c *fiber.Ctx) error {
	entry := make([]int32, 0, 10)
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("DeleteSchedule: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Delete(entry)
	if err != nil {
		h.l.Errorln("DeleteSchedule: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("delete successfully")
}
