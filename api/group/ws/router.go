package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	log := logFile.NewLogFile("router", "websocket.log")
	app := g.GetApp()
	wm := g.GetWebsocketManager()

	ws := app.Group("/ws")
	ws.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	ws.Get("/:id", websocket.New(func(c *websocket.Conn) {
		wm.Register(2, c)
		defer func() {
			wm.Unregister(2, c)
		}()
		//log.Info().Println(c.Locals("allowed"))
		//log.Info().Println(c.Params("id"))
		//log.Info().Println(c.Query("v"))
		//log.Info().Println(c.NetConn())
		//log.Info().Println(c.Cookies("session"))
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Info().Println("read:", err)
				break
			}
			log.Info().Println("send: %s, type: %d", msg, mt)
		}
	}))
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
	GetWebsocketManager() api.WebsocketManager
}
