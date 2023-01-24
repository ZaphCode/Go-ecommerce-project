package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	tasksCh chan func()
	app     *fiber.App
}

//* Constructor

func New() *Server {
	return &Server{
		tasksCh: make(chan func()),
		app:     fiber.New(),
	}
}

//* Methods

func (s *Server) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) InitBackgroundTaks() {
	ticker := time.NewTicker(3 * time.Hour)
	for {
		select {
		case task := <-s.tasksCh:
			task()
		case <-ticker.C:
			fmt.Println("Tick")
		}
	}
}

func (s *Server) SetGlobalMiddlewares() {
	cfg := config.Get()

	cc := cors.Config{
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
	}

	if config.IsProduction() {
		cc.AllowOrigins = cfg.Api.ClientOrigin
	}

	s.app.Use(cors.New(cc))
	s.app.Use(recover.New())
	s.app.Use(logger.New())

	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"healthy": true,
			"message": "Hello world!",
		})
	})
}

func (s *Server) TryRoute(req *http.Request) (*http.Response, error) {
	return s.app.Test(req)
}
