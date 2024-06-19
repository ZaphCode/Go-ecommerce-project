package card

import "github.com/gofiber/fiber/v2"

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
		return h.RespErr(c, 500, "error detaching cards", err.Error())
	}

	return h.RespOK(c, 200, "card removed")
}
