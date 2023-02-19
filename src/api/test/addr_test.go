package test

import (
	"net/http"
	"testing"

	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AddressRoutesSuite struct {
	ServerSuite
	bp string
}

func TestAddressRoutesSuite(t *testing.T) {
	prs := new(AddressRoutesSuite)
	prs.bp = "/api/address"
	suite.Run(t, prs)
}

func (s *AddressRoutesSuite) TestAddressRoutes_GetAll() {
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
			desc: "Auth user addresses",
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

func (s *AddressRoutesSuite) TestAddressRoutes_Create() {
	path := s.bp + "/create"

	testCases := []TryRouteTestCase{
		{
			desc:          "No token provided",
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
			desc: "Invalid request body",
			req: s.MakeReq("POST", path, dtos.NewAddressDTO{
				Name:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				PostalCode: "tacos",
				Line1:      "calle uwu",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Create success",
			req: s.MakeReq("POST", path, dtos.NewAddressDTO{
				Name:       "CSL main house",
				City:       "Cabo San Lucas",
				Country:    "Mexico",
				PostalCode: "23473",
				Line1:      "Calle ballena",
				Line2:      "Calle pargo",
				State:      "BCS",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusCreated,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *AddressRoutesSuite) TestAddressRoutes_Update() {
	path := s.bp + "/update/"

	testCases := []TryRouteTestCase{
		{
			desc:          "No token provided",
			req:           s.MakeReq("PUT", path+"k", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid uuid",
			req: s.MakeReq("PUT", path+"kfad", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusNotAcceptable,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid json",
			req: s.MakeReq("PUT", path+utils.AddrExp1.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Not found",
			req: s.MakeReq("PUT", path+uuid.New().String(), dtos.UpdateAddressDTO{}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("PUT", path+utils.AddrExp1.ID.String(), dtos.UpdateAddressDTO{
				Name:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				PostalCode: "tacos",
				Line1:      "calle uwu",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Not owner",
			req: s.MakeReq("PUT", path+utils.AddrExp1.ID.String(), dtos.NewAddressDTO{
				Country: "México",
				State:   "Baja California Sur",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Update success",
			req: s.MakeReq("PUT", path+utils.AddrExp1.ID.String(), dtos.NewAddressDTO{
				Country: "México",
				State:   "Baja California Sur",
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
