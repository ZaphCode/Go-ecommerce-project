package api

import "github.com/gofiber/fiber/v2"

func CreateRoutes(router fiber.Router) {

	router.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Hello world")
	})
}
