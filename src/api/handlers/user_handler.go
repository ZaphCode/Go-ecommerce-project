package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
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

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Invalid user id",
		})
	}

	user, err := h.usrSvc.GetByID(uid)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error getting user",
			Detail:  err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "User not found",
		})
	}

	return c.Status(fiber.StatusFound).JSON(dtos.RespOK[*domain.User]{
		Status:  dtos.StatusOK,
		Message: "User found",
		Data:    user,
	})
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.usrSvc.GetAll()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error getting users",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusFound).JSON(dtos.RespOK[[]domain.User]{
		Status:  dtos.StatusOK,
		Message: "All users",
		Data:    users,
	})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	body := dtos.NewUserDTO{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err.Error(),
		})
	}

	user := body.AdaptToUser()

	if err := h.usrSvc.Create(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Create user error",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dtos.RespOK[*domain.User]{
		Status:  dtos.StatusOK,
		Message: "User created!",
		Data:    &user,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	body := dtos.UpdateUserDTO{}
	id := c.Params("id")
	uid, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Invalid user id",
		})
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err.Error(),
		})
	}

	user := body.AdaptToUser()

	if err := h.usrSvc.Update(uid, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Create user error",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dtos.RespOK[*domain.User]{
		Status:  dtos.StatusOK,
		Message: "User updated!",
		Data:    &user,
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	uid, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.RespErr{
			Status:  dtos.StatusErr,
			Message: "Invalid user id",
		})
	}

	if err := h.usrSvc.Delete(uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.RespDetailErr{
			Status:  dtos.StatusErr,
			Message: "Error deleting user",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dtos.RespOK[*bool]{
		Status:  dtos.StatusOK,
		Message: "User deleted",
	})
}
