package command_template

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/dal/model"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command_template"
	"github.com/littlebluewhite/schedule_task_command/util"
)

type hOperate interface {
	List() ([]model.CommandTemplate, error)
	Find(ids []int32) ([]model.CommandTemplate, error)
	Create([]*e_command_template.CommandTemplateCreate) ([]model.CommandTemplate, error)
	Update([]*e_command_template.CommandTemplateUpdate) error
	Delete([]int32) error
	ReloadCache() error
	Execute(ctx context.Context, sc e_command_template.SendCommandTemplate) (id uint64, err error)
	Send(ctx context.Context, ctc e_command_template.CommandTemplateCreate) e_command.CommandPub
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

// GetCommandTemplates swagger
// @Summary     Show all command templates
// @Description Get all command templates
// @Tags        command_template
// @Produce     json
// @Success     200 {array} e_command_template.CommandTemplate
// @Router      /api/command_template/ [get]
func (h *Handler) GetCommandTemplates(c *fiber.Ctx) error {
	ct, err := h.o.List()
	result := e_command_template.Format(ct)
	if err != nil {
		h.l.Errorln("GetCommandTemplates: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Infof("GetCommandTemplates: success: %+v", result)
	return c.Status(200).JSON(result)
}

// GetCommandTemplateById swagger
// @Summary     Show command templates
// @Description Get command templates by id
// @Tags        command_template
// @Produce     json
// @Param       id  path     int true "command template id"
// @Success     200 {object} e_command_template.CommandTemplate
// @Router      /api/command_template/{id} [get]
func (h *Handler) GetCommandTemplateById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Errorln("GetCommandTemplateById: ", err)
		return util.Err(c, err, 1)
	}
	ht, err := h.o.Find([]int32{int32(id)})
	if err != nil {
		h.l.Errorln("GetCommandTemplateById: ", err)
		return util.Err(c, err, 2)
	}
	h.l.Infoln("GetCommandTemplateById: success")
	return c.Status(200).JSON(e_command_template.Format(ht)[0])
}

// AddCommandTemplate swagger
// @Summary Create command templates
// @Tags    command_template
// @Accept  json
// @Produce json
// @Param   command_template body     []e_command_template.CommandTemplateCreate true "command template body"
// @Success 200           {array} e_command_template.CommandTemplate
// @Router  /api/command_template/ [post]
func (h *Handler) AddCommandTemplate(c *fiber.Ctx) error {
	entry := []*e_command_template.CommandTemplateCreate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("AddCommandTemplate: ", err)
		return util.Err(c, err, 0)
	}
	result, err := h.o.Create(entry)
	if err != nil {
		h.l.Errorln("AddCommandTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_command_template.Format(result))
}

// UpdateCommandTemplate swagger
// @Summary Update command templates
// @Tags    command_template
// @Accept  json
// @Produce json
// @Param   command_template body     []e_command_template.CommandTemplateUpdate true "modify command template body"
// @Success 200           {string} string "update successfully"
// @Router  /api/command_template/ [patch]
func (h *Handler) UpdateCommandTemplate(c *fiber.Ctx) error {
	entry := []*e_command_template.CommandTemplateUpdate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("UpdateCommandTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Update(entry)
	if err != nil {
		h.l.Errorln("UpdateCommandTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("update successfully")
}

// DeleteCommandTemplate swagger
// @Summary Delete command templates
// @Tags    command_template
// @Produce json
// @Param ids body []int true "command_template id"
// @Success 200 {string} string "delete successfully"
// @Router  /api/command_template/ [delete]
func (h *Handler) DeleteCommandTemplate(c *fiber.Ctx) error {
	entry := make([]int32, 0, 10)
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("DeleteCommandTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Delete(entry)
	if err != nil {
		h.l.Errorln("DeleteCommandTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("delete successfully")
}

// ExecuteCommand swagger
// @Summary execute command templates
// @Tags    command_template
// @Produce json
// @Param id path int true "command_template id"
// @Param   sendCommand body  command_template.SendCommand true "send command body"
// @Success 200 {string} string "command id"
// @Router  /api/command_template/execute/{id} [post]
func (h *Handler) ExecuteCommand(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Errorln("ExecuteCommand: ", err)
		return util.Err(c, err, 1)
	}
	entry := SendCommand{}
	if err = c.BodyParser(&entry); err != nil {
		h.l.Errorln("ExecuteCommand: ", err)
		return util.Err(c, err, 2)
	}
	st := e_command_template.SendCommandTemplate{
		TemplateId:     int32(id),
		Source:         "Command Template",
		TriggerFrom:    entry.TriggerFrom,
		TriggerAccount: entry.TriggerAccount,
		Token:          entry.Token,
		Variables:      entry.Variables,
	}
	commandId, err := h.o.Execute(c.UserContext(), st)
	if err != nil {
		return util.Err(c, err, 3)
	}
	return c.Status(200).JSON(commandId)
}

// SendCommandTemplate swagger
// @Summary Test command and return result
// @Tags    command_template
// @Accept  json
// @Produce json
// @Param   command_template body     e_command_template.CommandTemplateCreate true "send command template body"
// @Success 200           {object} e_command.Command
// @Router  /api/command_template/send/ [post]
func (h *Handler) SendCommandTemplate(c *fiber.Ctx) error {
	entry := e_command_template.CommandTemplateCreate{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("SendCommandTemplate: ", err)
		return util.Err(c, err, 0)
	}
	com := h.o.Send(c.UserContext(), entry)
	return c.Status(200).JSON(com)
}
