package api

import (
	"github.com/ZaphCode/clean-arch/src/api/handlers"
	"github.com/ZaphCode/clean-arch/src/api/middlewares"
	"github.com/ZaphCode/clean-arch/src/utils"
)

func (s *Server) CreateAuthRoutes(
	authHdlr *handlers.AuthHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/auth")
	r.Get("/:provider/url", authHdlr.GetOAuthUrl)
	r.Get("/:provider/callback", authHdlr.SignInWihOAuth)
	r.Get("/refresh", authHdlr.RefreshToken)
	r.Get("/me", authMdlw.AuthRequired, authHdlr.GetAuthUser)
	r.Get("/signout", authHdlr.SignOut)
	r.Post("/signin", authHdlr.SignIn)
	r.Post("/signup", authHdlr.SignUp)
}

func (s *Server) CreateUserRoutes(
	usrHdlr *handlers.UserHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/user")
	r.Get("/all", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), usrHdlr.GetUsers)
	r.Get("/get/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), usrHdlr.GetUser)
	r.Post("/create", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), usrHdlr.CreateUser)
	r.Put("/update/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), usrHdlr.UpdateUser)
	r.Delete("/delete/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), usrHdlr.DeleteUser)
}

func (s *Server) CreateProductRoutes(
	prodHdlr *handlers.ProductHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/product")
	r.Get("/all", prodHdlr.GetProducts)
	r.Get("/get/:id", prodHdlr.GetProduct)
	r.Post("/create", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), prodHdlr.CreateProducts)
	r.Put("/update/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), prodHdlr.UpdateProduct)
	r.Delete("/delete/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), prodHdlr.DeleteProduct)
}

func (s *Server) CreateCategoryRoutes(
	catHdlr *handlers.CategoryHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/category")
	r.Get("/all", catHdlr.GetCategories)
	r.Post("/create", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), catHdlr.CreateCategory)
	r.Delete("/delete/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), catHdlr.DeleteCategory)
}

func (s *Server) CreateAddressesRoutes(
	addrHdlr *handlers.AddressHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/address")
	r.Get("/list", authMdlw.AuthRequired, addrHdlr.GetUserAddress)
	r.Post("/create", authMdlw.AuthRequired, addrHdlr.CreateAddress)
	r.Put("/update/:id", authMdlw.AuthRequired, addrHdlr.UpdateAddress)
	r.Delete("/delete/:id", authMdlw.AuthRequired, addrHdlr.DeleteAddress)
}

func (s *Server) CreateCardRoutes(
	crdHdlr *handlers.CardHandler,
	paymMdlw *middlewares.PaymentMiddleware,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/card")
	r.Get("/list", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, crdHdlr.GetUserCards)
	r.Post("/save", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, crdHdlr.SaveUserCard)
	r.Delete("/remove/:id", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, crdHdlr.RemoveUserCard)
}

func (s *Server) CreateOrderRoutes(
	ordHdlr *handlers.OrderHandler,
	paymMdlw *middlewares.PaymentMiddleware,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/order")
	r.Get("/list", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, ordHdlr.GetOrders)
	r.Post("/new", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, ordHdlr.CreateOrder)
}
