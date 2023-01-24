package main

import (
	"log"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api"
	"github.com/ZaphCode/clean-arch/src/api/handlers"
	"github.com/ZaphCode/clean-arch/src/api/middlewares"
	"github.com/ZaphCode/clean-arch/src/repositories/user"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/core"
	"github.com/ZaphCode/clean-arch/src/services/email"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/ZaphCode/clean-arch/src/utils"
)

func init() {
	config.MustLoadConfig("./config")
	config.MustLoadFirebaseConfig("./config")
}

func main() {
	cfg := config.Get()
	client := utils.GetFirestoreClient(config.GetFirebaseApp())

	//* Repos
	userRepo := user.NewFirestoreUserRepository(client, utils.UserColl)

	//* Services
	userSvc := core.NewUserService(userRepo)
	emailSvc := email.NewSmtpEmailService()
	vldSvc := validation.NewValidationService()
	jwtSvc := auth.NewJWTService()

	//* Midlewares
	authMdlw := middlewares.NewAuthMiddleware(jwtSvc)

	//* Handlers
	usrHdlr := handlers.NewUserHandler(userSvc, vldSvc)
	authHdlr := handlers.NewAuthHandler(userSvc, emailSvc, jwtSvc, vldSvc)

	//* Server
	server := api.New()

	server.SetGlobalMiddlewares()

	server.CreateAuthRoutes(authHdlr, authMdlw)
	server.CreateUserRoutes(usrHdlr, authMdlw)

	go server.InitBackgroundTaks()

	if err := server.Start(":" + cfg.Api.Port); err != nil {
		log.Fatal(err)
	}
}
