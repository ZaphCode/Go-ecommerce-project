package api

import (
	"github.com/ZaphCode/clean-arch/infrastructure/api/dtos"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *fiberServer) getUser(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid user id",
		})
	}

	user, err := s.userSvc.GetByID(uid)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error getting user",
			Detail:  err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "User not found",
		})
	}

	return c.Status(fiber.StatusFound).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "User found",
		Data:    user,
	})
}

func (s *fiberServer) getUsers(c *fiber.Ctx) error {
	users, err := s.userSvc.GetAll()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error getting users",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusFound).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "All users",
		Data:    users,
	})
}

func (s *fiberServer) createUser(c *fiber.Ctx) error {
	body := dtos.NewUserDTO{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := s.validationSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err,
		})
	}

	user := body.AdaptToUser()

	if err := s.userSvc.Create(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Create user error",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "User created!",
		Data:    user,
	})
}

func (s *fiberServer) updateUser(c *fiber.Ctx) error {
	body := dtos.UpdateUserDTO{}
	id := c.Params("id")
	uid, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid user id",
		})
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error parsing the request body",
			Detail:  err.Error(),
		})
	}

	if err := s.validationSvc.Validate(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "One or more fields are invalid",
			Detail:  err,
		})
	}

	user := body.AdaptToUser()

	if err := s.userSvc.Update(uid, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Create user error",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "User updated!",
		Data:    user,
	})
}

func (s *fiberServer) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Invalid user id",
		})
	}

	if err := s.userSvc.Delete(uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.RespErr{
			Status:  utils.StatusErr,
			Message: "Error deleting user",
			Detail:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.RespOk{
		Status:  utils.StatusOk,
		Message: "User deleted",
	})
}
