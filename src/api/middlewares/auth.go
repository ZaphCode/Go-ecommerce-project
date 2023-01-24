package middlewares

import (
	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtSvc auth.JWTService
}

func NewAuthMiddleware(jwtSvc auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtSvc}
}

func (m *AuthMiddleware) AuthRequired(c *fiber.Ctx) error {
	cfg := config.Get()

	token := c.Get(cfg.Api.AccessTokenHeader)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Missing access token",
		})
	}

	claims, err := m.jwtSvc.DecodeToken(token, cfg.Api.AccessTokenSecret)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Invalid token",
			Detail:  err.Error(),
		})
	}

	c.Locals("user-data", claims)

	return c.Next()
}

func (m *AuthMiddleware) RoleRequired(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ud, ok := c.Locals("user-data").(*auth.Claims)

		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespErr{
				Status:  dtos.StatusErr,
				Message: "Internal server error",
			})
		}

		if ud.Role == role || ud.Role == utils.AdminRole {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Missing permisions",
		})
	}
}
