package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	shared.Responder
	prodSvc domain.ProductService
	catSvc  domain.CategoryService
	vldSvc  validation.ValidationService
}

func NewCategoryHandler(
	prodSvc domain.ProductService,
	catSvc domain.CategoryService,
	vldSvc validation.ValidationService,
) *CategoryHandler {
	return &CategoryHandler{
		prodSvc: prodSvc,
		catSvc:  catSvc,
		vldSvc:  vldSvc,
	}
}

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

	return h.RespOK(c, 200, "all products", cs)
}

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

	ec, err := h.catSvc.GetByName(body.Name)

	if err != nil {
		return h.RespErr(c, 500, "error getting category", err.Error())
	}

	if ec != nil {
		return h.RespErr(c, 406, "That category already exists")
	}

	cat := body.AdaptToCategory()

	if err := h.catSvc.Create(&cat); err != nil {
		return h.RespErr(c, 500, "error creating category", err.Error())
	}

	return h.RespOK(c, 201, "category created", cat)
}

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
	id := c.Params("id")
	uid, err := uuid.Parse(id)

	if err != nil {
		return h.RespErr(c, 406, "invalid category id")
	}

	cat, err := h.catSvc.GetByID(uid)

	if err != nil {
		return h.RespErr(c, 500, "error looking for that category", err.Error())
	}

	if cat == nil {
		return h.RespErr(c, 400, "that category does'nt exist")
	}

	ps, err := h.prodSvc.GetByCategory(cat.Name)

	if err != nil {
		return h.RespErr(c, 500, "error checking for products of that category", err.Error())
	}

	if len(ps) > 0 {
		return h.RespErr(c, 400, "that category has products")
	}

	if err := h.catSvc.Delete(uid); err != nil {
		return h.RespErr(c, 500, "error creating category", err.Error())
	}

	return h.RespOK(c, 201, "category deleted")
}
