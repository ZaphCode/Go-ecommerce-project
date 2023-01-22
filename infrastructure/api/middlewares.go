package api

import (
	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/infrastructure/services/auth"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
	"github.com/gofiber/fiber/v2"
)

//* Auth middlewares

func (s *fiberServer) authRequired(c *fiber.Ctx) error {
	cfg := config.Get()

	token := c.Get(cfg.Api.AccessTokenHeader)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Missing access token",
		})
	}

	claims, err := s.jwtSvc.DecodeToken(token, cfg.Api.AccessTokenSecret)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid token",
			Detail:  err.Error(),
		})
	}

	c.Locals("user-data", claims)

	return c.Next()
}

func (s *fiberServer) roleRequired(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ud, ok := c.Locals("user-data").(*auth.Claims)

		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
				Status:  utils.StatusErr,
				Message: "Internal server error",
			})
		}

		if ud.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(utils.RespErr{
				Status:  utils.StatusErr,
				Message: "Missing permisions",
			})
		}

		return c.Next()
	}
}
