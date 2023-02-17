package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	shared.Responder
	prodSvc domain.ProductService
	catSvc  domain.CategoryService
	vldSvc  validation.ValidationService
}

func NewProdutHandler(
	prodSvc domain.ProductService,
	catSvc domain.CategoryService,
	vldSvc validation.ValidationService,
) *ProductHandler {
	return &ProductHandler{
		prodSvc: prodSvc,
		catSvc:  catSvc,
		vldSvc:  vldSvc,
	}
}

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

	return h.RespOK(c, 302, "product found", prod)
}

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
func (h *ProductHandler) CreateProducts(c *fiber.Ctx) error {
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

// * Delete product handler
// @Summary      Delete product
// @Description  Delete product
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "product   uuid" example(3afc3021-9395-11ed-a8b6-d8bbc1a27045)
// @Success      201  {object}  dtos.RespOKDTO
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
