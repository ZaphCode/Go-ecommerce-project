package product

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Update prod handler
// @Summary      Update product
// @Description  Update product
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "product   uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Param        product_data  body dtos.UpdateProductDTO true "product data"
// @Success      200  {object}  dtos.ProductRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /product/update/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid product id")
	}

	body := dtos.UpdateProductDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	uf := body.AdaptToUpdateFields()

	if err := h.prodSvc.Update(uid, uf); err != nil {
		return h.RespErr(c, 500, "error updating product", err.Error())
	}

	//TODO ---- Remove this ---------
	p, err := h.prodSvc.GetByID(uid)

	if err != nil {
		return h.RespErr(c, 500, "error getting updated product", err.Error())
	}
	//TODO --------------------------

	return h.RespOK(c, 200, "product updated", p)
}
