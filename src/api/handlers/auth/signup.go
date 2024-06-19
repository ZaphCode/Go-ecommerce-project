package auth

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/gofiber/fiber/v2"
)

// * Sign up handler
// @Summary      Sign up
// @Description  Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user_data body dtos.SignupDTO true "Sign up user"
// @Success      201  {object}  dtos.UserRespOKDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /auth/signup [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	body := dtos.SignupDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more field are invalid", err)
	}

	user := body.AdaptToUser()

	if err := h.usrSvc.Create(&user); err != nil {
		return h.RespErr(c, 500, "create user error", err.Error())
	}

	// // Send verification email
	// go func() {
	// 	tokenCode, err := h.jwtSvc.CreateToken(
	// 		auth.Claims{ID: user.ID, Role: user.Role},
	// 		time.Hour*24*3, cfg.Api.VerificationSecret,
	// 	)

	// 	if err != nil {
	// 		fmt.Println("Error sending email")
	// 		return
	// 	}

	// 	err = h.emailSvc.SendEmail(email.EmailData{
	// 		Email:    user.Email,
	// 		Subject:  "Pulse | Verify your email!",
	// 		Template: "change_password.html",
	// 		Data: fiber.Map{
	// 			"Name":  user.Username,
	// 			"Email": user.Email,
	// 			"Code":  tokenCode,
	// 		},
	// 	})

	// 	if err != nil {
	// 		fmt.Println("Error sending email")
	// 		return
	// 	}

	// 	fmt.Println(">>> Email Sent to:", user.Email)
	// }()

	return h.RespOK(c, 201, "sign up success", user)
}
