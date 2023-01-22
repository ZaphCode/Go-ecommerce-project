package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/infrastructure/api/dtos"
	"github.com/ZaphCode/clean-arch/infrastructure/services/auth"
	"github.com/ZaphCode/clean-arch/infrastructure/services/email"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

//* Auth Handlers

func (s *fiberServer) signUp(c *fiber.Ctx) error {
	body := dtos.SignupDTO{}
	cfg := config.Get()

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := s.validationSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err,
		})
	}

	user := body.AdaptToUser()

	if err := s.userSvc.Create(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Create user error",
			Detail:  err.Error(),
		})
	}

	// Send verification email
	go func() {
		tokenCode, err := s.jwtSvc.CreateToken(
			auth.Claims{ID: user.ID, Role: user.Role},
			time.Hour*24*3, cfg.Api.VerificationSecret,
		)

		if err != nil {
			fmt.Println("Error sending email")
			return
		}

		if err := s.emailSvc.SendEmail(email.EmailData{
			Email:    user.Email,
			Subject:  "Pulse | Verify your email!",
			Template: "change_password.html",
			Data: fiber.Map{
				"Name":  user.Username,
				"Email": user.Email,
				"Code":  tokenCode,
			},
		}); err != nil {
			fmt.Println("Error sending email")
			return
		}

		fmt.Println(">>> Email Sent to:", user.Email)
	}()

	return c.Status(fiber.StatusCreated).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Signup successfully",
		Data:    user,
	})
}

func (s *fiberServer) signIn(c *fiber.Ctx) error {
	body := dtos.SigninDTO{}
	cfg := config.Get()

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := s.validationSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err,
		})
	}

	user, err := s.userSvc.GetByEmail(body.Email)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Searching user error",
			Detail:  err,
		})
	}

	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid password",
		})
	}

	accessToken, atErr := s.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Minute*5, cfg.Api.AccessTokenSecret,
	)

	refreshToken, rtErr := s.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if atErr != nil || rtErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error creating tokens",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return c.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Signin successfully",
		Data: fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (s *fiberServer) signOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     config.Get().Api.RefreshTokenCookie,
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-(time.Hour)),
		SameSite: "lax",
	})

	return c.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Sign out successfully",
	})
}

func (s *fiberServer) signInWihOAuth(c *fiber.Ctx) error {
	code := c.Query("code")
	provider := c.Params("provider")
	providers := auth.GetOAuthProviders()
	cfg := config.Get()

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Missing OAuth code from query",
		})
	}

	if !utils.ItemInSlice(provider, providers) {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid oauth provider",
			Detail:  fmt.Sprintf("The avalible proveders are: %s", strings.Join(providers, ", ")),
		})
	}

	var oauthSvc auth.OAuthService

	switch provider {
	case auth.GoogleProvider:
		oauthSvc = auth.NewGoogleOAuthService()
	case auth.DiscordProvider:
		oauthSvc = auth.NewDiscordOAuthService()
	case auth.GithubProvider:
		oauthSvc = auth.NewGithubOAuthService()
	}

	oauthUser, err := oauthSvc.GetOAuthUser(code)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error getting user from " + provider,
			Detail:  err.Error(),
		})
	}

	user, err := s.userSvc.GetByEmail(oauthUser.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Searching user error",
			Detail:  err.Error(),
		})
	}

	if user == nil {
		if err := s.userSvc.Create(oauthUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
				Status:  utils.StatusErr,
				Message: "Creating user error",
				Detail:  err.Error(),
			})
		}
		user = oauthUser
	}

	if !user.VerifiedEmail {
		go func() {
			tokenCode, err := s.jwtSvc.CreateToken(
				auth.Claims{ID: user.ID, Role: user.Role},
				time.Hour*24*3, cfg.Api.VerificationSecret,
			)

			if err != nil {
				fmt.Println("Error sending email")
				return
			}

			if err := s.emailSvc.SendEmail(email.EmailData{
				Email:    user.Email,
				Subject:  "Pulse | Verify your email!",
				Template: "change_password.html",
				Data: fiber.Map{
					"Name":  user.Username,
					"Email": user.Email,
					"Code":  tokenCode,
				},
			}); err != nil {
				fmt.Println("Error sending email")
				return
			}

			fmt.Println(">>> Email Sent to:", user.Email)
		}()
	}

	accessToken, atErr := s.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Minute*5, cfg.Api.AccessTokenSecret,
	)

	refreshToken, rtErr := s.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if atErr != nil || rtErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error creating tokens",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return c.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Signin successfully",
		Data: fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (s *fiberServer) getOAuthUrl(c *fiber.Ctx) error {
	provider := c.Params("provider")
	providers := auth.GetOAuthProviders()

	if !utils.ItemInSlice(provider, providers) {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid oauth provider",
			Detail:  fmt.Sprintf("The avalible proveders are: %s", strings.Join(providers, ", ")),
		})
	}

	var oauthSvc auth.OAuthService

	switch provider {
	case auth.GoogleProvider:
		oauthSvc = auth.NewGoogleOAuthService()
	case auth.DiscordProvider:
		oauthSvc = auth.NewDiscordOAuthService()
	case auth.GithubProvider:
		oauthSvc = auth.NewGithubOAuthService()
	}

	return c.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "OAuth url for " + provider,
		Data:    oauthSvc.GetOAuthUrl(),
	})
}

func (s *fiberServer) getAuthUser(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Internal server error",
		})
	}

	user, err := s.userSvc.GetByID(ud.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Auth user",
		Data:    user,
	})
}

func (s *fiberServer) refreshToken(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Missing refresh token",
		})
	}

	claims, err := s.jwtSvc.DecodeToken(rt, cfg.Api.RefreshTokenSecret)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid refresh token",
			Detail:  err.Error(),
		})
	}

	at, err := s.jwtSvc.CreateToken(*claims, 5*time.Minute, cfg.Api.AccessTokenSecret)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Creating token error",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Token refreshed successfully",
		Data:    at,
	})
}
