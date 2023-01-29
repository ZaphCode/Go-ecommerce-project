package domain

import (
	"github.com/ZaphCode/clean-arch/src/domain/shared"
	"github.com/google/uuid"
)

//* Model

type Card struct {
	shared.DomainModel
	UserID    uuid.UUID `json:"user_id"`
	Country   string    `json:"country"`
	Name      string    `json:"name"`
	ExpMonth  uint16    `json:"exp_month"`
	ExpYear   uint16    `json:"exp_year"`
	Brand     string    `json:"brand"`
	Last4     string    `json:"last4"`
	PaymentID string    `json:"payment_id"`
}

//* Service

type CardService interface {
	shared.ServiceCrudOperations[Card]
	GetAllByUserID(ID uuid.UUID) ([]Card, error)
}

//* Repository

type PaymentMethodRepository interface {
	shared.RepositoryCrudOperations[Card]
}
