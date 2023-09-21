package dtos

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

type NewProductDTO struct {
	Category     string   `json:"category" validation:"required" example:"clothes"`
	Name         string   `json:"name" validate:"required,min=4,max=50" example:"Black T-Shirt Addidas"`
	Description  string   `json:"description" validate:"required,min=4,max=200" example:"The best T-shirt in the world."`
	Price        int64    `json:"price" validate:"required,number,gte=0" example:"2599"`
	DiscountRate int64    `json:"discount_rate" validate:"number,gte=0,lte=100" example:"23"`
	ImagesUrl    []string `json:"images_url" validate:"required,min=1,max=10,dive,url" example:"https://example.com/image1.png,https://example.com/image2.png"`
	Tags         []string `json:"tags" validate:"required,max=6" example:"t-shirts,clothes,addidas"`
	Avalible     bool     `json:"avalible"`
}

func (dto NewProductDTO) AdaptToProduct() (prod domain.Product) {
	prod.Category = dto.Category
	prod.Name = dto.Name
	prod.Description = dto.Description
	prod.Price = dto.Price
	prod.DiscountRate = dto.DiscountRate
	prod.ImagesUrl = dto.ImagesUrl
	prod.Tags = dto.Tags
	prod.Available = dto.Avalible
	return
}

type ProductDTO struct { //? Documentation
	NewProductDTO
	ID        uuid.UUID `json:"id" example:"8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"`
	CreatedAt int64     `json:"created_at" example:"1674405183"`
	UpdatedAt int64     `json:"updated_at" example:"1674405181"`
}

type UpdateProductDTO struct {
	Category     string   `json:"category,omitempty" example:"clothes"`
	Name         string   `json:"name,omitempty" validate:"omitempty,min=4,max=50" example:"Black T-Shirt Addidas"`
	Description  string   `json:"description,omitempty" validate:"omitempty,min=4,max=200" example:"The best T-shirt in the world."`
	Price        *int64   `json:"price,omitempty" validate:"omitempty,number,gte=0" example:"2599"`
	DiscountRate *int64   `json:"discount_rate,omitempty" validate:"omitempty,number,gte=0,lte=100" example:"23"`
	ImagesUrl    []string `json:"images_url,omitempty" validate:"omitempty,min=1,max=10,dive,url" example:"https://example.com/image1.png,https://example.com/image2.png"`
	Tags         []string `json:"tags,omitempty" validate:"omitempty,max=6" example:"t-shirts,clothes,addidas"`
	Available    *bool    `json:"available,omitempty"`
}

func (dto UpdateProductDTO) AdaptToUpdateFields() domain.UpdateFields {
	return utils.StructToMap(dto)
}
