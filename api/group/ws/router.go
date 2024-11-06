package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/entry/e_module"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	log := my_log.NewLog("router/websocket")
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
			log.Errorln(err)
		}
	}))
	ws.Get("/task", websocket.New(func(c *websocket.Conn) {
		err := hm.WsConnect(e_module.Task, c)
		if err != nil {
			log.Errorln(err)
		}
	}))
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
	GetWebsocketHub() api.HubManager
}
