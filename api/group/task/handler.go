package task

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
)

type hOperate interface {
	List() ([]model.TaskTemplate, error)
	Find(ids []int32) ([]model.TaskTemplate, error)
	Create([]*e_task_template.TaskTemplateCreate) ([]model.TaskTemplate, error)
	Update([]*e_task_template.TaskTemplateUpdate) error
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

// GetTasks swagger
// @Summary     Show all tasks
// @Description Get all tasks
// @Tags        task
// @Produce     json
// @Success     200 {array} e_task.Task
// @Router      /api/task/ [get]
func (h *Handler) GetTasks(c *fiber.Ctx) error {
	ht, err := h.o.List()
	if err != nil {
		h.l.Error().Println("GetTaskTemplates: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Info().Println("GetTaskTemplates: success")
	return c.Status(200).JSON(e_task_template.Format(ht))
}

// GetTaskByTaskId swagger
// @Summary     Show task
// @Description Get task by taskId
// @Tags        task
// @Produce     json
// @Param       taskId  path     int true "taskId"
// @Success     200 {object} e_task.Task
// @Router      /task/{taskId} [get]
func (h *Handler) GetTaskByTaskId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Error().Println("GetTaskTemplateById: ", err)
		return util.Err(c, err, 0)
	}
	ht, err := h.o.Find([]int32{int32(id)})
	if err != nil {
		h.l.Error().Println("GetTaskTemplateById: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Info().Println("GetTaskTemplateById: success")
	return c.Status(200).JSON(e_task_template.Format(ht)[0])
}

// CancelTask swagger
// @Summary     cancel task
// @Description Cancel task by taskId
// @Tags        task
// @Produce     json
// @Param       taskId  path     int true "taskId"
// @Success     200 {string} cancel successfully
// @Router      /task/{taskId} [Delete]
func (h *Handler) CancelTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.l.Error().Println("GetTaskTemplateById: ", err)
		return util.Err(c, err, 0)
	}
	ht, err := h.o.Find([]int32{int32(id)})
	if err != nil {
		h.l.Error().Println("GetTaskTemplateById: ", err)
		return util.Err(c, err, 0)
	}
	h.l.Info().Println("GetTaskTemplateById: success")
	return c.Status(200).JSON(e_task_template.Format(ht)[0])
}

// GetHistory swagger
// @Summary get task history
// @Tags    task
// @Accept  json
// @Produce json
// @Param       id  path     int true "time template id"
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Success 200 {array} e_task.task
// @Router  /api/task/history/{id} [get]
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	//id := c.Params("id")
	//start := c.Query("start")
	//if start == "" {
	//	return util.Err(c, NoStartTime, 0)
	//}
	//stop := c.Query("stop")
	//data, err := h.o.ReadFromHistory(id, start, stop)
	//if err != nil {
	//	return util.Err(c, err, 0)
	//}
	//return c.Status(200).JSON(data)
	return nil
}
