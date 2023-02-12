package domain

import (
	"github.com/google/uuid"
)

//* Model

type Card struct {
	Model
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
	ServiceCrudOperations[Card]
	GetAllByUserID(ID uuid.UUID) ([]Card, error)
}

//* Repository

type CardRepository interface {
	RepositoryCrudOperations[Card]
	FindWhere(fld, cond string, val interface{}) ([]Card, error)
}
