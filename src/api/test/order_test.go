package test

import (
	"net/http"
	"testing"

	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/stretchr/testify/suite"
)

type OrderRoutesSuite struct {
	ServerSuite
	bp string
}

func TestOrderRoutesSuite(t *testing.T) {
	ords := new(OrderRoutesSuite)
	ords.bp = "/api/order"
	suite.Run(t, ords)
}

func (s *OrderRoutesSuite) TestOrderRoutes_GetAll() {
	path := s.bp + "/list"

	testCases := []TryRouteTestCase{
		{
			desc:          "Not authenticated",
			req:           s.MakeReq("GET", path, nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Auth user cards",
			req: s.MakeReq("GET", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *OrderRoutesSuite) TestOrderRoutes_NewOrder() {
	path := s.bp + "/new"

	testCases := []TryRouteTestCase{
		{
			desc:          "Not authenticated",
			req:           s.MakeReq("POST", path, nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Unprocesable json",
			req: s.MakeReq("POST", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid body (empty)",
			req: s.MakeReq("POST", path, dtos.NewOrderDTO{}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid body (unexisting products)",
			req: s.MakeReq("POST", path, dtos.NewOrderDTO{
				PaymentID: s.paymentID,
				AddressID: utils.AddrExp1.ID,
				Products: []domain.OrderProduct{
					{ID: utils.AddrExp1.ID, Quantity: 13},
					{ID: utils.AddrExp2.ID, Quantity: 3},
				},
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid address id",
			req: s.MakeReq("POST", path, dtos.NewOrderDTO{
				PaymentID: s.paymentID,
				AddressID: utils.UserAdmin.ID,
				Products: []domain.OrderProduct{
					{ID: utils.ProductExp1.ID, Quantity: 13},
					{ID: utils.ProductExp1.ID, Quantity: 3},
				},
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid payment method",
			req: s.MakeReq("POST", path, dtos.NewOrderDTO{
				PaymentID: "adsfadf",
				AddressID: utils.AddrExp1.ID,
				Products: []domain.OrderProduct{
					{ID: utils.ProductExp1.ID, Quantity: 13},
					{ID: utils.ProductExp2.ID, Quantity: 3},
				},
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Attached payment method",
			req: s.MakeReq("POST", path, dtos.NewOrderDTO{
				PaymentID: "pm_1NKP27G8UXDxPRbaNZRE6Ajd",
				AddressID: utils.AddrExp1.ID,
				Products: []domain.OrderProduct{
					{ID: utils.ProductExp1.ID, Quantity: 13},
					{ID: utils.ProductExp2.ID, Quantity: 3},
				},
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Proper work with new payment method",
			req: s.MakeReq("POST", path, dtos.NewOrderDTO{
				PaymentID: s.paymentID,
				AddressID: utils.AddrExp1.ID,
				Products: []domain.OrderProduct{
					{ID: utils.ProductExp1.ID, Quantity: 2},
					{ID: utils.ProductExp2.ID, Quantity: 3},
				},
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
		{
			desc: "Proper work with saved card",
			req: s.MakeReq("POST", path, dtos.NewOrderDTO{
				PaymentID: "pm_1NKP27G8UXDxPRbaNZRE6Ajd",
				AddressID: utils.AddrExp1.ID,
				Products: []domain.OrderProduct{
					{ID: utils.ProductExp1.ID, Quantity: 2},
					{ID: utils.ProductExp2.ID, Quantity: 3},
				},
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}
