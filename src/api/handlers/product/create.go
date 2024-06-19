package product

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/gofiber/fiber/v2"
)

// * Create product handler
// @Summary      Create new product
// @Description  Create product
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product_data  body dtos.NewProductDTO true "product data"
// @Success      201  {object}  dtos.ProductRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /product/create [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	body := dtos.NewProductDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	prod := body.AdaptToProduct()

	if err := h.prodSvc.Create(&prod); err != nil {
		return h.RespErr(c, 500, "error creating product", err.Error())
	}

	return h.RespOK(c, 201, "product created", prod)
}
