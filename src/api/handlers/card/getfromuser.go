package card

import "github.com/gofiber/fiber/v2"

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
