package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	swagger "github.com/arsmn/fiber-swagger/v2"
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

func (s *Server) InitBackgroundTasks() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTSTP)
	docsTicker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case task := <-s.tasksCh:
			task()
		case <-docsTicker.C:
			fmt.Println("check docs on: http://localhost:9000/docs/index.html")
			docsTicker.Stop()
		case <-signalCh:
			fmt.Println("Shutting down the server...")
			close(s.tasksCh)
			if err := s.app.Shutdown(); err != nil {
				fmt.Println("Error shutting down the server")
			}
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

	s.app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"healthy": true,
			"message": "Hello world!",
		})
	})

	s.app.Get("/docs/*", swagger.HandlerDefault)
}

func (s *Server) TryRoute(req *http.Request) (*http.Response, error) {
	return s.app.Test(req, -1) // "WITHOUT TIMEOUT"
}
