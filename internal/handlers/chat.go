package handlers

import (
	"github.com/Sergio-dot/open-call/pkg/chat"
	w "github.com/Sergio-dot/open-call/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/websocket/v2"
)

// RoomChat renders the HTML template for the chat
func RoomChat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map{}, "layouts/main")
}

// RoomChatWebsocket handles Websocket connections for chat
func RoomChatWebsocket(c *websocket.Conn) {
	sess := c.Locals("session").(*session.Session)

	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	w.RoomsLock.Lock()
	room := w.Rooms[uuid]
	w.RoomsLock.Unlock()
	if room == nil {
		return
	}
	if room.Hub == nil {
		return
	}
	chat.PeerChatConn(c.Conn, room.Hub, sess) // TODO - pass session
}

// StreamChatWebsocket handles Websocket connections for stream
func StreamChatWebsocket(c *websocket.Conn) {
	sess := c.Locals("session").(*session.Session)

	suuid := c.Params("suuid")
	if suuid == "" {
		return
	}

	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		if stream.Hub == nil {
			hub := chat.NewHub()
			stream.Hub = hub
			go hub.Run()
		}
		chat.PeerChatConn(c.Conn, stream.Hub, sess) // TODO - pass session
		return
	}
	w.RoomsLock.Unlock()
}
