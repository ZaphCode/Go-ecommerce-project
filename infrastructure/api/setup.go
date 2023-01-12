package api

import (
	"strings"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var cfg = config.GetConfig()

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Api.ClientOrigin,
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

	CreateRoutes(app.Group("/api"))

	return app
}
