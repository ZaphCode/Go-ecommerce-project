package api

import (
	"strings"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/infrastructure/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup() *fiber.App {
	cfg := config.Get()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		//AllowOrigins: cfg.Api.ClientOrigin,
		AllowOrigins: "*",
		AllowHeaders: strings.Join([]string{
			"Origin",
			"Content-Type",
			"Accept",
			cfg.Api.RefreshTokenHeader,
			cfg.Api.AccessTokenHeader,
		}, ", "),
		AllowMethods:     cors.ConfigDefault.AllowMethods,
		AllowCredentials: true,
	}))

	api := app.Group("/api")

	// Health check
	api.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Hello world")
	})

	// Routes
	routes.CreateAuthRoutes(api.Group("/auth"))

	return app
}
