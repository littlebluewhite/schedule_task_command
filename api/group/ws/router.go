package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_module"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	log := logFile.NewLogFile("router", "websocket.log")
	app := g.GetApp()

	hm := g.GetWebsocketHub()
	hm.RegisterHub(e_module.Task)
	hm.RegisterHub(e_module.Command)

	ws := app.Group("/ws")

	ws.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	ws.Get("/command", websocket.New(func(c *websocket.Conn) {
		err := hm.WsConnect(e_module.Command, c)
		if err != nil {
			log.Error().Println(err)
		}
	}))
	ws.Get("/task", websocket.New(func(c *websocket.Conn) {
		err := hm.WsConnect(e_module.Task, c)
		if err != nil {
			log.Error().Println(err)
		}
	}))
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
	GetWebsocketHub() api.HubManager
}
