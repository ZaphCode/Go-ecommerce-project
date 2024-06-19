package user

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
)

type UserHandler struct {
	shared.Responder
	usrSvc domain.UserService
	vldSvc validation.ValidationService
}

func NewUserHandler(
	usrSvc domain.UserService,
	vldSvc validation.ValidationService,
) *UserHandler {
	return &UserHandler{
		usrSvc: usrSvc,
		vldSvc: vldSvc,
	}
}
