package product

import "github.com/gofiber/fiber/v2"

// * Get products handler
// @Summary      Get products
// @Description  Get all products
// @Tags         product
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.ProductsRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /product/all [get]
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	ps, err := h.prodSvc.GetAll()

	if err != nil {
		return h.RespErr(c, 500, "error getting products", err.Error())
	}

	return h.RespOK(c, 200, "all products", ps)
}
