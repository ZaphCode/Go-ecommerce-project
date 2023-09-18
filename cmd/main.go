package main

import (
	"log"

	"github.com/ZaphCode/clean-arch/config"
	_ "github.com/ZaphCode/clean-arch/docs"
	"github.com/ZaphCode/clean-arch/src/api"
	"github.com/ZaphCode/clean-arch/src/api/handlers"
	"github.com/ZaphCode/clean-arch/src/api/middlewares"
	"github.com/ZaphCode/clean-arch/src/repositories/address"
	"github.com/ZaphCode/clean-arch/src/repositories/category"
	"github.com/ZaphCode/clean-arch/src/repositories/order"
	"github.com/ZaphCode/clean-arch/src/repositories/product"
	"github.com/ZaphCode/clean-arch/src/repositories/user"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/core"
	"github.com/ZaphCode/clean-arch/src/services/email"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/ZaphCode/clean-arch/src/utils"
)

func init() {
	config.MustLoadConfig("./config")
	config.MustLoadFirebaseConfig("./config")
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name X-Access-Token

func main() {
	cfg := config.Get()
	client := utils.GetFirestoreClient(config.GetFirebaseApp())

	//* Repos
	userRepo := user.NewFirestoreUserRepository(client, utils.UserColl)
	prodRepo := product.NewFirestoreProductRepository(client, utils.ProdColl)
	catRepo := category.NewFirestoreCategoryRepository(client, utils.CategColl)
	addrRepo := address.NewFirestoreAddressRepository(client, utils.AddrColl)
	ordRepo := order.NewFirestoreOrderRepository(client, utils.OrderColl)

	//* Services
	userSvc := core.NewUserService(userRepo)
	prodSvc := core.NewProductService(prodRepo, catRepo)
	catSvc := core.NewCategoryService(catRepo, prodRepo)
	addrSvc := core.NewAddressService(addrRepo, userRepo)
	ordSvc := core.NewOrderService(ordRepo, addrRepo)
	pmSvc := payment.NewStripePaymentService(cfg.Stripe.SecretKey, userRepo)
	emailSvc := email.NewSmtpEmailService()
	vldSvc := validation.NewValidationService()
	jwtSvc := auth.NewJWTService()

	//* Midlewares
	authMdlw := middlewares.NewAuthMiddleware(jwtSvc)
	paymMdlw := middlewares.NewPaymentMiddleware(pmSvc)

	//* Handlers
	usrHdlr := handlers.NewUserHandler(userSvc, vldSvc)
	authHdlr := handlers.NewAuthHandler(userSvc, emailSvc, jwtSvc, vldSvc)
	prodHdlr := handlers.NewProdutHandler(prodSvc, catSvc, vldSvc)
	catHdlr := handlers.NewCategoryHandler(prodSvc, catSvc, vldSvc)
	addrHdlr := handlers.NewAddressHandler(userSvc, addrSvc, vldSvc)
	cardHdlr := handlers.NewCardHandler(userSvc, pmSvc, vldSvc)
	ordHdlr := handlers.NewOrderHandler(userSvc, ordSvc, prodSvc, pmSvc, vldSvc)

	//* Server
	server := api.New()

	//* Setup
	server.SetGlobalMiddlewares()

	//* Routes
	server.CreateAuthRoutes(authHdlr, authMdlw)
	server.CreateUserRoutes(usrHdlr, authMdlw)
	server.CreateProductRoutes(prodHdlr, authMdlw)
	server.CreateCategoryRoutes(catHdlr, authMdlw)
	server.CreateAdreesesRoutes(addrHdlr, authMdlw)
	server.CreateCardRoutes(cardHdlr, paymMdlw, authMdlw)
	server.CreateOrderRoutes(ordHdlr, paymMdlw, authMdlw)

	go server.InitBackgroundTaks()

	if err := server.Start(":" + cfg.Api.Port); err != nil {
		log.Fatal(err)
	}
}
