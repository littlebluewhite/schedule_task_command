package time_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_time_template"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
)

type hOperate interface {
	List() ([]model.TimeTemplate, error)
	Find(ids []int32) ([]model.TimeTemplate, error)
	Create([]*e_time_template.TimeTemplateCreate) ([]model.TimeTemplate, error)
	Update([]*e_time_template.TimeTemplateUpdate) error
	Delete([]int32) error
	ReloadCache() error
	CheckTime(id int, c CheckTime) bool
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

// GetTimeTemplates swagger
// @Summary     Show all time templates
// @Description Get all time templates
// @Tags        time_template
// @Produce     json
// @Success     200 {array} e_time_template.TimeTemplate
// @Router     /api/time_template/ [get]
func (h *Handler) GetTimeTemplates(c *fiber.Ctx) error {
	tt, err := h.o.List()
	if err != nil {
		h.l.Error().Println("GetTimeTemplates: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Info().Println("GetTimeTemplates: success")
	return c.Status(200).JSON(e_time_template.Format(tt))
}

// GetTimeTemplateById swagger
// @Summary     Show time templates
// @Description Get time templates by id
// @Tags        time_template
// @Produce     json
// @Param       id  path     int true "time template id"
// @Success     200 {object} e_time_template.TimeTemplate
// @Router      /api/time_template/{id} [get]
func (h *Handler) GetTimeTemplateById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Error().Println("GetTimeTemplateById: ", err)
		return util.Err(c, err, 0)
	}
	tt, err := h.o.Find([]int32{int32(id)})
	if err != nil {
		h.l.Error().Println("GetTimeTemplateById: ", err)
		return util.Err(c, err, 0)
	}
	result := e_time_template.Format(tt)
	h.l.Info().Println("GetTimeTemplateById: success")
	return c.Status(200).JSON(result[0])
}

// AddTimeTemplate swagger
// @Summary Create time templates
// @Tags    time_template
// @Accept  json
// @Produce json
// @Param   time_template body     []e_time_template.TimeTemplateCreate true "time template body"
// @Success 200           {array} e_time_template.TimeTemplate
// @Router  /api/time_template/ [post]
func (h *Handler) AddTimeTemplate(c *fiber.Ctx) error {
	entry := []*e_time_template.TimeTemplateCreate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("AddTimeTemplate: ", err)
		return util.Err(c, err, 0)
	}
	tt, err := h.o.Create(entry)
	if err != nil {
		h.l.Error().Println("AddTimeTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_time_template.Format(tt))
}

// UpdateTimeTemplate swagger
// @Summary Update time templates
// @Tags    time_template
// @Accept  json
// @Produce json
// @Param   time_template body     []e_time_template.TimeTemplateUpdate true "modify time template body"
// @Success 200           {string} string "update successfully"
// @Router  /api/time_template/ [patch]
func (h *Handler) UpdateTimeTemplate(c *fiber.Ctx) error {
	entry := []*e_time_template.TimeTemplateUpdate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("UpdateTimeTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Update(entry)
	if err != nil {
		h.l.Error().Println("UpdateTimeTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("update successfully")
}

// DeleteTimeTemplate swagger
// @Summary Delete time templates
// @Tags    time_template
// @Produce json
// @Param ids body []int true "time_template id"
// @Success 200 {string} string "delete successfully"
// @Router  /api/time_template/ [delete]
func (h *Handler) DeleteTimeTemplate(c *fiber.Ctx) error {
	entry := make([]int32, 0, 10)
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("DeleteTimeTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Delete(entry)
	if err != nil {
		h.l.Error().Println("DeleteTimeTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("delete successfully")
}

// CheckTime swagger
// @Summary Check time templates
// @Tags    time_template
// @Accept  json
// @Produce json
// @Param       id  path     int true "time template id"
// @Param   checkTime body     time_template.CheckTime true "check time body"
// @Success 200           {boolean} boolean
// @Router  /api/time_template/checkTime/{id} [post]
func (h *Handler) CheckTime(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Error().Println("CheckTime: ", err)
		return util.Err(c, err, 0)
	}
	entry := CheckTime{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("CheckTime: ", err)
		return util.Err(c, err, 0)
	}
	isTime := h.o.CheckTime(id, entry)
	return c.Status(200).JSON(isTime)
}
