package handlers

import (
	"fmt"
	"log"
	"os"
	"time"

	w "github.com/Sergio-dot/open-call/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Stream(c *fiber.Ctx) error {
	suuid := c.Params("suuid")
	if suuid == "" {
		c.Status(400)
		return nil
	}

	sess, err := Store.Get(c)
	if err != nil {
		log.Println(err)
		return c.Redirect("/")
	}

	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		ws = "wss"
	}

	w.RoomsLock.Lock()
	if _, ok := w.Streams[suuid]; ok {
		sess.Set("success-message", "Joined room")
		successMessage, _ := sess.Get("success-message").(string)
		// remove message from the session
		sess.Delete("success-message")
		sess.Delete("error-message")
		sess.Save()
		w.RoomsLock.Unlock()
		return c.Render("stream", fiber.Map{
			"StreamWebsocketAddr": fmt.Sprintf("%s://%s/stream/%s/websocket", ws, c.Hostname(), suuid),
			"ChatWebsocketAddr":   fmt.Sprintf("%s://%s/stream/%s/chat/websocket", ws, c.Hostname(), suuid),
			"ViewerWebsocketAddr": fmt.Sprintf("%s://%s/stream/%s/viewer/websocket", ws, c.Hostname(), suuid),
			"Type":                "stream",
			"PageTitle":           "OpenCall - Viewer",
			"UserID":              sess.Get("userID"),
			"Email":               sess.Get("email"),
			"Username":            sess.Get("username"),
			"CreatedAt":           sess.Get("createdAt"),
			"UpdatedAt":           sess.Get("updatedAt"),
			"ToastSuccess":        successMessage,
		}, "layouts/main")
	}

	sess.Set("error-message", "Room does not exist")
	sess.Save()

	w.RoomsLock.Unlock()

	return c.Redirect("/dashboard")
}

func StreamWebsocket(c *websocket.Conn) {
	suuid := c.Params("suuid")
	if suuid == "" {
		return
	}

	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		w.StreamConn(c, stream.Peers)
		return
	}
	w.RoomsLock.Unlock()
}

func StreamViewerWebsocket(c *websocket.Conn) {
	suuid := c.Params("suuid")
	if suuid == "" {
		return
	}

	w.RoomsLock.Lock()
	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		viewerConn(c, stream.Peers)
		return
	}
	w.RoomsLock.Unlock()
}

func viewerConn(c *websocket.Conn, p *w.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case <-ticker.C:
			nw, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			nw.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}
