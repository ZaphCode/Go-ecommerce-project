package domain

import (
	"github.com/ZaphCode/clean-arch/src/domain/shared"
	"github.com/google/uuid"
)

//* Model

type Address struct {
	shared.DomainModel
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	Line1      string    `json:"line1"`
	Line2      string    `json:"line2"`
	State      string    `json:"state"`
}

//* Service

type AddressService interface {
	shared.ServiceCrudOperations[Address]
	GetAllByUserID(ID uuid.UUID) ([]Address, error)
}

//* Repository

type AddressRepository interface {
	shared.RepositoryCrudOperations[Address]
}
