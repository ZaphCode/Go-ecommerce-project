package auth

import (
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/gofiber/fiber/v2"
)

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
