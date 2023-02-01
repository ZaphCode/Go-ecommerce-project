package main

import (
	"log"

	"github.com/ZaphCode/clean-arch/config"
	_ "github.com/ZaphCode/clean-arch/docs"
	"github.com/ZaphCode/clean-arch/src/api"
	"github.com/ZaphCode/clean-arch/src/api/handlers"
	"github.com/ZaphCode/clean-arch/src/api/middlewares"
	"github.com/ZaphCode/clean-arch/src/repositories/category"
	"github.com/ZaphCode/clean-arch/src/repositories/product"
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

// @title           Swagger Example API
// @version         1.6.9
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
//// @contact.name   API Support
//// @contact.url    http://www.swagger.io/support
//// @contact.email  support@swagger.io

//// @license.name  Apache 2.0
//// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name X-Access-Token

// @host      localhost:9000
// @BasePath  /api
func main() {
	cfg := config.Get()
	client := utils.GetFirestoreClient(config.GetFirebaseApp())

	//* Repos
	userRepo := user.NewFirestoreUserRepository(client, utils.UserColl)
	prodRepo := product.NewFirestoreProductRepository(client, utils.ProdColl)
	catRepo := category.NewFirestoreCategoryRepository(client, utils.CategColl)

	//* Services
	userSvc := core.NewUserService(userRepo)
	prodSvc := core.NewProductService(prodRepo)
	catSvc := core.NewCategoryService(catRepo)
	emailSvc := email.NewSmtpEmailService()
	vldSvc := validation.NewValidationService()
	jwtSvc := auth.NewJWTService()

	//* Midlewares
	authMdlw := middlewares.NewAuthMiddleware(jwtSvc)

	//* Handlers
	usrHdlr := handlers.NewUserHandler(userSvc, vldSvc)
	authHdlr := handlers.NewAuthHandler(userSvc, emailSvc, jwtSvc, vldSvc)
	prodHdlr := handlers.NewProdutHandler(prodSvc, catSvc, vldSvc)
	catHdlr := handlers.NewCategoryHandler(prodSvc, catSvc, vldSvc)

	//* Server
	server := api.New()

	server.SetGlobalMiddlewares()

	server.CreateAuthRoutes(authHdlr, authMdlw)
	server.CreateUserRoutes(usrHdlr, authMdlw)
	server.CreateProductRoutes(prodHdlr, authMdlw)
	server.CreateCategoryRoutes(catHdlr, authMdlw)

	go server.InitBackgroundTaks()

	if err := server.Start(":" + cfg.Api.Port); err != nil {
		log.Fatal(err)
	}
}
