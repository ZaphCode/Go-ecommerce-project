package dtos

import (
	"github.com/google/uuid"
)

// type NewCardDTO struct {
// 	Country  string `json:"country" validate:"required,max=30" example:"USA"`
// 	Name     string `json:"name" validate:"required,max=30" example:"Main card"`
// 	ExpMonth uint16 `json:"exp_month" validate:"required,gte=0,lte=12" example:"11"`
// 	ExpYear  uint16 `json:"exp_year" validate:"required,gte=2023" example:"2025"`
// 	Brand    string `json:"brand" validate:"required,max=15" example:"visa"`
// 	Last4    string `json:"last4" validate:"required,len=4" example:"1429"`
// }

// func (dto NewCardDTO) AdaptToCard(usrID uuid.UUID) (card payment.Card) {
// 	card.Brand = dto.Brand
// 	card.Country = dto.Country
// 	card.Last4 = dto.Last4
// 	card.Name = dto.Name
// 	card.ExpMonth = dto.ExpMonth
// 	card.ExpYear = dto.ExpYear
// 	return
// }

type SaveCardDTO struct {
	PaymentID string `json:"payment_id" validate:"required" example:"pm_1NKPiEG8UXDxPRbaEDuh6BrU"`
}

type CardDTO struct {
	ID        uuid.UUID `json:"id" example:"8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"`
	Country   string    `json:"country"  example:"USA"`
	Name      string    `json:"name" example:"Main card"`
	ExpMonth  uint16    `json:"exp_month"  example:"11"`
	ExpYear   uint16    `json:"exp_year" example:"2025"`
	Brand     string    `json:"brand" example:"visa"`
	Last4     string    `json:"last4" example:"1429"`
	CreatedAt int64     `json:"created_at" example:"1674405183"`
	UpdatedAt int64     `json:"updated_at" example:"1674405181"`
	PaymentID string    `json:"payment_id" example:"pm_1NKPiEG8UXDxPRbaEDuh6BrU"`
}

// type UpdateCardDTO struct {
// 	Country  string  `json:"country,omitempty" validate:"omitempty,max=30" example:"USA"`
// 	Name     string  `json:"name,omitempty" validate:"omitempty,max=30" example:"Main card"`
// 	ExpMonth *uint16 `json:"exp_month,omitempty" validate:"omitempty,gte=0,lte=12" example:"11"`
// 	ExpYear  *uint16 `json:"exp_year,omitempty" validate:"omitempty,gte=2023" example:"2025"`
// 	Brand    string  `json:"brand,omitempty" validate:"omitempty,max=15" example:"visa"`
// 	Last4    string  `json:"last4,omitempty" validate:"omitempty,len=4" example:"1429"`
// }

// func (dto UpdateCardDTO) AdaptToUpdateFields() (addr domain.UpdateFields) {
// 	return utils.StructToMap(dto)
// }
