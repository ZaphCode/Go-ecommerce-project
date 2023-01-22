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

func NewServer(
	userSvc domain.UserService,
	emailSvc email.EmailService,
	jwtSvc auth.JwtAuthService,
	validationSvc validation.ValidationService,
) Server {
	return &fiberServer{
		tasksCh:       make(chan func()),
		app:           fiberSetup(),
		userSvc:       userSvc,
		emailSvc:      emailSvc,
		jwtSvc:        jwtSvc,
		validationSvc: validationSvc,
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
	router := s.app.Group("/api")

	//* Test route
	router.Get("/", func(c *fiber.Ctx) error {
		s.tasksCh <- func() {
			time.Sleep(2 * time.Second)
			fmt.Println("Hello from background tasks")
		}
		return c.SendString("Hello world!")
	})

	{ //* Auth routes
		r := router.Group("/auth")
		r.Get("/:provider/url", s.getOAuthUrl)
		r.Get("/:provider/callback", s.signInWihOAuth)
		r.Get("/me", s.authRequired, s.getAuthUser)
		r.Get("/refresh", s.refreshToken)
		r.Get("/signout", s.signOut)
		r.Post("/signin", s.signIn)
		r.Post("/signup", s.signUp)
	}

	{ //* User routes
		r := router.Group("/user")
		r.Get("/all", s.authRequired, s.roleRequired(domain.ModeratorRole), s.getUsers)
		r.Get("/get/:id", s.authRequired, s.roleRequired(domain.ModeratorRole), s.getUser)
		r.Post("/create", s.authRequired, s.roleRequired(domain.AdminRole), s.createUser)
		r.Put("/update/:id", s.authRequired, s.roleRequired(domain.AdminRole), s.updateUser)
		r.Delete("/delete/:id", s.authRequired, s.roleRequired(domain.AdminRole), s.deleteUser)
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
