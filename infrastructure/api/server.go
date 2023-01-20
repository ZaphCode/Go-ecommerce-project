package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/domain"
	"github.com/ZaphCode/clean-arch/infrastructure/services/auth"
	"github.com/ZaphCode/clean-arch/infrastructure/services/email"
	"github.com/ZaphCode/clean-arch/infrastructure/services/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//* Server

type Server interface {
	Start(port string) error
	InitBackgroundTaks()
	CreateRoutes()
	TryRoute(req *http.Request) (*http.Response, error)
}

//* Fiber implementation

type fiberServer struct {
	tasksCh chan func()
	app     *fiber.App

	// Services
	userSvc       domain.UserService
	emailSvc      email.EmailService
	jwtSvc        auth.JwtAuthService
	validationSvc validation.ValidationService
}

//* Constructor

func NewServer() Server {
	return &fiberServer{
		tasksCh: make(chan func()),
		app:     fiberSetup(),
	}
}

//* Methods

func (s *fiberServer) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *fiberServer) InitBackgroundTaks() {
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

func (s *fiberServer) CreateRoutes() {
	r := s.app.Group("/api")

	{ //* Auth Routes
		r = r.Group("/auth")
		r.Get("/:provider/url", s.getOAuthUrl)
		r.Get("/:provider/callback", s.signInWihOAuth)
		r.Get("/signout", s.signOut)
		r.Post("/signin", s.signIn)
		r.Post("/singup", s.signUp)
	}
}

func (s *fiberServer) TryRoute(req *http.Request) (*http.Response, error) {
	return s.app.Test(req)
}

//* Util

func fiberSetup() *fiber.App {
	cfg := config.Get()

	app := fiber.New()

	corsCfg := cors.Config{
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
		corsCfg.AllowOrigins = cfg.Api.ClientOrigin
	}

	app.Use(cors.New(corsCfg))

	app.Use(logger.New())

	return app
}
