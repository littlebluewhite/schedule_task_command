package task

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task"
	"github.com/littlebluewhite/schedule_task_command/util"
)

type hOperate interface {
	List() ([]e_task.Task, error)
	Find(ids []uint64) ([]e_task.Task, error)
	Cancel(id uint64, message string) error
	GetHistory(id, templateId, start, stop, status string) ([]e_task.TaskPub, error)
	FindById(id uint64) (t e_task.TaskPub, err error)
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

// GetTasks swagger
// @Summary     Show all tasks
// @Description Get all tasks
// @Tags        task
// @Produce     json
// @Success     200 {array} e_task.TaskPub
// @Router      /api/task/ [get]
func (h *Handler) GetTasks(c *fiber.Ctx) error {
	tasks, err := h.o.List()
	if err != nil {
		h.l.Errorln("Error getting tasks")
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task.ToPubSlice(tasks))
}

// GetTaskById swagger
// @Summary     Show task
// @Description Get task by id include history data
// @Tags        task
// @Produce     json
// @Param       id  path     string true "id"
// @Success     200 {object} e_task.Task
// @Router      /api/task/{id} [get]
func (h *Handler) GetTaskById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := errors.New("id is error")
		h.l.Errorln("GetTaskById: ", e)
		return util.Err(c, e, 1)
	}
	t, err := h.o.FindById(uint64(id))
	if err != nil {
		h.l.Errorln("GetTaskById: ", err)
		return util.Err(c, err, 2)
	}
	return c.Status(200).JSON(t)
}

// GetSimpleTasks swagger
// @Summary     Show all simple tasks
// @Description Get all simple tasks
// @Tags        task
// @Produce     json
// @Success     200 {array} e_task.SimpleTask
// @Router      /api/task/simple/ [get]
func (h *Handler) GetSimpleTasks(c *fiber.Ctx) error {
	tasks, err := h.o.List()
	if err != nil {
		h.l.Errorln("Error getting tasks")
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task.ToSimpleTaskSlice(tasks))
}

// GetSimpleTasksById swagger
// @Summary     Show task
// @Description Get simple task by id include history data
// @Tags        task
// @Produce     json
// @Param       id  path     string true "id"
// @Success     200 {object} e_task.SimpleTask
// @Router      /api/task/simple/{id} [get]
func (h *Handler) GetSimpleTasksById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := errors.New("id is error")
		h.l.Errorln("GetTaskById: ", e)
		return util.Err(c, e, 0)
	}
	t, err := h.o.FindById(uint64(id))
	if err != nil {
		h.l.Errorln("GetTaskById: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task.TaskPubToSimpleTask(t))
}

// GetStageItemStatus swagger
// @Summary     Show stage item Status
// @Description Get stage item Status by task id
// @Tags        task
// @Produce     json
// @Param       id  path     string true "id"
// @Success     200 {array} int
// @Router      /api/task/stage_item/status/{id} [get]
func (h *Handler) GetStageItemStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := errors.New("id is error")
		h.l.Errorln("GetTaskById: ", e)
		return util.Err(c, e, 0)
	}
	t, err := h.o.FindById(uint64(id))
	if err != nil {
		h.l.Errorln("GetTaskById: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task.ToStageItemStatus(t))
}

// GetHistory swagger
// @Summary get task history
// @Tags    task
// @Accept  json
// @Produce json
// @Param       id  query     int false "task id"
// @Param       template_id  query     int false "task template id"
// @Param       status  query     string false "status" Enums(Success, Failure, Cancel)
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Success 200 {array} e_task.Task
// @Router  /api/task/history [get]
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	id := c.Query("id")
	templateId := c.Query("template_id")
	status := c.Query("status")
	start := c.Query("start")
	if start == "" {
		return util.Err(c, NoStartTime, 0)
	}
	stop := c.Query("stop")
	data, err := h.o.GetHistory(id, templateId, start, stop, status)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(data)
}

// CancelTask swagger
// @Summary     cancel task
// @Description Cancel task by id
// @Tags        task
// @Produce     json
// @Param       id  path     string true "id"
// @Param   message body     task.CancelBody true "cancel message"
// @Success     200 {string} cancel successfully
// @Router      /api/task/{id} [Delete]
func (h *Handler) CancelTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := errors.New("id is error")
		h.l.Errorln("GetTaskById: ", e)
		return util.Err(c, e, 0)
	}
	entry := CancelBody{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("CancelTask: ", err)
		return util.Err(c, err, 0)
	}
	err = h.o.Cancel(uint64(id), entry.ClientMessage)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("cancel successful")
}
