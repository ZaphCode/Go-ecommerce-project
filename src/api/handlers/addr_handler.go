package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddressHandler struct {
	shared.Responder
	usrSvc  domain.UserService
	addrSvc domain.AddressService
	vldSvc  validation.ValidationService
}

func NewAddressHandler(
	usrSvc domain.UserService,
	addrSvc domain.AddressService,
	vldSvc validation.ValidationService,
) *AddressHandler {
	return &AddressHandler{
		usrSvc:  usrSvc,
		addrSvc: addrSvc,
		vldSvc:  vldSvc,
	}
}

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
