package address

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
)

// GetUserAddress
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

	_, err := h.addrSvc.GetAllByUserID(ud.ID)

	if err != nil {
		return h.RespErr(c, 500, "error getting addresses", err.Error())
	}

	// TODO: implement get all address from user (remove this mock data)
	return h.RespOK(c, 200, "all address for user", []domain.Address{utils.AddrExp1, utils.AddrExp2})
}
