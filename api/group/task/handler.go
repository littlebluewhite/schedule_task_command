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
	Find(taskIds []string) ([]e_task.Task, error)
	Cancel(taskId string) error
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
	}
	return c.Status(200).JSON(e_task.ToPubSlice(tasks))
}

// GetTaskByTaskId swagger
// @Summary     Show task
// @Description Get task by taskId
// @Tags        task
// @Produce     json
// @Param       taskId  path     string true "taskId"
// @Success     200 {object} e_task.Task
// @Router      /api/task/{taskId} [get]
func (h *Handler) GetTaskByTaskId(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	if taskId == "" {
		e := errors.New("taskId is error")
		h.l.Error().Println("GetTaskByTaskId: ", e)
		return util.Err(c, e, 0)
	}
	ht, err := h.o.Find([]string{taskId})
	if err != nil {
		h.l.Error().Println("GetTaskByTaskId: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_task.ToPub(ht[0]))
}

// CancelTask swagger
// @Summary     cancel task
// @Description Cancel task by taskId
// @Tags        task
// @Produce     json
// @Param       taskId  path     string true "taskId"
// @Success     200 {string} cancel successfully
// @Router      /api/task/{taskId} [Delete]
func (h *Handler) CancelTask(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	if taskId == "" {
		e := errors.New("taskId is error")
		h.l.Error().Println("GetTaskByTaskId: ", e)
		return util.Err(c, e, 0)
	}
	err := h.o.Cancel(taskId)
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
