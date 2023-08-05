package handlers

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/ahmad-khatib0/go/websockets/video-call-chat/pkg/chat"
	w "github.com/ahmad-khatib0/go/websockets/video-call-chat/pkg/webrtc"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)

func RoomCreate(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("/room/%s", uuid.New().String()))
}

func Room(c *fiber.Ctx) error {
	uid := c.Params("uuid")
	if uid == "" {
		c.Status(400)
		return nil
	}

	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		ws = "wss"
	}

	uuid1, suuid, _ := createOrGetRoom(uid)

	return c.Render("peer", fiber.Map{
		"RoomWebsocketAddr":   fmt.Sprintf("%s://%s/room/%s/websocket", ws, c.Hostname(), uuid1),
		"RoomLink":            fmt.Sprintf("%s://%s/room/%s", c.Protocol(), c.Hostname(), uuid1),
		"ChatWebsocketAddr":   fmt.Sprintf("%s://%s/room/%s/chat/websocket", ws, c.Hostname(), uuid1),
		"ViewerWebsocketAddr": fmt.Sprintf("%s://%s/room/%s/viewer/websocket", ws, c.Hostname(), uuid1),
		"StreamLink":          fmt.Sprintf("%s://%s/stream/%s", c.Protocol(), c.Hostname(), suuid),
		"Type":                "room",
	}, "layouts/main")
}

func RoomWebsocket(c *websocket.Conn) {
	uid := c.Params("uuid")
	if uid == "" {
		return
	}

	_, _, room := createOrGetRoom(uid)
	w.RoomConn(c, room.Peers)
}

func createOrGetRoom(uuid string) (string, string, *w.Room) {
	w.RoomsLock.Lock()
	defer w.RoomsLock.Unlock()

	h := sha256.New()
	h.Write([]byte(uuid))
	suuid := fmt.Sprintf("%x", h.Sum(nil))

	if room := w.Rooms[uuid]; room != nil {
		if _, ok := w.Streams[suuid]; ok {
			w.Streams[suuid] = room
		}
		return uuid, suuid, room
	}

	hub := chat.NewHub()
	p := &w.Peers{}
	p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &w.Room{
		Peers: p,
		Hub:   hub,
	}

	w.Rooms[uuid] = room
	w.Streams[suuid] = room
	go hub.Run()

	return uuid, suuid, room
}

func RoomViewerWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	w.RoomsLock.Lock()
	if peer, ok := w.Rooms[uuid]; ok {
		w.RoomsLock.Unlock()
		roomViewerConn(c, peer.Peers)
		return
	}
	w.RoomsLock.Unlock()
}

func roomViewerConn(c *websocket.Conn, p *w.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case <-ticker.C:
			wr, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			wr.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
