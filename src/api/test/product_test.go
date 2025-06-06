package test

import (
	"net/http"
	"testing"

	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ProductRoutesSuite struct {
	ServerSuite
	bp string
}

func TestProductRoutesSuite(t *testing.T) {
	prs := new(ProductRoutesSuite)
	prs.bp = "/api/product"
	suite.Run(t, prs)
}

func (s *ProductRoutesSuite) TestProductRoutes_GetAll() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Get all products",
			req:           s.MakeReq("GET", s.bp+"/all", nil),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *ProductRoutesSuite) TestProductRoutes_GetByID() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Invalid product id",
			req:           s.MakeReq("GET", s.bp+"/get/dafadf", nil),
			showResp:      true,
			wantStatus:    http.StatusNotAcceptable,
			bodyValidator: s.CheckFail,
		},
		{
			desc:          "Not found product",
			req:           s.MakeReq("GET", s.bp+"/get/"+uuid.New().String(), nil),
			showResp:      true,
			wantStatus:    http.StatusNotFound,
			bodyValidator: s.CheckFail,
		},
		{
			desc:          "Get product success",
			req:           s.MakeReq("GET", s.bp+"/get/"+utils.ProductExpToDev1.ID.String(), nil),
			showResp:      true,
			wantStatus:    http.StatusFound,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *ProductRoutesSuite) TestProductRoutes_Delete() {
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
			desc: "Invalid product id",
			req: s.MakeReq("DELETE", s.bp+"/delete/dafadf", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusNotAcceptable,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Not found product",
			req: s.MakeReq("DELETE", s.bp+"/delete/"+uuid.New().String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Delete product success",
			req: s.MakeReq("DELETE", s.bp+"/delete/"+utils.ProductExp1.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *ProductRoutesSuite) TestProductRoutes_Create() {
	testCases := []TryRouteTestCase{
		{
			desc:          "No token provided",
			req:           s.MakeReq("POST", s.bp+"/create", nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User has not permissions",
			req: s.MakeReq("POST", s.bp+"/create", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.userAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Mod has not permissions",
			req: s.MakeReq("POST", s.bp+"/create", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.modAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusForbidden,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Unprocesable json",
			req: s.MakeReq("POST", s.bp+"/create", nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("POST", s.bp+"/create", dtos.NewProductDTO{
				Price:        0,
				DiscountRate: 120,
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Unexisting category",
			req: s.MakeReq("POST", s.bp+"/create", dtos.NewProductDTO{
				Category:     "toys",
				Name:         "Woody toy story",
				Description:  "Best toy ever",
				ImagesUrl:    []string{"https://toy.com/woody"},
				Tags:         []string{"toys"},
				Price:        14,
				DiscountRate: 10,
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
			req: s.MakeReq("POST", s.bp+"/create", dtos.NewProductDTO{
				Category:     "headsets",
				Name:         "Logitech g613",
				Description:  "Best cheap headsets",
				ImagesUrl:    []string{"https://logi.com/headsets"},
				Tags:         []string{"headsets", "tech"},
				Price:        1400,
				DiscountRate: 0,
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

func (s *ProductRoutesSuite) TestProductRoutes_Update() {
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
			req: s.MakeReq("PUT", path+utils.ProductExp2.ID.String(), nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("PUT", path+utils.ProductExp2.ID.String(), dtos.UpdateProductDTO{
				Category:     " a ",
				DiscountRate: utils.PTR[int64](200),
				ImagesUrl:    []string{"A"},
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Product not found",
			req: s.MakeReq("PUT", path+utils.AddrExp1.ID.String(), dtos.UpdateProductDTO{
				DiscountRate: utils.PTR[int64](10),
			}, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
				"Content-Type":              "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Category to update not found",
			req: s.MakeReq("PUT", path+utils.ProductExp2.ID.String(), dtos.UpdateProductDTO{
				Category: "toys",
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
			req: s.MakeReq("PUT", path+utils.ProductExpToDev1.ID.String(), dtos.UpdateProductDTO{
				Category:     "clothes",
				Price:        utils.PTR[int64](1500),
				DiscountRate: utils.PTR[int64](12),
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
