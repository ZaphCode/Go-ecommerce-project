package category

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/gofiber/fiber/v2"
)

// * Create category handler
// @Summary      Create new category
// @Description  Create category
// @Tags         category
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category_data  body dtos.NewCategoryDTO true "category data"
// @Success      201  {object}  dtos.CategoryRespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      422  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.RespErrDTO
// @Failure      400  {object}  dtos.ValidationRespErrDTO
// @Router       /category/create [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	body := dtos.NewCategoryDTO{}

	if err := c.BodyParser(&body); err != nil {
		return h.RespErr(c, 422, "error parsing the request body", err.Error())
	}

	if err := h.vldSvc.Validate(&body); err != nil {
		return h.RespValErr(c, 400, "one or more fields are invalid", err)
	}

	cat := body.AdaptToCategory()

	if err := h.catSvc.Create(&cat); err != nil {
		return h.RespErr(c, 500, "error creating category", err.Error())
	}

	return h.RespOK(c, 201, "category created", cat)
}
