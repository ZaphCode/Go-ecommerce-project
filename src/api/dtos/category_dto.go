package dtos

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type NewCategoryDTO struct {
	Name string `json:"name" validate:"required,max=15" example:"clothes"`
}

func (dto NewCategoryDTO) AdaptToCategory() domain.Category {
	return domain.Category{Name: dto.Name}
}

type CategoryDTO struct {
	NewCategoryDTO
	ID        uuid.UUID `json:"id" example:"8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"`
	CreatedAt int64     `json:"created_at" example:"1674405183"`
	UpdatedAt int64     `json:"updated_at" example:"1674405181"`
}
