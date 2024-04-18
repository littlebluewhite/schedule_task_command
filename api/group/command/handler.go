package command

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
)

type hOperate interface {
	List() ([]e_command.Command, error)
	Find(ids []uint64) ([]e_command.Command, error)
	Cancel(id uint64, message string) error
	GetHistory(id, templateId, status, start, stop string) ([]e_command.CommandPub, error)
	FindById(id uint64) (e_command.CommandPub, error)
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

// GetCommands swagger
// @Summary     Show all commands
// @Description Get all commands
// @Tags        Command
// @Produce     json
// @Success     200 {array} e_command.Command
// @Router      /api/command/ [get]
func (h *Handler) GetCommands(c *fiber.Ctx) error {
	commands, err := h.o.List()
	if err != nil {
		h.l.Error().Println("Error getting commands")
	}
	return c.Status(200).JSON(e_command.ToPubSlice(commands))
}

// GetCommandById swagger
// @Summary     Show Command
// @Description Get Command by id include history data
// @Tags        Command
// @Produce     json
// @Param       id  path     string true "id"
// @Success     200 {object} e_command.Command
// @Router      /api/command/{id} [get]
func (h *Handler) GetCommandById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := errors.New("id is error")
		h.l.Error().Println("GetCommandById: ", e)
		return util.Err(c, e, 1)
	}
	com, err := h.o.FindById(uint64(id))
	if err != nil {
		h.l.Error().Println("GetCommandById: ", err)
		return util.Err(c, err, 2)
	}
	b, e := json.Marshal(com)
	fmt.Println(e)
	fmt.Println(b)
	return c.Status(200).JSON(com)
}

// CancelCommand swagger
// @Summary     cancel Command
// @Description Cancel Command by id
// @Tags        Command
// @Produce     json
// @Param       id  path     string true "id"
// @Param   message body     command.CancelBody true "cancel message"
// @Success     200 {string} cancel successfully
// @Router      /api/command/{id} [Delete]
func (h *Handler) CancelCommand(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		e := errors.New("id is error")
		h.l.Error().Println("GetCommandById: ", e)
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
// @Summary get Command history
// @Tags    Command
// @Accept  json
// @Produce json
// @Param       id  query     int false "Command id"
// @Param       template_id  query     int false "Command template id"
// @Param       status  query     string false "status" Enums(Success, Failure, Cancel)
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Success 200 {array} e_command.Command
// @Router  /api/command/history [get]
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
