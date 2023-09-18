package test

import (
	"net/http"
	"testing"

	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/stretchr/testify/suite"
	//"github.com/stripe/stripe-go/v74"
	//"github.com/stripe/stripe-go/v74/paymentmethod"
)

type CardRoutesSuite struct {
	ServerSuite
	bp string
}

func TestCardRoutesSuite(t *testing.T) {
	crds := new(CardRoutesSuite)
	crds.bp = "/api/card"
	suite.Run(t, crds)
}

func (s *CardRoutesSuite) TestCardRoutes_BGetAll() {
	path := s.bp + "/list"

	s.T().Logf("\n\n>>> %s\n\n", s.paymentID)

	utils.PrintBlueTesting(s.T(), "B GETTING")

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

func (s *CardRoutesSuite) TestCardRoutes_ASaveCard() {
	path := s.bp + "/save"

	utils.PrintBlueTesting(s.T(), "A SAVING")

	testCases := []TryRouteTestCase{
		{
			desc:          "Not authenticated",
			req:           s.MakeReq("POST", path, nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("POST", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid payment id",
			req: s.MakeReq("POST", path, dtos.SaveCardDTO{
				PaymentID: "adsfads",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Valid payment id",
			req: s.MakeReq("POST", path, dtos.SaveCardDTO{
				PaymentID: s.paymentID,
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

func (s *CardRoutesSuite) TestCardRoutes_CRemoveCard() {
	path := s.bp + "/remove"

	utils.PrintBlueTesting(s.T(), "C REMOVING")

	testCases := []TryRouteTestCase{
		{
			desc:          "Not authenticated",
			req:           s.MakeReq("DELETE", path+"/FDAF", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid payment id",
			req: s.MakeReq("DELETE", path+"/asd", dtos.SaveCardDTO{
				PaymentID: "adsfads",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Valid payment id",
			req: s.MakeReq("DELETE", path+"/"+s.paymentID, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}
