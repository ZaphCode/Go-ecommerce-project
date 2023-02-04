package api

import (
	"bytes"
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	server *Server
	cfg    *config.Config
}

//* Main

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

//* Life cycle

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

//* Tests

type jsonMap map[string]interface{}

type tryRouteTestCase struct {
	desc          string
	reqMaker      func() *http.Request
	bodyValidator func(jsm jsonMap)
	wantStatus    int
	showResp      bool
}

func (s *ServerSuite) TestRandomRoutes() {
	testCases := []tryRouteTestCase{
		{
			desc: "Not found",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/sexogratis", nil)
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusNotFound,
		},
		{
			desc: "Healty check",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusOK,
		},
	}

	s.runRequests(testCases)
}

func (s *ServerSuite) TestProductRoutes() {
	testCases := []tryRouteTestCase{
		{
			desc: "get all products",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/product/all", nil)
				s.Require().NoError(err, "req error")
				return req
			},
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				s.Equal(status, "success", "status sould be success")
			},
			wantStatus: http.StatusOK,
		},
		{
			desc: "get product: invalid id",
			reqMaker: func() *http.Request {
				randStr := utils.RandomString(20)
				req, err := http.NewRequest(http.MethodGet, "/api/product/get/"+randStr, nil)
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusNotAcceptable,
		},
		{
			desc: "get product: valid id but not found",
			reqMaker: func() *http.Request {
				idStr := uuid.New().String()
				req, err := http.NewRequest(http.MethodGet, "/api/product/get/"+idStr, nil)
				s.Require().NoError(err, "req error")
				return req
			},
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				_, ok = jsm["message"]
				s.True(ok, "should contain message")
				s.Equal(status, "failure", "status sould be success")
			},
			wantStatus: http.StatusNotFound,
		},
	}

	s.runRequests(testCases)
}

func (s *ServerSuite) runRequests(testCases []tryRouteTestCase) {
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			res, err := s.server.TryRoute(tC.reqMaker())

			s.NoError(err, "request error!")

			s.Equal(tC.wantStatus, res.StatusCode, "wrong status code!")

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			s.Require().NoError(err, "readig response error!")

			if tC.showResp {
				s.showResponse(res.Request.URL.Path, resBody)
			}

			if tC.bodyValidator != nil {

				s.Equal("application/json", res.Header.Get("Content-Type"), "response should be json!")

				jsm := make(jsonMap)

				s.NoError(json.Unmarshal(resBody, &jsm), "unmarshall err")

				tC.bodyValidator(jsm)

			}
		})
	}
}

func (s *ServerSuite) showResponse(path string, resp []byte) {
	result := string(resp)
	buff := new(bytes.Buffer)

	if err := json.Indent(buff, resp, "", "    "); err == nil {
		result = "(JSON) " + buff.String()
	}

	s.T().Logf("\n\n >>> %s : %s \n\n", path, result)
}
