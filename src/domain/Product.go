package domain

import (
	"github.com/ZaphCode/clean-arch/src/domain/shared"
	"github.com/google/uuid"
)

//* Model

type Product struct {
	shared.DomainModel
	Category     string   `json:"category"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        int64    `json:"price"`
	DiscountRate float32  `json:"discount_rate"`
	ImagesUrl    []string `json:"images_url"`
	Tags         []string `json:"tags"`
	Avalible     bool     `json:"avalible"`
}

//* Service

type ProductService interface {
	shared.ServiceCrudOperations[Product]
	GetLatestProds(lim int) ([]Product, error)
	GetByTags(tags ...string) ([]Product, error)
	GetByCategory(c string) ([]Product, error)
	SetAvalible(ID uuid.UUID, avl bool) error
}

//* Repository

type ProductRepository interface {
	shared.RepositoryCrudOperations[Product]
	FindOrderBy(field string, ord string) ([]Product, error)
	FindWhere(field string, cond string, val any) ([]Product, error)
	UpdateField(ID uuid.UUID, field string, val any) error
}
