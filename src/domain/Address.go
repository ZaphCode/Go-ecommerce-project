package domain

import (
	"github.com/google/uuid"
)

//* Model

type Address struct {
	Model
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
	ServiceCrudOperations[Address]
	GetAllByUserID(ID uuid.UUID) ([]Address, error)
}

//* Repository

type AddressRepository interface {
	RepositoryCrudOperations[Address]
}
