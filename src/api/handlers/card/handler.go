package card

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/ZaphCode/clean-arch/src/services/validation"
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
