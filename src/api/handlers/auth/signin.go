package auth

import (
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/gofiber/fiber/v2"
)

// * Sign in handler
// @Summary      Sign in
// @Description  Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param 	     credentials body dtos.SigninDTO true "Sign in user"
// @Success      200  {object}  dtos.SignInRespOKDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      404  {object}  dtos.RespErrDTO
// @Failure      401  {object}  dtos.RespErrDTO
// @Router       /auth/signin [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	body := dtos.SigninDTO{}
	cfg := config.Get()

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	user, err := h.usrSvc.GetByCredentials(body.Email, body.Password)

	if err != nil {
		return h.RespErr(c, 500, "error getting user", err.Error())
	}

	accessToken, atErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		shared.AccessTokenExp, cfg.Api.AccessTokenSecret,
	)

	refreshToken, rtErr := h.jwtSvc.CreateToken(
		auth.Claims{ID: user.ID, Role: user.Role},
		time.Hour*24*5, cfg.Api.RefreshTokenSecret,
	)

	if atErr != nil || rtErr != nil {
		return h.RespErr(c, 500, "error creating tokens", "something went wrong")
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.Api.RefreshTokenCookie,
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 5),
		SameSite: "lax",
	})

	return h.RespOK(c, 200, "sign in successfully", fiber.Map{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
