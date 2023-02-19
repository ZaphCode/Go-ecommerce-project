package dtos

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

type NewAddressDTO struct {
	Name       string `json:"name" validate:"required,max=20" example:"main card"`
	City       string `json:"city" validate:"required,max=30" example:"Los Angeles"`
	Country    string `json:"country" validate:"required,max=30" example:"USA"`
	PostalCode string `json:"postal_code" validate:"required,len=5,numeric" example:"24156"`
	Line1      string `json:"line1" validate:"required,max=20" example:"Lolipop"`
	Line2      string `json:"line2,omitempty" validate:"omitempty,max=20" example:"Wolfstreet"`
	State      string `json:"state" validate:"required,max=40" example:"California"`
}

func (a NewAddressDTO) AdaptToAddress(usrID uuid.UUID) (addr domain.Address) {
	addr.Name = a.Name
	addr.UserID = usrID
	addr.City = a.City
	addr.Country = a.Country
	addr.PostalCode = a.PostalCode
	addr.Line1 = a.Line1
	addr.Line2 = a.Line2
	addr.State = a.State
	return
}

type AddressDTO struct {
	NewAddressDTO
	ID        uuid.UUID `json:"id" example:"8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"`
	CreatedAt int64     `json:"created_at" example:"1674405183"`
	UpdatedAt int64     `json:"updated_at" example:"1674405181"`
}

type UpdateAddressDTO struct {
	Name       string `json:"name,omitempty" validate:"omitempty,max=20" example:"main card"`
	City       string `json:"city,omitempty" validate:"omitempty,max=30" example:"Los Angeles"`
	Country    string `json:"country,omitempty" validate:"omitempty,max=30" example:"USA"`
	PostalCode string `json:"postal_code,omitempty" validate:"omitempty,len=5,numeric" example:"24156"`
	Line1      string `json:"line1,omitempty" validate:"omitempty,max=20" example:"Lolipop"`
	Line2      string `json:"line2,omitempty" validate:"omitempty,max=20" example:"Wolfstreet"`
	State      string `json:"state,omitempty" validate:"omitempty,max=40" example:"California"`
}

func (dto UpdateAddressDTO) AdaptToUpdateFields() (addr domain.UpdateFields) {
	return utils.StructToMap(dto)
}
