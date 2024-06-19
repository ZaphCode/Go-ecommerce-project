package address

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
)

// * Create address handler
// @Summary      Create new address
// @Description  Create address
// @Tags         address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        address_data  body dtos.NewAddressDTO true "address data"
// @Success      201  {object}  dtos.AddressRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /address/create [post]
func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	body := dtos.NewAddressDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	addr := body.AdaptToAddress(ud.ID)

	if err := h.addrSvc.Create(&addr); err != nil {
		return h.RespErr(c, 500, "error saving address", err.Error())
	}

	return h.RespOK(c, 201, "address saved", addr)
}
