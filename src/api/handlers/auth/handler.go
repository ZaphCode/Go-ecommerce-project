package auth

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/email"
	"github.com/ZaphCode/clean-arch/src/services/validation"
)

type AuthHandler struct {
	shared.Responder
	usrSvc   domain.UserService
	emailSvc email.EmailService
	jwtSvc   auth.JWTService
	vldSvc   validation.ValidationService
}

func NewAuthHandler(
	usrSvc domain.UserService,
	emailSvc email.EmailService,
	jwtSvc auth.JWTService,
	vldSvc validation.ValidationService,
) *AuthHandler {
	return &AuthHandler{
		usrSvc:   usrSvc,
		emailSvc: emailSvc,
		jwtSvc:   jwtSvc,
		vldSvc:   vldSvc,
	}
}
