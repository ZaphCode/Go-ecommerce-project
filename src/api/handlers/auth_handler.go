package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/email"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	shared.Responder
	usrSvc   domain.UserService
	emailSvc email.EmailService
	jwtSvc   auth.JWTService
	vldSvc   validation.ValidationService
}

func NewAuthHandler(
	usrSvc domain.UserService,
	emailSvc email.EmailService,
	jwtSvc auth.JWTService,
	vldSvc validation.ValidationService,
) *AuthHandler {
	return &AuthHandler{
		usrSvc:   usrSvc,
		emailSvc: emailSvc,
		jwtSvc:   jwtSvc,
		vldSvc:   vldSvc,
	}
}

// * Sign up handler
// @Summary      Sign up
// @Description  Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user_data body dtos.SignupDTO true "Sign up user"
// @Success      201  {object}  dtos.UserRespOKDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /auth/signup [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	body := dtos.SignupDTO{}
	cfg := config.Get()

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more field are invalid", err)
	}

	user := body.AdaptToUser()

	if err := h.usrSvc.Create(&user); err != nil {
		return h.RespErr(c, 500, "create user error", err.Error())
	}

	// Send verification email
	go func() {
		tokenCode, err := h.jwtSvc.CreateToken(
			auth.Claims{ID: user.ID, Role: user.Role},
			time.Hour*24*3, cfg.Api.VerificationSecret,
		)

		if err != nil {
			fmt.Println("Error sending email")
			return
		}

		err = h.emailSvc.SendEmail(email.EmailData{
			Email:    user.Email,
			Subject:  "Pulse | Verify your email!",
			Template: "change_password.html",
			Data: fiber.Map{
				"Name":  user.Username,
				"Email": user.Email,
				"Code":  tokenCode,
			},
		})

		if err != nil {
			fmt.Println("Error sending email")
			return
		}

		fmt.Println(">>> Email Sent to:", user.Email)
	}()

	return h.RespOK(c, 201, "sign up success", user)
}

// * Sign in handler
// @Summary      Sign in
// @Description  Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param 	     credentials body dtos.SigninDTO true "Sign in user"
// @Success      200  {object}  dtos.SignInRespOKDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      404  {object}  dtos.RespErrDTO
// @Failure      403  {object}  dtos.RespErrDTO
// @Router       /auth/signin [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	body := dtos.SigninDTO{}
	cfg := config.Get()

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	user, err := h.usrSvc.GetByEmail(body.Email)

	if err != nil {
		return h.RespErr(c, 500, "searching user error", err.Error())
	}

	if user == nil {
		return h.RespErr(c, 404, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return h.RespErr(c, 403, "invalid password")
	}

	accessToken, atErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		shared.AccessTokenExp, cfg.Api.AccessTokenSecret,
	)

	refreshToken, rtErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if atErr != nil || rtErr != nil {
		return h.RespErr(c, 500, "error creating tokens", "something went wrong")
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return h.RespOK(c, 200, "sign in successfully", fiber.Map{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// * Sign out handler
// @Summary      Sign out
// @Description  Logout user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dtos.RespOKDTO
// @Router       /auth/signout [get]
func (h *AuthHandler) SignOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     config.Get().Api.RefreshTokenCookie,
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-(time.Hour)),
		SameSite: "lax",
	})

	return h.RespOK(c, 200, "sign out successfully")
}

