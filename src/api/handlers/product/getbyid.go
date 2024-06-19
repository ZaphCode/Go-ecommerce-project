package product

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Get product by ID handler
// @Summary      Get product
// @Description  Get product by id
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path string true "product   uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      200  {object}  dtos.ProductRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.DetailRespErrDTO
// @Failure      404  {object}  dtos.RespErrDTO
// @Router       /product/get/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	uid, err := uuid.Parse(id)

	if err != nil {
		return h.RespErr(c, 406, "invalid product id")
	}

	prod, err := h.prodSvc.GetByID(uid)

	if err != nil {
		return h.RespErr(c, 500, "error getting product", err.Error())
	}

	if prod == nil {
		return h.RespErr(c, 404, "product not found")
	}

	return h.RespOK(c, http.StatusFound, "product found", prod)
}
