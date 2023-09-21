package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ZaphCode/clean-arch/config"
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
	"github.com/stretchr/testify/suite"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentmethod"
)

type TryRouteTestCase struct {
	desc          string
	req           *http.Request
	wantStatus    int
	showResp      bool
	bodyValidator func(jsm map[string]any)
}

type ServerSuite struct {
	suite.Suite
	server           *api.Server
	cfg              config.Config
	adminAccessToken string
	modAccessToken   string
	userAccessToken  string
	paymentID        string
}

func (s *ServerSuite) SetupSuite() {
	s.T().Logf("\n----------- SETUP ------------")

	config.MustLoadConfig("./../../../config")
	s.cfg = config.Get()

	// Repos
	userRepo := user.NewMemoryUserRepository(utils.UserAdmin, utils.UserExp1, utils.UserExp2)
	prodRepo := product.NewMemoryProductRepository(utils.ProductExp1, utils.ProductExp2)
	catRepo := category.NewMemoryCategoryRepository(utils.CategoryExp1, utils.CategoryExp2, utils.CategoryExp3)
	addrRepo := address.NewMemoryAddressRepository(utils.AddrExp1, utils.AddrExp2)
	ordRepo := order.NewMemoryOrderRepository()

	// Services
	userSvc := core.NewUserService(userRepo)
	prodSvc := core.NewProductService(prodRepo, catRepo)
	catSvc := core.NewCategoryService(catRepo, prodRepo)
	addrSvc := core.NewAddressService(addrRepo, userRepo)
	ordSvc := core.NewOrderService(ordRepo, addrRepo)
	pmSvc := payment.NewStripePaymentService(s.cfg.Stripe.SecretKey, userRepo)
	emailSvc := email.NewSmtpEmailService()
	vldSvc := validation.NewValidationService()
	jwtSvc := auth.NewJWTService()

	// Midlewares
	authMdlw := middlewares.NewAuthMiddleware(jwtSvc)
	paymMdlw := middlewares.NewPaymentMiddleware(pmSvc)

	// Handlers
	usrHdlr := handlers.NewUserHandler(userSvc, vldSvc)
	addrHdlr := handlers.NewAddressHandler(userSvc, addrSvc, vldSvc)
	authHdlr := handlers.NewAuthHandler(userSvc, emailSvc, jwtSvc, vldSvc)
	prodHdlr := handlers.NewProductHandler(prodSvc, catSvc, vldSvc)
	catHdlr := handlers.NewCategoryHandler(prodSvc, catSvc, vldSvc)
	crdHdlr := handlers.NewCardHandler(userSvc, pmSvc, vldSvc)
	ordHdlr := handlers.NewOrderHandler(userSvc, ordSvc, prodSvc, pmSvc, vldSvc)

	// Server
	server := api.New()

	// Setup
	server.SetGlobalMiddlewares()

	// Routes
	server.CreateAuthRoutes(authHdlr, authMdlw)
	server.CreateUserRoutes(usrHdlr, authMdlw)
	server.CreateProductRoutes(prodHdlr, authMdlw)
	server.CreateCategoryRoutes(catHdlr, authMdlw)
	server.CreateAddressesRoutes(addrHdlr, authMdlw)
	server.CreateOrderRoutes(ordHdlr, paymMdlw, authMdlw)
	server.CreateCardRoutes(crdHdlr, paymMdlw, authMdlw)

	s.server = server

	// Admin token
	at, err := jwtSvc.CreateToken(auth.Claims{
		ID:   utils.UserAdmin.ID,
		Role: utils.UserAdmin.Role,
	}, time.Minute*1, s.cfg.Api.AccessTokenSecret)

	s.NoError(err, "should not be error")

	s.adminAccessToken = at

	// Mod token
	mt, err := jwtSvc.CreateToken(auth.Claims{
		ID:   utils.UserExp2.ID,
		Role: utils.UserExp2.Role,
	}, time.Minute*1, s.cfg.Api.AccessTokenSecret)

	s.NoError(err, "should not be error")

	s.modAccessToken = mt

	// User token
	ut, err := jwtSvc.CreateToken(auth.Claims{
		ID:   utils.UserExp1.ID,
		Role: utils.UserExp1.Role,
	}, time.Minute*1, s.cfg.Api.AccessTokenSecret)

	s.NoError(err, "should not be error")

	s.userAccessToken = ut

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String("5555555555554444"),
			ExpMonth: stripe.Int64(4),
			ExpYear:  stripe.Int64(2025),
			CVC:      stripe.String("314"),
		},
		BillingDetails: &stripe.PaymentMethodBillingDetailsParams{
			Name: stripe.String("FROM SERVER SUITE"),
		},
		Type: stripe.String("card"),
	}

	pm, err := paymentmethod.New(params)

	s.NoError(err, "should not be error")

	s.paymentID = pm.ID
}

func (s *ServerSuite) RunRequests(testCases []TryRouteTestCase) {
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			res, err := s.server.TryRoute(tC.req)

			s.NoError(err, "request error!")

			s.Equal(tC.wantStatus, res.StatusCode, "wrong status code!")

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			s.Require().NoError(err, "reading response error!")

			if tC.showResp {
				s.ShowResponse(res.Request.URL.Path, resBody)
			}

			if tC.bodyValidator != nil {

				s.Equal("application/json", res.Header.Get("Content-Type"), "response should be json!")

				jsm := make(map[string]any)

				s.NoError(json.Unmarshal(resBody, &jsm), "unmarshall err")

				tC.bodyValidator(jsm)
			}
		})
	}
}

func (s *ServerSuite) ShowResponse(path string, resp []byte) {
	result := string(resp)
	buff := new(bytes.Buffer)

	if err := json.Indent(buff, resp, "", "    "); err == nil {
		result = "(JSON) " + buff.String()
	}

	s.T().Logf("\n\n >>> %s : %s \n\n", path, result)
}

func (s *ServerSuite) MakeReq(met, path string, body any, hdrs ...map[string]string) *http.Request {
	var reqBody io.Reader = nil

	if body != nil {
		marsh, err := json.Marshal(body)

		if err != nil {
			s.FailNow("json marshal error")
		}

		reqBody = bytes.NewBuffer(marsh)
	}
	req, err := http.NewRequest(met, path, reqBody)

	if err != nil {
		s.FailNow("json marshal error")
	}

	if len(hdrs) > 0 {
		for h, v := range hdrs[0] {
			req.Header.Set(h, v)
		}
	}

	return req
}

func (s *ServerSuite) CheckSuccess(jsm map[string]any) {
	status, ok := jsm["status"]
	s.True(ok, "should contain status")
	s.Equal("success", status, "status should be success")
}

func (s *ServerSuite) CheckFail(jsm map[string]any) {
	status, ok := jsm["status"]
	s.True(ok, "should contain failure")
	s.Equal("failure", status, "status should be failure")
}
