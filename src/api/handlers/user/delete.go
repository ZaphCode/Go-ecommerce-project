package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Delete user handler
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "user uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      200  {object}  dtos.UserRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Router       /user/delete/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid user id")
	}

	if err := h.usrSvc.Delete(uid); err != nil {
		return h.RespErr(c, 500, "error deleting user", err.Error())
	}

	return h.RespOK(c, 200, "user deleted")
}
