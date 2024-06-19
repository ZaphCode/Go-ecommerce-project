package category

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// * Delete category handler
// @Summary      Delete category
// @Description  Delete category
// @Tags         category
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "category uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      201  {object}  dtos.RespOKDTO
// @Failure      401  {object}  dtos.AuthRespErrDTO
// @Failure      400  {object}  dtos.RespErrDTO
// @Failure      500  {object}  dtos.DetailRespErrDTO
// @Failure      406  {object}  dtos.DetailRespErrDTO
// @Router       /category/delete/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.RespErr(c, 406, "invalid category id")
	}

	if err := h.catSvc.Delete(uid); err != nil {
		return h.RespErr(c, 500, "error creating category", err.Error())
	}

	return h.RespOK(c, 200, "category deleted")
}
