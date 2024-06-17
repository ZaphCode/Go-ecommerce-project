package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/gofiber/fiber/v2"
)

type CardHandler struct {
	shared.Responder
	usrSvc domain.UserService
	pmSvc  payment.PaymentService
	vldSvc validation.ValidationService
}

func NewCardHandler(
	usrSvc domain.UserService,
	pmSvc payment.PaymentService,
	vldSvc validation.ValidationService,
) *CardHandler {
	return &CardHandler{
		usrSvc: usrSvc,
		pmSvc:  pmSvc,
		vldSvc: vldSvc,
	}
}

// * Get user cards handler
// @Summary      Get auth user cards
// @Description  Get all cards from auth user
// @Tags         card
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.CardsRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /card/list [get]
func (h *CardHandler) GetUserCards(c *fiber.Ctx) error {
	cusID, ok := c.Locals("customer-id").(string)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	cs, err := h.pmSvc.GetCustomerCards(cusID)

	if err != nil {
		return h.RespErr(c, 500, "error getting cards", err.Error())
	}

	return h.RespOK(c, 200, "all cards for user", cs)
}

// * Save card handler
// @Summary      Save card
// @Description  Save card to user acount
// @Tags         card
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.CardsRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /card/list [post]
func (h *CardHandler) SaveUserCard(c *fiber.Ctx) error {
	cusID, ok := c.Locals("customer-id").(string)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	body := dtos.SaveCardDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	err := h.pmSvc.AttachCardToCustomer(body.PaymentID, cusID)

	if err != nil {
		return h.RespErr(c, 500, "error saving card", "something is wrong with the payment method id")
	}

	return h.RespOK(c, 200, "card saved")
}

// * Remove card handler
// @Summary      Remove card
// @Description  Remove a card from user acount
// @Tags         card
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.CardsRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /card/list [delete]
func (h *CardHandler) RemoveUserCard(c *fiber.Ctx) error {
	cusID, ok := c.Locals("customer-id").(string)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	cardID := c.Params("id")

	err := h.pmSvc.DetachCardFromCustomer(cardID, cusID)

	if err != nil {
		return h.RespErr(c, 500, "error detacing cards", err.Error())
	}

	return h.RespOK(c, 200, "card removed")
}