// * Sing in with OAuth handler
// @Summary      Sign in OAuth
// @Description  Sing in by OAuth provider (google/github/discord)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        code query string true "OAuth code"
// @Param        provider path string true "OAuth provider" Enums(google, discord, github)
// @Success      200  {object}  dtos.SignInRespOKDTO
// @Failure      406  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.RespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /auth/{provider}/callback [get]
func (h *AuthHandler) SignInWihOAuth(c *fiber.Ctx) error {
	code := c.Query("code")
	provider := c.Params("provider")
	providers := utils.GetOAuthProviders()
	cfg := config.Get()

	if code == "" {
		return h.RespErr(c, 400, "missing oauth code from query")
	}

	if !utils.ItemInSlice(provider, providers) {
		return h.RespErr(c, 406,
			"invalid oauth provider",
			fmt.Sprintf("The avalible proveders are: %s", strings.Join(providers, ", ")),
		)
	}

	var oauthSvc auth.OAuthService

	switch provider {
	case utils.GoogleProvider:
		oauthSvc = auth.NewGoogleOAuthService()
	case utils.DiscordProvider:
		oauthSvc = auth.NewDiscordOAuthService()
	case utils.GithubProvider:
		oauthSvc = auth.NewGithubOAuthService()
	}

	oauthUser, err := oauthSvc.GetOAuthUser(code)

	if err != nil {
		return h.RespErr(c, 500, "error getting user from "+provider, err.Error())
	}

	user, err := h.usrSvc.GetByEmail(oauthUser.Email)

	if err != nil {
		return h.RespErr(c, 500, "searching user error", err.Error())
	}

	if user == nil {
		if err := h.usrSvc.Create(oauthUser); err != nil {
			return h.RespErr(c, 500, "creating user error", err.Error())
		}
		user = oauthUser
	}

	if !user.VerifiedEmail {
		go func() {
			tokenCode, err := h.jwtSvc.CreateToken(
				auth.Claims{ID: user.ID, Role: user.Role},
				time.Hour*24*3, cfg.Api.VerificationSecret,
			)

			if err != nil {
				fmt.Println("Error sending email")
				return
			}

			err = h.emailSvc.SendEmail(email.EmailData{
				Email:    user.Email,
				Subject:  "Pulse | Verify your email!",
				Template: "change_password.html",
				Data: fiber.Map{
					"Name":  user.Username,
					"Email": user.Email,
					"Code":  tokenCode,
				},
			})

			if err != nil {
				fmt.Println("Error sending email")
				return
			}

			fmt.Println(">>> Email Sent to:", user.Email)
		}()
	}

	accessToken, atErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		shared.AccessTokenExp, cfg.Api.AccessTokenSecret,
	)

	refreshToken, rtErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if atErr != nil || rtErr != nil {
		return h.RespErr(c, 500, "error creating tokens", "something went wrong")
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return h.RespOK(c, 200, "sign in successfully", fiber.Map{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// * Get OAuth url handler
// @Summary      Get OAuth url
// @Description  Get OAuth url from provider (google/github/discord)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param 		 provider path string true "OAuth provider" Enums(google, discord, github)
// @Success      200  {object}  dtos.URLRespOKDTO
// @Failure      406  {object}  dtos.DetailRespErrDTO
// @Router       /auth/{provider}/url [get]
func (h *AuthHandler) GetOAuthUrl(c *fiber.Ctx) error {
	provider := c.Params("provider")
	providers := utils.GetOAuthProviders()

	if !utils.ItemInSlice(provider, providers) {
		return h.RespErr(c, 406, "invalid oauth provider", fmt.Sprintf("The avalible proveders are: %s", strings.Join(providers, ", ")))
	}

	var oauthSvc auth.OAuthService

	switch provider {
	case utils.GoogleProvider:
		oauthSvc = auth.NewGoogleOAuthService()
	case utils.DiscordProvider:
		oauthSvc = auth.NewDiscordOAuthService()
	case utils.GithubProvider:
		oauthSvc = auth.NewGithubOAuthService()
	}

	return h.RespOK(c, 200, "OAuth url for "+provider, oauthSvc.GetOAuthUrl())
}

// * Get auth user handler
// @Summary      Get auth user
// @Description  Get the current authenticated user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200  {object}  dtos.UserRespOKDTO
// @Failure      500  {object}  dtos.RespErrDTO
// @Router       /auth/me [get]
func (h *AuthHandler) GetAuthUser(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return h.RespErr(c, 500, "internal server error")
	}

	user, err := h.usrSvc.GetByID(ud.ID)

	if err != nil {
		return h.RespErr(c, 500, "internal server error")
	}

	return h.RespOK(c, 200, "auth user", user)
}

// * Refresh token handler
// @Summary      Refresh token
// @Description  Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param method query string false "Send refresh token method"
// @Success      200  {object}  dtos.TokenRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      403  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.RespErrDTO
// @Router       /auth/refresh [get]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	method := c.Query("method", "cookie")
	cfg := config.Get()

	var rt string

	switch method {
	case "cookie":
		rt = c.Cookies(cfg.Api.RefreshTokenCookie)
	case "header":
		rt = c.Get(cfg.Api.RefreshTokenHeader)
	default:
		rt = c.Cookies(cfg.Api.RefreshTokenCookie)
	}

	if rt == "" {
		return h.RespErr(c, 400, "missing refresh token", "send the token by headers or cookies")
	}

	claims, err := h.jwtSvc.DecodeToken(rt, cfg.Api.RefreshTokenSecret)

	if err != nil {
		return h.RespErr(c, 400, "invalid refresh token", err.Error())
	}

	at, err := h.jwtSvc.CreateToken(*claims, 10*time.Minute, cfg.Api.AccessTokenSecret)

	if err != nil {
		return h.RespErr(c, 500, "creating token error", err.Error())
	}

	return h.RespOK(c, 200, "token refreshed successfully", at)
}
