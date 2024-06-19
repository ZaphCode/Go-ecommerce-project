package address

import (
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
)

// * Get user address handler
// @Summary      Get auth user addresses
// @Description  Get all addresses from auth user
// @Tags         address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.AddressesRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /address/list [get]
func (h *AddressHandler) GetUserAddress(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	as, err := h.addrSvc.GetAllByUserID(ud.ID)

	if err != nil {
		return h.RespErr(c, 500, "error getting addresses", err.Error())
	}

	return h.RespOK(c, 200, "all address for user", as)
}
