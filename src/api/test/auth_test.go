package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/stretchr/testify/suite"
)

type AuthRoutesSuite struct {
	ServerSuite
	bp                  string
	adminRefreshToken   string
	expiredRefreshToken string
	expiredAccessToken  string
}

func TestAuthRoutesSuite(t *testing.T) {
	config.MustLoadConfig("./../../../config")

	rs := new(AuthRoutesSuite)
	jwtSvc := auth.NewJWTService()

	rt, err1 := jwtSvc.CreateToken(auth.Claims{
		ID:   utils.UserAdmin.ID,
		Role: utils.UserAdmin.Role,
	}, time.Minute*1, config.Get().Api.RefreshTokenSecret)

	expRt, err2 := jwtSvc.CreateToken(auth.Claims{
		ID:   utils.UserAdmin.ID,
		Role: utils.UserAdmin.Role,
	}, time.Nanosecond*0, config.Get().Api.RefreshTokenSecret)

	expAt, err3 := jwtSvc.CreateToken(auth.Claims{
		ID:   utils.UserAdmin.ID,
		Role: utils.UserAdmin.Role,
	}, time.Nanosecond*0, config.Get().Api.AccessTokenSecret)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("error creating a testing token")
	}

	rs.bp = "/api/auth"
	rs.adminRefreshToken = rt
	rs.expiredRefreshToken = expRt
	rs.expiredAccessToken = expAt

	suite.Run(t, rs)
}

func (s *AuthRoutesSuite) TestAuthRoutesSuite_SignIn() {
	path := s.bp + "/signin"

	testCases := []TryRouteTestCase{
		{
			desc:          "Unprocesable json",
			req:           s.MakeReq("POST", path, nil),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("POST", path, dtos.SigninDTO{
				Email:    "test@gmailcom",
				Password: "",
			}, map[string]string{
				"Content-Type": "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User not found",
			req: s.MakeReq("POST", path, dtos.SigninDTO{
				Email:    "test@gmail.com",
				Password: "password",
			}, map[string]string{
				"Content-Type": "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusNotFound,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Wrong password",
			req: s.MakeReq("POST", path, dtos.SigninDTO{
				Email:    "zaph@fapi.com",
				Password: "menosfapi3",
			}, map[string]string{
				"Content-Type": "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Signin success!",
			req: s.MakeReq("POST", path, dtos.SigninDTO{
				Email:    "zaph@fapi.com",
				Password: "menosfapi33",
			}, map[string]string{
				"Content-Type": "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *AuthRoutesSuite) TestAuthRoutesSuite_SignUp() {
	path := s.bp + "/signup"

	testCases := []TryRouteTestCase{
		{
			desc:          "Unprocesable json",
			req:           s.MakeReq("POST", path, nil),
			showResp:      true,
			wantStatus:    http.StatusUnprocessableEntity,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid request body",
			req: s.MakeReq("POST", path, dtos.SignupDTO{
				Email:    "test@gmailcom",
				Password: "13",
				Age:      1,
			}, map[string]string{
				"Content-Type": "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "User that already exists",
			req: s.MakeReq("POST", path, dtos.SignupDTO{
				Email:    "zaph@fapi.com",
				Password: "password",
				Username: "Zaphkiel 33",
				Age:      24,
			}, map[string]string{
				"Content-Type": "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusInternalServerError,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Signup success!",
			req: s.MakeReq("POST", path, dtos.SignupDTO{
				Email:    "juan@testing.com",
				Password: "juan1234",
				Username: "Juan Gonzales",
				Age:      18,
			}, map[string]string{
				"Content-Type": "application/json",
			}),
			showResp:      true,
			wantStatus:    http.StatusCreated,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *AuthRoutesSuite) TestAuthRoutesSuite_Signout() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Success signout",
			req:           s.MakeReq("GET", s.bp+"/signout", nil),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *AuthRoutesSuite) TestAuthRoutesSuite_AuthMe() {
	path := s.bp + "/me"

	testCases := []TryRouteTestCase{
		{
			desc:          "No token provided",
			req:           s.MakeReq("GET", path, nil),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid access token",
			req: s.MakeReq("GET", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: "dsfal;l;fkda;lfka;lf33v/fadsfa",
			}),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Expired access token",
			req: s.MakeReq("GET", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.expiredAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid signature",
			req: s.MakeReq("GET", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminRefreshToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusUnauthorized,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Success auth current user",
			req: s.MakeReq("GET", path, nil, map[string]string{
				s.cfg.Api.AccessTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *AuthRoutesSuite) TestAuthRoutesSuite_Refresh() {
	path := s.bp + "/refresh"

	testCases := []TryRouteTestCase{
		{
			desc:          "Refresh token not recibed netheir cookie and header",
			req:           s.MakeReq("GET", path, nil),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid refresh token",
			req: s.MakeReq("GET", path+"?method=header", nil, map[string]string{
				s.cfg.Api.RefreshTokenHeader: "dsfal;l;fkda;lfka;lf33v/fadsfa",
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Expired access token",
			req: s.MakeReq("GET", path+"?method=header", nil, map[string]string{
				s.cfg.Api.RefreshTokenHeader: s.expiredRefreshToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Invalid signature",
			req: s.MakeReq("GET", path+"?method=header", nil, map[string]string{
				s.cfg.Api.RefreshTokenHeader: s.adminAccessToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc: "Success refresh",
			req: s.MakeReq("GET", path+"?method=header", nil, map[string]string{
				s.cfg.Api.RefreshTokenHeader: s.adminRefreshToken,
			}),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}

func (s *AuthRoutesSuite) TestAuthRoutesSuite_GetOAuthUrl() {
	testCases := []TryRouteTestCase{
		{
			desc:          "Invalid oauth provider",
			req:           s.MakeReq("GET", s.bp+"/facebook/url", nil),
			showResp:      true,
			wantStatus:    http.StatusBadRequest,
			bodyValidator: s.CheckFail,
		},
		{
			desc:          "Success github url",
			req:           s.MakeReq("GET", s.bp+"/"+utils.GithubProvider+"/url", nil),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
		{
			desc:          "Success discord url",
			req:           s.MakeReq("GET", s.bp+"/"+utils.DiscordProvider+"/url", nil),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
		{
			desc:          "Success google url",
			req:           s.MakeReq("GET", s.bp+"/"+utils.GoogleProvider+"/url", nil),
			showResp:      true,
			wantStatus:    http.StatusOK,
			bodyValidator: s.CheckSuccess,
		},
	}
	s.RunRequests(testCases)
}
