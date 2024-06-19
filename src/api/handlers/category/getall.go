package category

import "github.com/gofiber/fiber/v2"

// * Get categories handler
// @Summary      Get categories
// @Description  Get all categories
// @Tags         category
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.CategoriesRespOKDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Router       /category/all [get]
func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	cs, err := h.catSvc.GetAll()

	if err != nil {
		return h.RespErr(c, 500, "error getting categories", err.Error())
	}

	return h.RespOK(c, 200, "all categories", cs)
}
