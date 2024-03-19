package header_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_header_template"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
)

type hOperate interface {
	List() ([]model.HeaderTemplate, error)
	Find(ids []int32) ([]model.HeaderTemplate, error)
	Create([]*e_header_template.HeaderTemplateCreate) ([]model.HeaderTemplate, error)
	Update([]*e_header_template.HeaderTemplateUpdate) error
	Delete([]int32) error
	ReloadCache() error
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

// GetheaderTemplates swagger
// @Summary     Show all header templates
// @Description Get all header templates
// @Tags        header_template
// @Produce     json
// @Success     200 {array} e_header_template.HeaderTemplate
// @Router      /api/header_template/ [get]
func (h *Handler) GetheaderTemplates(c *fiber.Ctx) error {
	ht, err := h.o.List()
	if err != nil {
		h.l.Error().Println("GetheaderTemplates: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Info().Println("GetheaderTemplates: success")
	return c.Status(200).JSON(ht)
}

// GetHeaderTemplateById swagger
// @Summary     Show header templates
// @Description Get header templates by id
// @Tags        header_template
// @Produce     json
// @Param       id  path     int true "header template id"
// @Success     200 {object} e_header_template.HeaderTemplate
// @Router      /api/header_template/{id} [get]
func (h *Handler) GetHeaderTemplateById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Error().Println("GetHeaderTemplateById: ", err)
		return util.Err(c, err, 1)
	}
	ht, err := h.o.Find([]int32{int32(id)})
	if err != nil {
		h.l.Error().Println("GetHeaderTemplateById: ", err)
		return util.Err(c, err, 2)
	}
	h.l.Info().Println("GetHeaderTemplateById: success")
	return c.Status(200).JSON(ht[0])
}

// AddHeaderTemplate swagger
// @Summary Create header templates
// @Tags    header_template
// @Accept  json
// @Produce json
// @Param   header_template body     []e_header_template.HeaderTemplateCreate true "header template body"
// @Success 200           {array} e_header_template.HeaderTemplate
// @Router  /api/header_template/ [post]
func (h *Handler) AddHeaderTemplate(c *fiber.Ctx) error {
	entry := []*e_header_template.HeaderTemplateCreate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("AddHeaderTemplate: ", err)
		return util.Err(c, err, 0)
	}
	result, err := h.o.Create(entry)
	if err != nil {
		h.l.Error().Println("AddHeaderTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(result)
}

// UpdateHeaderTemplate swagger
// @Summary Update header templates
// @Tags    header_template
// @Accept  json
// @Produce json
// @Param   header_template body     []e_header_template.HeaderTemplateUpdate true "modify header template body"
// @Success 200           {string} string "update successfully"
// @Router  /api/header_template/ [patch]
func (h *Handler) UpdateHeaderTemplate(c *fiber.Ctx) error {
	entry := []*e_header_template.HeaderTemplateUpdate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("UpdateHeaderTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Update(entry)
	if err != nil {
		h.l.Error().Println("UpdateHeaderTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("update successfully")
}

// DeleteHeaderTemplate swagger
// @Summary Delete header templates
// @Tags    header_template
// @Produce json
// @Param ids body []int true "header_template id"
// @Success 200 {string} string "delete successfully"
// @Router  /api/header_template/ [delete]
func (h *Handler) DeleteHeaderTemplate(c *fiber.Ctx) error {
	entry := make([]int32, 0, 10)
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("DeleteHeaderTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Delete(entry)
	if err != nil {
		h.l.Error().Println("DeleteHeaderTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("delete successfully")
}
