package order

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/ZaphCode/clean-arch/src/services/validation"
)

type OrderHandler struct {
	shared.Responder
	usrSvc  domain.UserService
	ordSvc  domain.OrderService
	pmSvc   payment.PaymentService
	prodSvc domain.ProductService
	vldSvc  validation.ValidationService
}

func NewOrderHandler(
	usrSvc domain.UserService,
	ordSvc domain.OrderService,
	prodSvc domain.ProductService,
	pmSvc payment.PaymentService,
	vldSvc validation.ValidationService,
) *OrderHandler {
	return &OrderHandler{
		usrSvc:  usrSvc,
		prodSvc: prodSvc,
		ordSvc:  ordSvc,
		pmSvc:   pmSvc,
		vldSvc:  vldSvc,
	}
}
