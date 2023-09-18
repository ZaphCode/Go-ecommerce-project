package payment

import (
	"github.com/google/uuid"
)

//* Service

type PaymentService interface {
	//CreatePaymentIntent(cusID string, amount int64) (string, error)
	GetOrCreateCustomerID(uid uuid.UUID) (string, error)
	DeleteCustomer(cusID string) error
	MakePayment(cusID, pmID string, amount int64) error
	GetCustomerCards(custID string) ([]Card, error)
	AttachCardToCustomer(cardID, cusID string) error
	DetachCardFromCustomer(cardID, cusID string) error
}

//* Models

type Card struct {
	CustomerID string `json:"customer_id"`
	Country    string `json:"country"`
	Name       string `json:"name"`
	ExpMonth   uint16 `json:"exp_month"`
	ExpYear    uint16 `json:"exp_year"`
	Brand      string `json:"brand"`
	Last4      string `json:"last4"`
	PaymentID  string `json:"-"` // todo: change for "-"
}
