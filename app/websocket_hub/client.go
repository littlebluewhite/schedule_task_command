package websocket_hub

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"schedule_task_command/api"
)

type client struct {
	conn *websocket.Conn
	box  chan []byte
	l    api.Logger
}

func newClient(conn *websocket.Conn, log api.Logger) *client {
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
				c.l.Errorln("reade err:", err)
			}
			break
		} else {
			c.l.Infof("recv: %s", msg)
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
			c.l.Infof("client: %v send 2 start", c.conn.RemoteAddr())
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.l.Errorln("NextWriter err:", err)
				return
			}
			n, err := w.Write(msg)
			if err != nil {
				c.l.Errorln("w.Write err:", err)
				return
			}
			if n != len(msg) {
				c.l.Errorln("w.Write length err:", n)
				return
			}
			if err = w.Close(); err != nil {
				c.l.Errorln("w.close err:", err)
				return
			}
			c.l.Infof("client: %v send 2 end", c.conn.RemoteAddr())
		}
	}
}

func (c *client) close() {
	_ = c.conn.Close()
}

func (c *client) send(msg []byte) {
	c.l.Infof("client: %v send 1", c.conn.RemoteAddr())
	c.box <- msg
}
