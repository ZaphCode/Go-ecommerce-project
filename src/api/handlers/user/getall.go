package user

import "github.com/gofiber/fiber/v2"

// * Get Users handler
// @Summary      Get users
// @Description  Get all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.UsersRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /user/all [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.usrSvc.GetAll()

	if err != nil {
		return h.RespErr(c, 500, "error getting users", err.Error())
	}

	return h.RespOK(c, 200, "all users", users)
}
