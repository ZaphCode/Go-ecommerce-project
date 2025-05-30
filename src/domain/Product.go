package domain

import (
	"github.com/google/uuid"
)

//* Model

type Product struct {
	Model
	Category     string   `json:"category"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        int64    `json:"price"`
	DiscountRate int64    `json:"discount_rate"`
	ImagesUrl    []string `json:"images_url"`
	Tags         []string `json:"tags"`
	Available    bool     `json:"available"`
}

//* Service

type ProductService interface {
	ServiceCrudOperations[Product]
	CalculateTotalPrice(ops []OrderProduct) (int64, error)
	GetLatestProds(lim ...int) ([]Product, error)
	GetByTags(tags ...string) ([]Product, error)
	GetByCategory(c string) ([]Product, error)
	SetAvailable(ID uuid.UUID, avl bool) error
}

//* Repository

type ProductRepository interface {
	RepositoryCrudOperations[Product]
	FindOrderBy(field string, ord string) ([]Product, error)
	FindWhere(field string, cond string, val any) ([]Product, error)
	UpdateField(ID uuid.UUID, field string, val any) error
}
