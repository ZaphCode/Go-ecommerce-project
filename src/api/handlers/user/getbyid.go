package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Get user handler
// @Summary      Get user
// @Description  Get user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "user uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      302  {object}  dtos.UserRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      404  {object}  dtos.RespErrDTO
// @Failure      401  {object}  dtos.DetailRespErrDTO
// @Router       /user/get/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid user id")
	}

	user, err := h.usrSvc.GetByID(uid)

	if err != nil {
		return h.RespErr(c, 500, "error getting user", err.Error())
	}

	if user == nil {
		return h.RespErr(c, 404, "user not found")
	}

	return h.RespOK(c, 200, "user found", user)
}
