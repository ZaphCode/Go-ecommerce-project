package dtos

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type NewAddressDTO struct {
	Name       string `json:"name" validation:"required,max=20" example:"main card"`
	City       string `json:"city" validation:"required,max=30" example:"Los Angeles"`
	Country    string `json:"country" validation:"required,max=30" example:"USA"`
	PostalCode string `json:"postal_code" validation:"required,postcode_iso3166_alpha" example:"24156"`
	Line1      string `json:"line1" validation:"required,max=20" example:"Lolipop"`
	Line2      string `json:"line2" validation:"required,max=20" example:"Wolfstreet"`
	State      string `json:"state" validation:"required,max=40" example:"California"`
}

func (a NewAddressDTO) AdaptToAddress() (addr domain.Address) {
	addr.Name = a.Name
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
	Name       string `json:"name" validation:"max=20" example:"main card"`
	City       string `json:"city" validation:"max=30" example:"Los Angeles"`
	Country    string `json:"country" validation:"max=30" example:"USA"`
	PostalCode string `json:"postal_code" validation:"postcode_iso3166_alpha" example:"24156"`
	Line1      string `json:"line1" validation:"max=20" example:"Lolipop"`
	Line2      string `json:"line2" validation:"max=20" example:"Wolfstreet"`
	State      string `json:"state" validation:"max=40" example:"California"`
}

func (a UpdateAddressDTO) AdaptToAddress() (addr domain.Address) {
	addr.Name = a.Name
	addr.City = a.City
	addr.Country = a.Country
	addr.PostalCode = a.PostalCode
	addr.Line1 = a.Line1
	addr.Line2 = a.Line2
	addr.State = a.State
	return
}
