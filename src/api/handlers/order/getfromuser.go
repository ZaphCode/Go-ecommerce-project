package order

import (
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
)

// * Get user orders handler
// @Summary      Get auth user orders
// @Description  Get all orders from auth user
// @Tags         order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.OrdersRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /order/list [get]
func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	os, err := h.ordSvc.GetAllByUserID(ud.ID)

	if err != nil {
		return h.RespErr(c, 500, "error getting user orders", err.Error())
	}

	return h.RespOK(c, 200, "all orders", os)
}
