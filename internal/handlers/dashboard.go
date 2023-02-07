package handlers

import "github.com/gofiber/fiber/v2"

func Dashboard(ctx *fiber.Ctx) error {
	return ctx.Render("dashboard", fiber.Map{
		"PageTitle": "OpenCall - Home",
	}, "layouts/main")
}
