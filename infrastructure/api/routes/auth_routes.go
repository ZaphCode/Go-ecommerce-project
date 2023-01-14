package routes

import (
	"github.com/ZaphCode/clean-arch/config"
	authCtroller "github.com/ZaphCode/clean-arch/infrastructure/controllers/auth"
	userRepo "github.com/ZaphCode/clean-arch/infrastructure/repositories/user"
	authSvc "github.com/ZaphCode/clean-arch/infrastructure/services/auth"
	userSvc "github.com/ZaphCode/clean-arch/infrastructure/services/user"
	validationSvc "github.com/ZaphCode/clean-arch/infrastructure/services/validation"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateAuthRoutes(router fiber.Router) {
	userRepository := userRepo.NewFirestoreUserRepository(
		utils.GetFirestoreClient(config.GetFirebaseApp()), "users",
	)
	userService := userSvc.NewUserService(userRepository)
	validationService := validationSvc.NewValidationService()
	jwtAuthService := authSvc.NewJwtAuthService()
	authController := authCtroller.NewAuthController(
		userService,
		validationService,
		jwtAuthService,
	)

	router.Post("/signin", authController.SignIn)
	router.Post("/signup", authController.SignUp)
	router.Get("/signout", authController.SignOut)
	router.Get("/:provider/callback", authController.SignInWihOAuth)
	router.Get("/:provider/url", authController.GetOAuthUrl)
}
