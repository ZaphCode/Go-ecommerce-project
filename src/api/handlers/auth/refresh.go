package auth

import (
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/gofiber/fiber/v2"
)

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
