package middlewares

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/gofiber/fiber/v2"
)

type PaymentMiddleware struct {
	shared.Responder
	paymSvc payment.PaymentService
}

func NewPaymentMiddleware(
	paymSvc payment.PaymentService,
) *PaymentMiddleware {
	return &PaymentMiddleware{
		paymSvc: paymSvc,
	}
}

func (m *PaymentMiddleware) CustomerIDRequired(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return m.RespErr(c, 500, "parsing user claims error")
	}

	cusID, err := m.paymSvc.GetOrCreateCustomerID(ud.ID)

	if err != nil {
		return m.RespErr(c, 500, "error getting customer id", err.Error())
	}

	c.Locals("customer-id", cusID)

	return c.Next()
}
