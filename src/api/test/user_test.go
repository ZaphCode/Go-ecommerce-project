package test

import (
	"net/http"
	"testing"

	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserRoutesSuite struct {
	ServerSuite
	bp string
}

func TestUserRoutesSuite(t *testing.T) {
	prs := new(UserRoutesSuite)
	prs.bp = "/api/user"
	suite.Run(t, prs)
}

func (s *UserRoutesSuite) TestUserRoutes_GetAll() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Get all users",
			req:           s.MakeReq("GET", s.bp+"/all", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User not authorized",
			req: s.MakeReq("GET", s.bp+"/all", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Mod/Admin authorized",
			req: s.MakeReq("GET", s.bp+"/all", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *UserRoutesSuite) TestUserRoutes_GetByID() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Get all users",
			req:           s.MakeReq("GET", s.bp+"/all", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User not authorized",
			req: s.MakeReq("GET", s.bp+"/all", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid user id",
			req: s.MakeReq("GET", s.bp+"/get/dafadf", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusNotAcceptable,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User not found",
			req: s.MakeReq("GET", s.bp+"/get/"+uuid.New().String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusNotFound,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Get user success",
			req: s.MakeReq("GET", s.bp+"/get/"+utils.UserExp1.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusFound,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *UserRoutesSuite) TestUserRoutes_Delete() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Not token provided",
			req:           s.MakeReq("DELETE", s.bp+"/delete/dafadf", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User has not permissions",
			req: s.MakeReq("DELETE", s.bp+"/delete/dafadf", nil, map[string]string{
				"" + s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Mod has not permissions",
			req: s.MakeReq("DELETE", s.bp+"/delete/dafadf", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid user id",
			req: s.MakeReq("DELETE", s.bp+"/delete/dafadf", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusNotAcceptable,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User not found",
			req: s.MakeReq("DELETE", s.bp+"/delete/"+uuid.New().String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Delete user success",
			req: s.MakeReq("DELETE", s.bp+"/delete/"+utils.UserExp2.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *UserRoutesSuite) TestUserRoutes_Create() {
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
			desc: "User has not permissions",
			req: s.MakeReq("POST", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Mod has not permissions",
			req: s.MakeReq("POST", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Unprocesable json",
			req: s.MakeReq("POST", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("POST", path, dtos.NewUserDTO{
				Email:         "omar#gmail.com",
				VerifiedEmail: true,
				Role:          "seller",
				Password:      "abc",
				ImageUrl:      "a",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Existing email",
			req: s.MakeReq("POST", path, dtos.NewUserDTO{
				Username:      "juan camanei",
				Email:         "john@testing.com",
				VerifiedEmail: true,
				Password:      "password12345",
				Role:          "user",
				Age:           23,
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Create success",
			req: s.MakeReq("POST", path, dtos.NewUserDTO{
				Username:      "Gavino",
				Email:         "gabi@hotmail.com.mx",
				VerifiedEmail: false,
				Password:      "contrasenauwu",
				Role:          "moderator",
				Age:           25,
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusCreated,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *UserRoutesSuite) TestUserRoutes_Update() {
	path := s.bp + "/update/"

	testCases := []TryRouteTestCase{
		{
			desc:          "No token provided",
			req:           s.MakeReq("PUT", path+"a", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User has not permissions",
			req: s.MakeReq("PUT", path+"z", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Mod has not permissions",
			req: s.MakeReq("PUT", path+"z", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid id",
			req: s.MakeReq("PUT", path+"wakawakweheh", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusNotAcceptable,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Unprocesable json",
			req: s.MakeReq("PUT", path+utils.UserExp2.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("PUT", path+utils.UserExp1.ID.String(), dtos.UpdateUserDTO{
				Username: "Fizz bar",
				Age:      utils.PTR[uint16](12),
				Role:     "helper",
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User not found",
			req: s.MakeReq("PUT", path+utils.AddrExp1.ID.String(), dtos.UpdateProductDTO{
				DiscountRate: utils.PTR[uint](10),
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Update success",
			req: s.MakeReq("PUT", path+utils.UserExp1.ID.String(), dtos.UpdateUserDTO{
				Role:          "moderator",
				Username:      "Fizz bar",
				VerifiedEmail: utils.PTR(false),
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}
