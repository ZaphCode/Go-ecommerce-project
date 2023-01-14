package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/domain"
	"github.com/ZaphCode/clean-arch/infrastructure/controllers/auth/dtos"
	"github.com/ZaphCode/clean-arch/infrastructure/services/auth"
	"github.com/ZaphCode/clean-arch/infrastructure/services/validation"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type authController struct {
	userSvc       domain.UserService
	validationSvc validation.ValidationService
	jwtAuthSvc    auth.JwtAuthService
}

func NewAuthController(
	userService domain.UserService,
	validationService validation.ValidationService,
	jwtAuthService auth.JwtAuthService,
) *authController {
	return &authController{
		userSvc:       userService,
		validationSvc: validationService,
		jwtAuthSvc:    jwtAuthService,
	}
}

// TODO: Add send email to verify functionality
func (c *authController) SignUp(ctx *fiber.Ctx) error {
	body := dtos.SignupDTO{}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := c.validationSvc.Validate(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err,
		})
	}

	user := body.AdaptToUser()

	if err := c.userSvc.Create(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Create user error",
			Detail:  err,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Signup successfully",
		Data:    user,
	})
}

func (c *authController) SignIn(ctx *fiber.Ctx) error {
	body := dtos.SigninDTO{}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := c.validationSvc.Validate(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err,
		})
	}

	user, err := c.userSvc.GetByEmail(body.Email)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Searching user error",
			Detail:  err,
		})
	}

	if user == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid password",
		})
	}

	accessToken, refreshToken, err := c.jwtAuthSvc.CreateTokens(
		auth.Claims{
			ID:   user.ID,
			Role: user.Role,
		},
		time.Minute*5, time.Hour*24*5,
	)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error creating tokens",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     config.Get().Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Signin successfully",
		Data: fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (c *authController) SignOut(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     config.Get().Api.RefreshTokenCookie,
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-(time.Hour)),
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Sign out successfully",
	})
}

// TODO: Add send email to verify functionality
func (c *authController) SignInWihOAuth(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	provider := ctx.Params("provider")
	providers := auth.GetOAuthProviders()

	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Missing OAuth code from query",
		})
	}

	if !utils.ItemInSlice(provider, providers) {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error getting user from " + provider,
			Detail:  err.Error(),
		})
	}

	user, err := c.userSvc.GetByEmail(oauthUser.Email)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Searching user error",
			Detail:  err.Error(),
		})
	}

	if user == nil {
		if err := c.userSvc.CreateFromOAuth(oauthUser); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
				Status:  utils.StatusErr,
				Message: "Creating user error",
				Detail:  err.Error(),
			})
		}
		user = oauthUser
	}

	accessToken, refreshToken, err := c.jwtAuthSvc.CreateTokens(
		auth.Claims{
			ID:   user.ID,
			Role: user.Role,
		},
		time.Minute*5, time.Hour*24*5,
	)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error creating tokens",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     config.Get().Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "Signin successfully",
		Data: fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (c *authController) GetOAuthUrl(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider")
	providers := auth.GetOAuthProviders()

	if !utils.ItemInSlice(provider, providers) {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
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

	return ctx.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "OAuth url for " + provider,
		Data:    oauthSvc.GetOAuthUrl(),
	})
}
