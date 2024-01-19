package websocket_hub

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"schedule_task_command/util/logFile"
)

type client struct {
	conn *websocket.Conn
	box  chan []byte
	l    logFile.LogFile
}

func newClient(conn *websocket.Conn, log logFile.LogFile) *client {
	return &client{
		conn: conn,
		box:  make(chan []byte, 256),
		l:    log,
	}
}

func (c *client) readPump() {
	for {
		if _, msg, err := c.conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.l.Error().Println("reade err:", err)
			}
			break
		} else {
			c.l.Info().Printf("recv: %s", msg)
			// do some command
		}
	}
}

func (c *client) writePump(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.box:
			c.l.Info().Printf("client: %v send 2 start", c.conn.RemoteAddr())
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.l.Error().Println("NextWriter err:", err)
				return
			}
			n, err := w.Write(msg)
			if err != nil {
				c.l.Error().Println("w.Write err:", err)
				return
			}
			if n != len(msg) {
				c.l.Error().Println("w.Write length err:", n)
				return
			}
			if err = w.Close(); err != nil {
				c.l.Error().Println("w.close err:", err)
				return
			}
			c.l.Info().Printf("client: %v send 2 end", c.conn.RemoteAddr())
		}
	}
}

func (c *client) close() {
	_ = c.conn.Close()
}

func (c *client) send(msg []byte) {
	c.l.Info().Printf("client: %v send 1", c.conn.RemoteAddr())
	c.box <- msg
}
