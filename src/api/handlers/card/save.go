package card

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/gofiber/fiber/v2"
)

// * Save card handler
// @Summary      Save card
// @Description  Save card to user acount
// @Tags         card
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        card_data  body dtos.SaveCardDTO true "card data"
// @Success      200  {object}  dtos.CardsRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /card/save [post]
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
