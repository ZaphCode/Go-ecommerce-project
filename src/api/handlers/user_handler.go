package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

//* Handler

type UserHandler struct {
	shared.Responder
	usrSvc domain.UserService
	vldSvc validation.ValidationService
}

//* Constructor

func NewUserHandler(
	usrSvc domain.UserService,
	vldSvc validation.ValidationService,
) *UserHandler {
	return &UserHandler{
		usrSvc: usrSvc,
		vldSvc: vldSvc,
	}
}

// * Get user handler
// @Summary      Get user
// @Description  Get user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "user uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      302  {object}  dtos.UserRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      404  {object}  dtos.RespErrDTO
// @Failure      401  {object}  dtos.DetailRespErrDTO
// @Router       /user/get/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid user id")
	}

	user, err := h.usrSvc.GetByID(uid)

	if err != nil {
		return h.RespErr(c, 500, "error getting user", err.Error())
	}

	if user == nil {
		return h.RespErr(c, 404, "user not found")
	}

	return h.RespOK(c, 302, "user found", user)
}

// * Get Users handler
// @Summary      Get users
// @Description  Get all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dtos.UsersRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /user/all [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.usrSvc.GetAll()

	if err != nil {
		return h.RespErr(c, 500, "error getting users", err.Error())
	}

	return h.RespOK(c, 200, "all users", users)
}

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

// * Delete user handler
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "user uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      200  {object}  dtos.UserRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Router       /user/delete/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid user id")
	}

	if err := h.usrSvc.Delete(uid); err != nil {
		return h.RespErr(c, 500, "error deleting user", err.Error())
	}

	return h.RespOK(c, 200, "user deleted")
}
