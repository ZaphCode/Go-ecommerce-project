package middlewares

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/gofiber/fiber/v2"
)

type PaymentMiddleware struct {
	shared.Responder
	paymSvc payment.PaymentService
	usrSvc  domain.UserService
}

func NewPaymentMiddleware(paymSvc payment.PaymentService) *PaymentMiddleware {
	return &PaymentMiddleware{paymSvc: paymSvc}
}

func (m *PaymentMiddleware) CustomerIDRequired(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return m.RespErr(c, 500, "internal server error", "parsing user claims error")
	}

	usr, err := m.usrSvc.GetByID(ud.ID)

	if usr == nil || err != nil {
		return m.RespErr(c, 500, "internal server error", "there is a problem in the claims")
	}

	customer_id := usr.CustomerID

	if customer_id == "" {
		customer_id, err = m.paymSvc.CreateCustomerID(ud.ID, usr.Email)

		if err != nil {
			return m.RespErr(c, 500, "internal server error", err.Error())
		}

		err = m.usrSvc.Update(ud.ID, domain.UpdateFields{
			"CustomerID": customer_id,
		})

		if err != nil {
			return m.RespErr(c, 500, "internal server error", err.Error())
		}
	}

	c.Locals("customer-id", customer_id)

	return c.Next()
}
