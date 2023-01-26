package shared

import (
	"github.com/gofiber/fiber/v2"
)

type Responder struct{}

func (Responder) RespOK(c *fiber.Ctx, code int, msg string, data ...interface{}) error {
	res := RespOK{
		Status:  StatusOK,
		Message: msg,
	}

	if data[0] != nil {
		res.Data = data[0]
	}

	return c.Status(code).JSON(res)
}

func (Responder) RespErr(c *fiber.Ctx, code int, msg string, detail ...string) error {
	res := RespErr{
		Status:  StatusOK,
		Message: msg,
	}

	if detail[0] != "" {
		res.Detail = detail[0]
	}
	return c.Status(code).JSON(res)
}

func (Responder) RespValErr(c *fiber.Ctx, code int, msg string, errors error) error {
	res := RespErr{
		Status:  StatusOK,
		Message: msg,
		Errors:  errors,
	}
	return c.Status(code).JSON(res)
}
