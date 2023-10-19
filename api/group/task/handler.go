package task

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
)

type hOperate interface {
	List() ([]e_task.Task, error)
	Find(ids []uint64) ([]e_task.Task, error)
	Cancel(id uint64, message string) error
	GetHistory(templateId, start, stop, status string) ([]e_task.TaskPub, error)
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

// GetTasks swagger
// @Summary     Show all tasks
// @Description Get all tasks
// @Tags        task
// @Produce     json
// @Success     200 {array} e_task.Task
// @Router      /api/task/ [get]
func (h *Handler) GetTasks(c *fiber.Ctx) error {
	tasks, err := h.o.List()
	if err != nil {
		h.l.Error().Println("Error getting tasks")
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task.ToPubSlice(tasks))
}

// GetTaskById swagger
// @Summary     Show task
// @Description Get task by id
// @Tags        task
// @Produce     json
// @Param       id  path     string true "id"
// @Success     200 {object} e_task.Task
// @Router      /api/task/{id} [get]
func (h *Handler) GetTaskById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := errors.New("id is error")
		h.l.Error().Println("GetTaskById: ", e)
		return util.Err(c, e, 0)
	}
	ht, err := h.o.Find([]uint64{uint64(id)})
	if err != nil {
		h.l.Error().Println("GetTaskById: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task.ToPub(ht[0]))
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
		h.l.Error().Println("GetTaskById: ", e)
		return util.Err(c, e, 0)
	}
	entry := CancelBody{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Error().Println("CancelTask: ", err)
		return util.Err(c, err, 0)
	}
	err = h.o.Cancel(uint64(id), entry.ClientMessage)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("cancel successful")
}

// GetHistory swagger
// @Summary get task history
// @Tags    task
// @Accept  json
// @Produce json
// @Param       template_id  query     int false "task template id"
// @Param       status  query     string false "status" Enums(Success, Failure, Cancel)
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Success 200 {array} e_task.Task
// @Router  /api/task/history [get]
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	templateId := c.Query("template_id")
	status := c.Query("status")
	start := c.Query("start")
	if start == "" {
		return util.Err(c, NoStartTime, 0)
	}
	stop := c.Query("stop")
	data, err := h.o.GetHistory(templateId, start, stop, status)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(data)
}
