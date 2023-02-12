package payment

import (
	"github.com/google/uuid"
)

//* Service

type PaymentService interface {
	CreatePaymentSecret(cusID string, amount int64) (string, error)
	CreateCustomerID(uid uuid.UUID, email string) (string, error)
}
