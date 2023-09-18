package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	shared.Responder
	usrSvc  domain.UserService
	ordSvc  domain.OrderService
	pmSvc   payment.PaymentService
	prodSvc domain.ProductService
	vldSvc  validation.ValidationService
}

func NewOrderHandler(
	usrSvc domain.UserService,
	ordSvc domain.OrderService,
	prodSvc domain.ProductService,
	pmSvc payment.PaymentService,
	vldSvc validation.ValidationService,
) *OrderHandler {
	return &OrderHandler{
		usrSvc:  usrSvc,
		prodSvc: prodSvc,
		ordSvc:  ordSvc,
		pmSvc:   pmSvc,
		vldSvc:  vldSvc,
	}
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	ud, ok := c.Locals("user-data").(*auth.Claims)

	if !ok {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	os, err := h.ordSvc.GetAllByUserID(ud.ID)

	if err != nil {
		return h.RespErr(c, 500, "error getting user orders", err.Error())
	}

	return h.RespOK(c, 200, "all orders", os)
}

// * Create new order handler
// @Summary      Create new order
// @Description  Create order
// @Tags         order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_data  body dtos.NewOrderDTO true "order data"
// @Success      200  {object}  dtos.OrderDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /order/new [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	usrData, ok1 := c.Locals("user-data").(*auth.Claims)
	cusID, ok2 := c.Locals("customer-id").(string)

	if !ok1 || !ok2 {
		return h.RespErr(c, 500, "internal server error", "something went wrong")
	}

	body := dtos.NewOrderDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	price, err := h.prodSvc.CalculateTotalPrice(body.Products)

	if err != nil {
		return h.RespErr(c, 400, "some product are invalid", err.Error())
	}

	order := body.AdaptToOrder(price, usrData.ID)

	if err := h.ordSvc.Create(&order); err != nil {
		return h.RespErr(c, 500, "error creating order", err.Error())
	}

	err = h.pmSvc.MakePayment(cusID, body.PaymentID, price)

	if err != nil {
		return h.RespErr(c, 500, "error making the payment", err.Error())
	}

	go func() {
		if err = h.ordSvc.SetPaidStatus(order.ID, true); err != nil {
			utils.PrintColor("red", "Error updating order status")
		}
	}()

	return h.RespOK(c, 200, "order created", fiber.Map{
		"order": order,
	})
}
