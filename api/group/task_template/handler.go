package task_template

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/dal/model"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task_template"
	"github.com/littlebluewhite/schedule_task_command/util"
)

type hOperate interface {
	List() ([]model.TaskTemplate, error)
	Find(ids []int32) ([]model.TaskTemplate, error)
	Create([]*e_task_template.TaskTemplateCreate) ([]model.TaskTemplate, error)
	Update([]*e_task_template.TaskTemplateUpdate) error
	Delete([]int32) error
	reloadCache() error
	Execute(ctx context.Context, st e_task_template.SendTaskTemplate) (id uint64, err error)
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

// GetTaskTemplates swagger
// @Summary     Show all task templates
// @Description Get all task templates
// @Tags        task_template
// @Produce     json
// @Success     200 {array} e_task_template.TaskTemplate
// @Router      /api/task_template/ [get]
func (h *Handler) GetTaskTemplates(c *fiber.Ctx) error {
	ht, err := h.o.List()
	if err != nil {
		h.l.Errorln("GetTaskTemplates: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Infoln("GetTaskTemplates: success")
	return c.Status(200).JSON(e_task_template.Format(ht))
}

// GetTaskTemplateById swagger
// @Summary     Show task templates
// @Description Get task templates by id
// @Tags        task_template
// @Produce     json
// @Param       id  path     int true "task template id"
// @Success     200 {object} e_task_template.TaskTemplate
// @Router      /api/task_template/{id} [get]
func (h *Handler) GetTaskTemplateById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Errorln("GetTaskTemplateById: ", err)
		return util.Err(c, err, 1)
	}
	ht, err := h.o.Find([]int32{int32(id)})
	if err != nil {
		h.l.Errorln("GetTaskTemplateById: ", err)
		return util.Err(c, err, 2)
	}
	h.l.Infoln("GetTaskTemplateById: success")
	return c.Status(200).JSON(e_task_template.Format(ht)[0])
}

// AddTaskTemplate swagger
// @Summary Create task templates
// @Tags    task_template
// @Accept  json
// @Produce json
// @Param   task_template body     []e_task_template.TaskTemplateCreate true "task template body"
// @Success 200           {array} e_task_template.TaskTemplate
// @Router  /api/task_template/ [post]
func (h *Handler) AddTaskTemplate(c *fiber.Ctx) error {
	entry := []*e_task_template.TaskTemplateCreate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("AddTaskTemplate: ", err)
		return util.Err(c, err, 0)
	}
	tt, err := h.o.Create(entry)
	if err != nil {
		h.l.Errorln("AddTaskTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task_template.Format(tt))
}

// UpdateTaskTemplate swagger
// @Summary Update task templates
// @Tags    task_template
// @Accept  json
// @Produce json
// @Param   task_template body     []e_task_template.TaskTemplateUpdate true "modify task template body"
// @Success 200           {string} string "update successfully"
// @Router  /api/task_template/ [patch]
func (h *Handler) UpdateTaskTemplate(c *fiber.Ctx) error {
	entry := []*e_task_template.TaskTemplateUpdate{nil}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("UpdateTaskTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Update(entry)
	if err != nil {
		h.l.Errorln("UpdateTaskTemplate: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("update successfully")
}

// DeleteTaskTemplate swagger
// @Summary Delete task templates
// @Tags    task_template
// @Produce json
// @Param ids body []int true "task_template id"
// @Success 200 {string} string "delete successfully"
// @Router  /api/task_template/ [delete]
func (h *Handler) DeleteTaskTemplate(c *fiber.Ctx) error {
	entry := make([]int32, 0, 10)
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("DeleteTaskTemplate: ", err)
		return util.Err(c, err, 0)
	}
	err := h.o.Delete(entry)
	if err != nil {
		h.l.Errorln("DeleteTaskTemplate: ", err)
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
			return util.Err(c, err, 10)
		}
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("delete successfully")
}

// ExecuteTask swagger
// @Summary execute task templates
// @Tags    task_template
// @Produce json
// @Param id path int true "task_template id"
// @Param   sendTask body  task_template.SendTask true "send task body"
// @Success 200 {string} string "task id"
// @Router  /api/task_template/execute/{id} [post]
func (h *Handler) ExecuteTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Errorln("ExecuteTask: ", err)
		return util.Err(c, err, 1)
	}
	entry := SendTask{}
	if err = c.BodyParser(&entry); err != nil {
		h.l.Errorln("ExecuteTask: ", err)
		return util.Err(c, err, 2)
	}
	st := e_task_template.SendTaskTemplate{
		TemplateId:     id,
		Source:         "Task Template",
		TriggerFrom:    entry.TriggerFrom,
		TriggerAccount: entry.TriggerAccount,
		Token:          entry.Token,
		Variables:      entry.Variables,
	}
	tid, err := h.o.Execute(c.UserContext(), st)
	if err != nil {
		return util.Err(c, err, 3)
	}
	return c.Status(200).JSON(tid)
}
