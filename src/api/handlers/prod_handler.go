package handlers

import (
	"github.com/ZaphCode/clean-arch/src/api/dtos"
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/gofiber/fiber/v2"
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
// @Tags         products
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

// * Create product handler
// @Summary      Create new product
// @Description  Create product
// @Tags         products
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

	cat, err := h.catSvc.GetByName(body.Category)

	if err != nil {
		return h.RespErr(c, 500, "error getting category", err.Error())
	}

	if cat == nil {
		return h.RespErr(c, 406, "That category does'nt exist")
	}

	prod := body.AdaptToProduct()

	if err := h.prodSvc.Create(&prod); err != nil {
		return h.RespErr(c, 500, "error creating product", err.Error())
	}

	return h.RespOK(c, 201, "product created", prod)
}
