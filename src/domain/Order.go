package domain

import (
	"github.com/google/uuid"
)

//* Model

type Order struct {
	Model
	UserID    uuid.UUID      `json:"user_id"`
	AddressID uuid.UUID      `json:"address_id"`
	PaymentID string         `json:"payment_id"`
	Amount    int64          `json:"amount"`
	Status    string         `json:"status"`
	Paid      bool           `json:"paid"`
	Products  []OrderProduct `json:"products"`
}

type OrderProduct struct {
	ID       uuid.UUID `json:"product_id"`
	Quantity uint      `json:"quantity"`
}

//* Service

type OrderService interface {
	Create(ord *Order) error
	GetAll() ([]Order, error)
	GetAllByUserID(ursID uuid.UUID) ([]Order, error)
	UpdateStatus(ID uuid.UUID, status string) error
	SetPaidStatus(ID uuid.UUID, paid bool) error
	Delete(ID uuid.UUID) error
}

//* Repository

type OrderRepository interface {
	Save(ord *Order) error
	Find() ([]Order, error)
	FindWhere(field, cond string, val any) ([]Order, error)
	UpdateField(ID uuid.UUID, field string, val any) error
	Remove(ID uuid.UUID) error
}
