package order

import (
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
)

// TODO: Add documentation

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
