package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"log"
)

// Session is the middleware used to manage sessions
func Session(store session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return err
		}

		ctx.Locals("session", sess)
		return ctx.Next()
	}
}

// Authentication ensure that interested endpoints
// are protected against unauthenticated users
func Authentication(ctx *fiber.Ctx) error {
	sess, ok := ctx.Locals("session").(*session.Session)
	if !ok {
		log.Println("Error casting session from locals")
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	userID := sess.Get("userID")
	if userID == nil {
		err := ctx.SendStatus(fiber.StatusUnauthorized)
		if err != nil {
			return err
		}
		return ctx.Redirect("/")
	}
	return ctx.Next()
}
