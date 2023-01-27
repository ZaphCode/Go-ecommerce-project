package shared

import (
	"github.com/gofiber/fiber/v2"
)

type RespOK struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RespErr struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Errors  error  `json:"errors,omitempty"`
}

type Responder struct{}

func (Responder) RespOK(c *fiber.Ctx, code int, msg string, data ...interface{}) error {
	res := RespOK{
		Status:  StatusOK,
		Message: msg,
	}

	if len(data) >= 1 {
		res.Data = data[0]
	}

	return c.Status(code).JSON(res)
}

func (Responder) RespErr(c *fiber.Ctx, code int, msg string, detail ...string) error {
	res := RespErr{
		Status:  StatusErr,
		Message: msg,
	}

	if len(detail) >= 1 {
		res.Detail = detail[0]
	}

	return c.Status(code).JSON(res)
}

func (Responder) RespValErr(c *fiber.Ctx, code int, msg string, errors error) error {
	res := RespErr{
		Status:  StatusErr,
		Message: msg,
		Errors:  errors,
	}
	return c.Status(code).JSON(res)
}
