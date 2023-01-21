package main

import (
	"log"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/infrastructure/api"
	"github.com/ZaphCode/clean-arch/infrastructure/repositories/user"
	"github.com/ZaphCode/clean-arch/infrastructure/services/auth"
	"github.com/ZaphCode/clean-arch/infrastructure/services/core"
	"github.com/ZaphCode/clean-arch/infrastructure/services/email"
	"github.com/ZaphCode/clean-arch/infrastructure/services/validation"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
)

func init() {
	config.MustLoadConfig("./config")
	config.MustLoadFirebaseConfig("./config")
}

func main() {
	cfg := config.Get()
	client := utils.GetFirestoreClient(config.GetFirebaseApp())
	// utils.PrettyPrint(cfg)

	//* Repos
	userRepo := user.NewFirestoreUserRepository(client, "users")

	//* Services
	userSvc := core.NewUserService(userRepo)
	emailSvc := email.NewSmtpEmailService()
	validationSvc := validation.NewValidationService()
	jwtAuthSvc := auth.NewJwtAuthService()

	server := api.NewServer(
		userSvc,
		emailSvc,
		jwtAuthSvc,
		validationSvc,
	)

	server.CreateRoutes()

	go server.InitBackgroundTaks()

	if err := server.Start(":" + cfg.Api.Port); err != nil {
		log.Fatal(err)
	}
}
