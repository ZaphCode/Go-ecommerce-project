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
	r.Post("/create", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), prodHdlr.CreateProducts)
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
