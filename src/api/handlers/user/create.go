package user

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/gofiber/fiber/v2"
)

// * Create User handler
// @Summary      Create user
// @Description  Create new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_data  body dtos.NewUserDTO true "user data"
// @Success      201  {object}  dtos.UserRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /user/create [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	body := dtos.NewUserDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	user := body.AdaptToUser()

	if err := h.usrSvc.Create(&user); err != nil {
		return h.RespErr(c, 500, "create user error", err.Error())
	}

	return h.RespOK(c, 201, "user created!", user)
}
