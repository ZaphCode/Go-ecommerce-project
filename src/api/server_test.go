package api

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ZaphCode/clean-arch/config"
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
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	server *Server
	cfg    *config.Config
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func (s *ServerSuite) SetupSuite() {
	s.T().Logf("\n----------- SETUP ------------")

	config.MustLoadConfig("./../../config")
	config.MustLoadFirebaseConfig("./../../config")

	s.cfg = config.Get()

	client := utils.GetFirestoreClient(config.GetFirebaseApp())

	// Repos
	userRepo := user.NewFirestoreUserRepository(client, utils.UserColl)
	prodRepo := product.NewFirestoreProductRepository(client, utils.ProdColl)
	catRepo := category.NewFirestoreCategoryRepository(client, utils.CategColl)

	// Services
	userSvc := core.NewUserService(userRepo)
	prodSvc := core.NewProductService(prodRepo)
	catSvc := core.NewCategoryService(catRepo)
	emailSvc := email.NewSmtpEmailService()
	vldSvc := validation.NewValidationService()
	jwtSvc := auth.NewJWTService()

	// Midlewares
	authMdlw := middlewares.NewAuthMiddleware(jwtSvc)

	// Handlers
	usrHdlr := handlers.NewUserHandler(userSvc, vldSvc)
	authHdlr := handlers.NewAuthHandler(userSvc, emailSvc, jwtSvc, vldSvc)
	prodHdlr := handlers.NewProdutHandler(prodSvc, catSvc, vldSvc)
	catHdlr := handlers.NewCategoryHandler(prodSvc, catSvc, vldSvc)

	// Server
	server := New()

	// Setup
	server.SetGlobalMiddlewares()

	// Routes
	server.CreateAuthRoutes(authHdlr, authMdlw)
	server.CreateUserRoutes(usrHdlr, authMdlw)
	server.CreateProductRoutes(prodHdlr, authMdlw)
	server.CreateCategoryRoutes(catHdlr, authMdlw)

	s.server = server
}

type jsonMap map[string]interface{}

type tryRouteTestCase struct {
	desc          string
	reqMaker      func() *http.Request
	bodyValidator func(jsm jsonMap)
	checkReqBody  bool
	wantStatus    int
}

func (s *ServerSuite) TestTableRandomRoutes() {
	testCases := []tryRouteTestCase{
		{
			desc: "Not found",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/sexogratis", nil)
				s.NoError(err, "req error")
				return req
			},
			//wantErr:      false,
			checkReqBody: true,
			wantStatus:   http.StatusNotFound,
		},
		{
			desc: "Healty check",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
				s.NoError(err, "req error")
				return req
			},
			checkReqBody: true,
			//wantErr:      false,
			wantStatus: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			res, err := s.server.TryRoute(tC.reqMaker())

			s.Require().NoError(err, "response error")

			s.Require().Equal(res.StatusCode, tC.wantStatus, "wrong status code")

			resBody, err := io.ReadAll(res.Body)

			s.Require().NoError(err, "response error")

			defer res.Body.Close()

			s.T().Logf("\n\n XXX Response: %s \n\n", string(resBody))
		})
	}
}

func (s *ServerSuite) TestTableProductRoutes() {
	r := s.Require()

	testCases := []tryRouteTestCase{
		{
			desc: "Get all products",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/product/all", nil)
				r.NoError(err, "req error")
				return req
			},
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]

				r.True(ok, "should contain status")

				r.Equal(status, "success", "status sould be success")
			},
			checkReqBody: true,
			wantStatus:   http.StatusOK,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			res, err := s.server.TryRoute(tC.reqMaker())

			r.NoError(err, "response error")

			r.Equal(res.StatusCode, tC.wantStatus, "wrong status code")

			resBody, err := io.ReadAll(res.Body)

			r.NoError(err, "response error")

			defer res.Body.Close()

			s.T().Logf("\n\n XXX Response: %s \n\n", string(resBody))

			if tC.checkReqBody {
				jsm := make(jsonMap)

				r.NoError(json.Unmarshal(resBody, &jsm), "Unmarshall err")

				tC.bodyValidator(jsm)
			}
		})
	}
}
