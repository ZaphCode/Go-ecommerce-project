package test

import (
	"net/http"
	"testing"

	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CategoryRoutesSuite struct {
	ServerSuite
	bp string
}

func TestCategoryRoutesSuite(t *testing.T) {
	prs := new(CategoryRoutesSuite)
	prs.bp = "/api/category"
	suite.Run(t, prs)
}

func (s *CategoryRoutesSuite) TestCategoryRoutes_GetAll() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Get all categories",
			req:           s.MakeReq("GET", s.bp+"/all", nil),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *CategoryRoutesSuite) TestCategoryRoutes_Delete() {
	path := s.bp + "/delete/"

	testCases := []TryRouteTestCase{
		{
			desc:          "Not token provided",
			req:           s.MakeReq("DELETE", path+"dafadf", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User has not permissions",
			req: s.MakeReq("DELETE", path+"dafadf", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid category id",
			req: s.MakeReq("DELETE", path+"dafadf", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusNotAcceptable,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Not found category",
			req: s.MakeReq("DELETE", path+uuid.New().String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Cannot delete category with products",
			req: s.MakeReq("DELETE", path+utils.CategoryExp3.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Delete category success",
			req: s.MakeReq("DELETE", path+utils.CategoryExp1.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}
