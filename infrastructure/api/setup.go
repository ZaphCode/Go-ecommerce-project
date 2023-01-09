package api

import "github.com/gofiber/fiber/v2"

func Setup() *fiber.App {
	app := fiber.New()

	CreateRoutes(app.Group("/api"))

	return app
}
