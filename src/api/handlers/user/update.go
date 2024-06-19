package user

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Update User handler
// @Summary      Update user
// @Description  Upadate existing user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "user uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Param        user_data  body dtos.UpdateUserDTO true "user data"
// @Success      200  {object}  dtos.UserRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /user/update/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid user id")
	}

	body := dtos.UpdateUserDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	uf := body.AdaptToUpdateFields()

	if err := h.usrSvc.Update(uid, uf); err != nil {
		return h.RespErr(c, 500, "create user error", err.Error())
	}

	upUsr, err := h.usrSvc.GetByID(uid)

	if err != nil {
		return h.RespErr(c, 500, "retriving updated user error", err.Error())
	}

	return h.RespOK(c, 200, "user updated!", upUsr)
}
