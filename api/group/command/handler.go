package command

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/util"
	"schedule_task_command/util/logFile"
)

type hOperate interface {
	List() ([]e_command.Command, error)
	Find(commandIds []string) ([]e_command.Command, error)
	Cancel(commandId string) error
	GetHistory(templateId, status, start, stop string) ([]e_command.CommandPub, error)
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

// GetCommandByCommandId swagger
// @Summary     Show Command
// @Description Get Command by commandId
// @Tags        Command
// @Produce     json
// @Param       commandId  path     string true "commandId"
// @Success     200 {object} e_command.Command
// @Router      /api/command/{commandId} [get]
func (h *Handler) GetCommandByCommandId(c *fiber.Ctx) error {
	commandId := c.Params("commandId")
	if commandId == "" {
		e := errors.New("commandId is error")
		h.l.Error().Println("GetCommandByCommandId: ", e)
		return util.Err(c, e, 0)
	}
	ht, err := h.o.Find([]string{commandId})
	if err != nil {
		h.l.Error().Println("GetCommandByCommandId: ", err)
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(e_command.ToPub(ht[0]))
}

// CancelCommand swagger
// @Summary     cancel Command
// @Description Cancel Command by commandId
// @Tags        Command
// @Produce     json
// @Param       commandId  path     string true "commandId"
// @Success     200 {string} cancel successfully
// @Router      /api/command/{commandId} [Delete]
func (h *Handler) CancelCommand(c *fiber.Ctx) error {
	commandId := c.Params("commandId")
	if commandId == "" {
		e := errors.New("commandId is error")
		h.l.Error().Println("GetCommandByCommandId: ", e)
		return util.Err(c, e, 0)
	}
	err := h.o.Cancel(commandId)
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
// @Param       template_id  query     int false "Command template id"
// @Param       status  query     string false "status" Enums(Success, Failure, Cancel)
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Success 200 {array} e_command.Command
// @Router  /api/command/history [get]
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
