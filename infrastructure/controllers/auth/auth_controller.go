package user

import (
	"github.com/ZaphCode/clean-arch/domain"
	"github.com/ZaphCode/clean-arch/infrastructure/controllers/auth/dto"
	"github.com/ZaphCode/clean-arch/infrastructure/services/validation"
	"github.com/gofiber/fiber/v2"
)

type userController struct {
	userService       domain.UserService
	validationService validation.ValidationService
}

func NewUserController(
	userService domain.UserService,
	validationService validation.ValidationService,
) *userController {
	return &userController{userService, validationService}
}

func (c *userController) SignUp(ctx *fiber.Ctx) error {
	body := dto.SignupDTO{}

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	user := body.AdaptToUser()

	if err := c.validationService.Validate(&user); err != nil {
		return err
	}

	if err := c.userService.Create(&user); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
