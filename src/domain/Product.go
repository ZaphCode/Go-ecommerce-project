package domain

import (
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Model

type Product struct {
	utils.DBModel
	CategoryID   uuid.UUID `json:"category_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        int64     `json:"price"`
	DiscountRate float32   `json:"discount_rate"`
	ImagesUrl    []string  `json:"images_url"`
	Tags         []string  `json:"tags"`
	Avalible     bool      `json:"avalible"`
}

//* Service

type ProductService interface {
	utils.ServiceCrudOperations[Product]
	GetLatestProds(lim int) ([]Product, error)
	GetByTags(tags ...string) ([]Product, error)
}

//* Repository

type ProductRepository interface {
	utils.RepositoryCrudOperations[Product]
	FindOrderBy(field string, ord string) ([]Product, error)
	FindWhere(field string, cond string, val any) ([]Product, error)
	UpdateField(ID uuid.UUID, field string, val any) error
}
