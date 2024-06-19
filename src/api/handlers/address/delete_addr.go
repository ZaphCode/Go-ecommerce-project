package address

import (
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Delete address handler
// @Summary      Delete address
// @Description  Delete address
// @Tags         address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "address  uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      200  {object}  dtos.RespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.DetailRespErrDTO
// @Router       /address/delete/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid address id")
	}

	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	if err := h.addrSvc.Delete(uid, ud.ID); err != nil {
		return h.RespErr(c, 500, "error deleting address", err.Error())
	}

	return h.RespOK(c, 200, "address deleted")
}
