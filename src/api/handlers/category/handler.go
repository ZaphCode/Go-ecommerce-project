package category

import (
	"github.com/ZaphCode/clean-arch/src/api/shared"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/services/validation"
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
