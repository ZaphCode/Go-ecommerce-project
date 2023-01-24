package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/email"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
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

// @Summary      Sign up
// @Description  Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user_data body dtos.SignupDTO true "Sign up user"
// @Success      201  {object}  dtos.RespOK[domain.User]
// @Failure      422  {object}  dtos.RespDetailErr
// @Failure      400  {object}  dtos.RespValErr
// @Failure      500  {object}  dtos.RespDetailErr
// @Router       /auth/signup [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	body := dtos.SignupDTO{}
	cfg := config.Get()

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespValErr{
			Status:  dtos.StatusErr,
			Message: "One or more fields are invalid",
			Errors:  err.(validation.ValidationErrors),
		})
	}

	user := body.AdaptToUser()

	if err := h.usrSvc.Create(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Create user error",
			Detail:  err.Error(),
		})
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

	return c.Status(fiber.StatusCreated).JSON(dtos.RespOK[domain.User]{
		Status:  dtos.StatusOK,
		Message: "Signup successfully",
		Data:    user,
	})
}

// @Summary      Sign in
// @Description  Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param credentials body dtos.SigninDTO true "Sign in user"
// @Success      200  {object}  dtos.RespOK[dtos.SignInResp]
// @Failure      422  {object}  dtos.RespDetailErr
// @Failure      400  {object}  dtos.RespValErr
// @Failure      500  {object}  dtos.RespDetailErr
// @Failure      404  {object}  dtos.RespErr
// @Failure      403  {object}  dtos.RespErr
// @Router       /auth/signin [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	body := dtos.SigninDTO{}
	cfg := config.Get()

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespValErr{
			Status:  dtos.StatusErr,
			Message: "One or more fields are invalid",
			Errors:  err.(validation.ValidationErrors),
		})
	}

	user, err := h.usrSvc.GetByEmail(body.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Searching user error",
			Detail:  err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Invalid password",
		})
	}

	accessToken, atErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Minute*10, cfg.Api.AccessTokenSecret,
	)

	refreshToken, rtErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if atErr != nil || rtErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error creating tokens",
			Detail:  "somethin went wrong",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return c.Status(fiber.StatusOK).JSON(dtos.RespOK[dtos.SignInResp]{
		Status:  dtos.StatusOK,
		Message: "Signin successfully",
		Data: dtos.SignInResp{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// @Summary      Sign out
// @Description  Logout user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dtos.RespOK[dtos.None]
// @Router       /auth/signout [get]
func (h *AuthHandler) SignOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     config.Get().Api.RefreshTokenCookie,
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-(time.Hour)),
		SameSite: "lax",
	})

	return c.Status(fiber.StatusOK).JSON(dtos.RespOK[interface{}]{
		Status:  dtos.StatusOK,
		Message: "Sign out successfully",
	})
}

// @Summary      Sign in OAuth
// @Description  Sing in by OAuth provider (google/github/discord)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param code query string true "OAuth code"
// @Param provider path string true "OAuth provider"
// @Success      200  {object}  dtos.RespOK[dtos.SignInResp]
// @Failure      406  {object}  dtos.RespDetailErr
// @Failure      400  {object}  dtos.RespErr
// @Failure      500  {object}  dtos.RespDetailErr
// @Router       /auth/{provider}/callback [get]
func (h *AuthHandler) SignInWihOAuth(c *fiber.Ctx) error {
	code := c.Query("code")
	provider := c.Params("provider")
	providers := utils.GetOAuthProviders()
	cfg := config.Get()

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespErr{
			Status:  dtos.StatusOK,
			Message: "Missing OAuth code from query",
		})
	}

	if !utils.ItemInSlice(provider, providers) {
		return c.Status(fiber.StatusNotAcceptable).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Invalid oauth provider",
			Detail:  fmt.Sprintf("The avalible proveders are: %s", strings.Join(providers, ", ")),
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error getting user from " + provider,
			Detail:  err.Error(),
		})
	}

	user, err := h.usrSvc.GetByEmail(oauthUser.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Searching user error",
			Detail:  err.Error(),
		})
	}

	if user == nil {
		if err := h.usrSvc.Create(oauthUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
				Status:  dtos.StatusErr,
				Message: "Creating user error",
				Detail:  err.Error(),
			})
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
		time.Minute*10, cfg.Api.AccessTokenSecret,
	)

	refreshToken, rtErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if atErr != nil || rtErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error creating tokens",
			Detail:  "something went wrong",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return c.Status(fiber.StatusOK).JSON(dtos.RespOK[dtos.SignInResp]{
		Status:  dtos.StatusOK,
		Message: "Signin successfully",
		Data: dtos.SignInResp{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// @Summary      Get OAuth url
// @Description  Get OAuth url from provider (google/github/discord)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param 		 provider path string true "OAuth provider"
// @Success      200  {object}  dtos.RespOK[string]
// @Failure      406  {object}  dtos.RespDetailErr
// @Router       /auth/{provider}/url [get]
func (h *AuthHandler) GetOAuthUrl(c *fiber.Ctx) error {
	provider := c.Params("provider")
	providers := utils.GetOAuthProviders()

	if !utils.ItemInSlice(provider, providers) {
		return c.Status(fiber.StatusNotAcceptable).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Invalid oauth provider",
			Detail:  fmt.Sprintf("The avalible proveders are: %s", strings.Join(providers, ", ")),
		})
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

	return c.Status(fiber.StatusOK).JSON(dtos.RespOK[string]{
		Status:  dtos.StatusOK,
		Message: "OAuth url for " + provider,
		Data:    oauthSvc.GetOAuthUrl(),
	})
}

// @Summary      Get auth user
// @Description  Get the current authenticated user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.RespOK[domain.User]
// @Failure      500  {object}  dtos.RespErr
// @Router       /auth/me [get]
func (h *AuthHandler) GetAuthUser(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Internal server error",
		})
	}

	user, err := h.usrSvc.GetByID(ud.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dtos.RespOK[domain.User]{
		Status:  dtos.StatusOK,
		Message: "Auth user",
		Data:    *user,
	})
}

// @Summary      Refresh token
// @Description  Refresh access jwtoken
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param method query string false "Send refresh token method"
// @Success      200  {object}  dtos.RespOK[string]
// @Failure      500  {object}  dtos.RespDetailErr
// @Failure      403  {object}  dtos.RespDetailErr
// @Failure      400  {object}  dtos.RespErr
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
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Missing refresh token",
		})
	}

	claims, err := h.jwtSvc.DecodeToken(rt, cfg.Api.RefreshTokenSecret)

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Invalid refresh token",
			Detail:  err.Error(),
		})
	}

	at, err := h.jwtSvc.CreateToken(*claims, 10*time.Minute, cfg.Api.AccessTokenSecret)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Creating token error",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dtos.RespOK[string]{
		Status:  dtos.StatusOK,
		Message: "Token refreshed successfully",
		Data:    at,
	})
}
