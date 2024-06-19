package auth

import (
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
)

// * Get auth user handler
// @Summary      Get auth user
// @Description  Get the current authenticated user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200  {object}  dtos.UserRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
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
