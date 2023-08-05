package handlers

import (
	"github.com/ahmad-khatib0/go/websockets/video-call-chat/pkg/chat"
	w "github.com/ahmad-khatib0/go/websockets/video-call-chat/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RoomChat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map{}, "layouts/main")
}

func RoomChatWebsocket(c *websocket.Conn) {
	uid := c.Params("uuid")
	if uid == "" {
		return
	}

	w.RoomsLock.Lock()
	room := w.Rooms[uid]
	w.RoomsLock.Unlock()

	if room == nil {
		return
	}

	if room.Hub == nil {
		return
	}

	chat.PeerChatConn(c.Conn, room.Hub)
}

func StreamChatWebsocket(c *websocket.Conn) {
	suid := c.Params("uuid")
	if suid == "" {
		return
	}

	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suid]; ok {
		w.RoomsLock.Unlock()
		if stream.Hub == nil {
			hub := chat.NewHub()
			stream.Hub = hub
			go hub.Run()
		}
		chat.PeerChatConn(c.Conn, stream.Hub)
		return
	}

	w.Rooms.Unlock()
}
