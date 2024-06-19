package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Delete product handler
// @Summary      Delete product
// @Description  Delete product
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "product   uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      200  {object}  dtos.RespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.DetailRespErrDTO
// @Router       /product/delete/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid product id")
	}

	if err := h.prodSvc.Delete(uid); err != nil {
		return h.RespErr(c, 500, "error deleting product", err.Error())
	}

	return h.RespOK(c, 200, "product deleted")
}
