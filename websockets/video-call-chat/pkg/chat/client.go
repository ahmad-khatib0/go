package chat

import (
	"bytes"
	"log"
	"time"

	"github.com/fasthttp/websocket"
)

var (
	newLine = []byte("\n")
	space   = []byte(" ")
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error websocket connection:  %v", err)
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newLine, space, -1))
		c.Hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}
			wr, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			wr.Write(msg)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				wr.Write(newLine)
				wr.Write(<-c.Send)
			}

			if err := wr.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func PeerChatConn(c *websocket.Conn, hub *Hub) {
	client := &Client{Hub: hub, Conn: c, Send: make(chan []byte, 256)}
	client.Hub.register <- client

	go client.writePump()
	client.readPump()
}
