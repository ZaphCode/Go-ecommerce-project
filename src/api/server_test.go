package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/handlers"
	"github.com/ZaphCode/clean-arch/src/api/middlewares"
	"github.com/ZaphCode/clean-arch/src/domain"
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
	server      *Server
	cfg         *config.Config
	accessToken string
}

//* Main

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

//* Life cycle

var rootUser = domain.User{
	Model: domain.Model{
		ID:        uuid.MustParse("e44ef83a-a1c7-11ed-a865-7e82d40d4740"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	CustomerID:    "",
	Username:      "Zaphkiel",
	Email:         "zaph@fapi.com",
	Password:      "$2a$10$D/cBZACDWfS4r910QhyIhucC/IKKD.4ilevC44j2CozjW0fscNBaG",
	Role:          utils.AdminRole,
	VerifiedEmail: true,
	ImageUrl:      "https://i.etsystatic.com/15149849/r/il/16852c/2126930485/il_570xN.2126930485_oub4.jpg",
	Age:           19,
}

func (s *ServerSuite) SetupSuite() {
	s.T().Logf("\n----------- SETUP ------------")

	config.MustLoadConfig("./../../config")
	s.cfg = config.Get()

	// Repos
	userRepo := user.NewMemoryUserRepository(rootUser)
	prodRepo := product.NewMemoryProductRepository()
	catRepo := category.NewMemoryProductRepository()

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

	token, err := jwtSvc.CreateToken(auth.Claims{
		ID:   rootUser.ID,
		Role: rootUser.Role,
	}, time.Minute*1, s.cfg.Api.AccessTokenSecret)

	s.NoError(err, "should not be error")

	s.accessToken = token
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
				req, err := http.NewRequest(http.MethodGet, "/api/random", nil)
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

//* Product routes

func (s *ServerSuite) TestProductRoutes_GetAll() {
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
				s.Equal("success", status, "status sould be success")
			},
			showResp:   true,
			wantStatus: http.StatusOK,
		},
	}
	s.runRequests(testCases)
}

func (s *ServerSuite) TestProductRoutes_GetByID() {
	testCases := []tryRouteTestCase{
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
				s.Equal("failure", status, "status sould be success")
			},
			wantStatus: http.StatusNotFound,
		},
	}

	s.runRequests(testCases)
}

//* Auth routes

func (s *ServerSuite) TestAuthRoutes_SignUp() {
	testCases := []tryRouteTestCase{
		{
			desc: "sign up properly",
			reqMaker: func() *http.Request {
				reqBody, err := json.Marshal(dtos.SignupDTO{
					Username: "zaph_mini",
					Email:    "zaph@mini.com",
					Password: "password1234",
					Age:      24,
				})
				s.Require().NoError(err, "req error")
				req, err := http.NewRequest(http.MethodPost, "/api/auth/signup/", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusCreated,
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				s.Equal("success", status, "status sould be success")
			},
		},
		{
			desc: "sign up with invalid json body",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/api/auth/signup/", nil)
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusUnprocessableEntity,
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				s.Equal("failure", status, "status sould be failure")
			},
		},
		{
			desc: "sign up with invalid data",
			reqMaker: func() *http.Request {
				reqBody, err := json.Marshal(dtos.SignupDTO{
					Username: "zap",
					Email:    "zaph@minicom",
					Password: "1234",
					Age:      14,
				})
				s.Require().NoError(err, "req error")
				req, err := http.NewRequest(http.MethodPost, "/api/auth/signup/", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusBadRequest,
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				s.Equal("failure", status, "status sould be failure")
				_, ok = jsm["errors"]
				s.True(ok, "should contain errors")
				// _, ok = errs.([]interface{})
				// s.True(ok, "should be array")
			},
		},
	}

	s.runRequests(testCases)
}

func (s *ServerSuite) TestAuthRoutes_SignIn() {
	testCases := []tryRouteTestCase{
		{
			desc: "sign up properly",
			reqMaker: func() *http.Request {
				reqBody, err := json.Marshal(dtos.SigninDTO{
					Email:    "zaph@fapi.com",
					Password: "menosfapi33",
				})
				s.Require().NoError(err, "req error")
				req, err := http.NewRequest(http.MethodPost, "/api/auth/signin/", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusOK,
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				s.Equal("success", status, "status sould be success")
				data, ok := jsm["data"]
				s.True(ok, "should contain data")
				dMap, ok := data.(map[string]interface{})
				s.True(ok, "data should be a map")
				_, ok = dMap["access_token"]
				s.True(ok, "should contain access token")
			},
		},
		{
			desc: "sign up with invalid json body",
			reqMaker: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/api/auth/signin/", nil)
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusUnprocessableEntity,
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				s.Equal("failure", status, "status sould be failure")
			},
		},
		{
			desc: "sign up with invalid data",
			reqMaker: func() *http.Request {
				reqBody, err := json.Marshal(dtos.SigninDTO{
					Email: "zaph@testcom",
				})
				s.Require().NoError(err, "req error")
				req, err := http.NewRequest(http.MethodPost, "/api/auth/signin/", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				s.Require().NoError(err, "req error")
				return req
			},
			showResp:   true,
			wantStatus: http.StatusBadRequest,
			bodyValidator: func(jsm jsonMap) {
				status, ok := jsm["status"]
				s.True(ok, "should contain status")
				s.Equal("failure", status, "status sould be failure")
				_, ok = jsm["errors"]
				s.True(ok, "should contain errors")
				// _, ok = errs.([]interface{})
				// s.True(ok, "should be array")
			},
		},
	}

	s.runRequests(testCases)
}

//* Helpers

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
