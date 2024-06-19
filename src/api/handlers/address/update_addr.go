package address

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Update address handler
// @Summary      Update address
// @Description  Update address
// @Tags         address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "address  uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Param        address_data  body dtos.UpdateAddressDTO true "address data"
// @Success      200  {object}  dtos.AddressRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /address/update/{id} [put]
func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid address id")
	}

	body := dtos.UpdateAddressDTO{}

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

	uf := body.AdaptToUpdateFields()

	if err := h.addrSvc.Update(uid, ud.ID, uf); err != nil {
		return h.RespErr(c, 500, "error updating address", err.Error())
	}

	addr, err := h.addrSvc.GetByID(uid)

	if err != nil {
		return h.RespErr(c, 500, "error getting the updated address", err.Error())
	}

	return h.RespOK(c, 200, "address saved", addr)
}
