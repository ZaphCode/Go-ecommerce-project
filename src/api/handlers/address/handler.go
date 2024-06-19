package address

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
)

type AddressHandler struct {
	shared.Responder
	usrSvc  domain.UserService
	addrSvc domain.AddressService
	vldSvc  validation.ValidationService
}

func NewAddressHandler(
	usrSvc domain.UserService,
	addrSvc domain.AddressService,
	vldSvc validation.ValidationService,
) *AddressHandler {
	return &AddressHandler{
		usrSvc:  usrSvc,
		addrSvc: addrSvc,
		vldSvc:  vldSvc,
	}
}
