package middlewares

import (
	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	shared.Responder
	jwtSvc auth.JWTService
}

func NewAuthMiddleware(jwtSvc auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtSvc: jwtSvc}
}

func (m *AuthMiddleware) AuthRequired(c *fiber.Ctx) error {
	cfg := config.Get()

	token := c.Get(cfg.Api.AccessTokenHeader)

	if token == "" {
		return m.RespErr(c, 401, "missing access token", "send the token by headers")
	}

	claims, err := m.jwtSvc.DecodeToken(token, cfg.Api.AccessTokenSecret)

	if err != nil {
		return m.RespErr(c, 401, "invalid token", err.Error())
	}

	c.Locals("user-data", claims)

	return c.Next()
}

func (m *AuthMiddleware) RoleRequired(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ud, ok := c.Locals("user-data").(*auth.Claims)

		if !ok {
			return m.RespErr(c, 500, "internal server error")
		}

		if ud.Role == role || ud.Role == utils.AdminRole {
			return c.Next()
		}

		return m.RespErr(c, 403, "missing permisions")
	}
}
