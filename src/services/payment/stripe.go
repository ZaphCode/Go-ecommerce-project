package payment

import (
	"fmt"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
)

//! THOOSE TEST DO'NT WORK TOGETHER BECAUSE THE PM IS USED TWICE

//* Implementation

type stripeServiceImpl struct {
	usrRepo domain.UserRepository
}

//* Contructor

func NewStripePaymentService(sk string, usrRepo domain.UserRepository) PaymentService {
	stripe.Key = sk
	return &stripeServiceImpl{usrRepo: usrRepo}
}

func (s *stripeServiceImpl) MakePayment(cusID, pmID string, amount int64) error {
	pm, err := paymentmethod.Get(
		pmID,
		nil,
	)

	if err != nil {
		return err
	}

	if pm.Customer != nil {
		if pm.Customer.ID != cusID {
			return fmt.Errorf("this payment method is associated to another customer")
		}
	}

	piID, err := s.createPaymentIntent(cusID, amount)

	if err != nil {
		return err
	}

	_, err = paymentintent.Confirm(piID, &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(pm.ID),
	})

	if err != nil {
		return fmt.Errorf("error confirming the payment intent: %w", err)
	}

	return nil
}

func (s *stripeServiceImpl) GetOrCreateCustomerID(uid uuid.UUID) (string, error) {
	usr, err := s.usrRepo.FindByID(uid)

	if usr == nil || err != nil {
		return "", fmt.Errorf("error getting user")
	}

	if usr.CustomerID != "" {
		return usr.CustomerID, nil
	}

	params := &stripe.CustomerParams{
		Email:       stripe.String(usr.Email),
		Name:        stripe.String(usr.Username),
		Description: stripe.String("From user: " + uid.String()),
	}

	cus, err := customer.New(params)

	if err != nil {
		return "", fmt.Errorf("error creating customer: %s", err.Error())
	}

	if err := s.usrRepo.UpdateField(uid, "CustomerID", cus.ID); err != nil {
		return "", fmt.Errorf("error updating customer id")
	}

	return cus.ID, nil
}

func (*stripeServiceImpl) DeleteCustomer(cusID string) error {
	_, err := customer.Del(cusID, nil)

	if err != nil {
		return fmt.Errorf("error deleting customer from stripe api. %w", err)
	}

	return nil
}

func (*stripeServiceImpl) AttachCardToCustomer(cardID, cusID string) error {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(cusID),
	}

	_, err := paymentmethod.Attach(
		cardID,
		params,
	)

	if err != nil {
		return err
	}

	return nil
}

func (*stripeServiceImpl) GetCustomerCards(custID string) ([]Card, error) {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(custID),
		Type:     stripe.String("card"),
	}

	cards := []Card{}

	i := paymentmethod.List(params)

	for i.Next() {
		pm := i.PaymentMethod()

		//utils.PrettyPrint(pm)

		card := Card{
			CustomerID: pm.Customer.ID,
			Country:    pm.Card.Country,
			ExpMonth:   uint16(pm.Card.ExpMonth),
			ExpYear:    uint16(pm.Card.ExpYear),
			Brand:      string(pm.Card.Brand),
			PaymentID:  pm.ID,
			Name:       pm.BillingDetails.Name,
			Last4:      pm.Card.Last4,
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (*stripeServiceImpl) DetachCardFromCustomer(cardID, cusID string) error {
	_, err := paymentmethod.Detach(
		cardID,
		nil,
	)

	return err
}

func (s *stripeServiceImpl) createPaymentIntent(cusID string, amount int64) (string, error) {
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

	return pi.ID, nil
}
