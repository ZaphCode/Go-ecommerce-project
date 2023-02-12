package payment

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

//* Implementation

type stripeServiceImpl struct{}

//* Contructor

func NewStripePaymentService(sk string) PaymentService {
	stripe.Key = sk
	return &stripeServiceImpl{}
}

func (s *stripeServiceImpl) CreatePaymentSecret(cusID string, amount int64) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Customer: stripe.String(cusID),
	}

	pi, err := paymentintent.New(params)

	if err != nil {
		return "", fmt.Errorf("error creating payment intent: %s", err.Error())
	}

	return pi.ClientSecret, nil
}

func (s *stripeServiceImpl) CreateCustomerID(uid uuid.UUID, email string) (string, error) {
	params := &stripe.CustomerParams{
		Email:       stripe.String(email),
		Description: stripe.String("from user: " + uid.String()),
	}

	cus, err := customer.New(params)

	if err != nil {
		return "", fmt.Errorf("error creating customer: %s", err.Error())
	}

	return cus.ID, nil
}
