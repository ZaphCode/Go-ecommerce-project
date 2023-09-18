package payment

import (
	"testing"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/repositories/user"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentmethod"

	//"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/stretchr/testify/suite"
)

type StripeServiceSuite struct {
	suite.Suite
	service PaymentService
}

//* Main

func TestStripeServiceSuite(t *testing.T) {
	suite.Run(t, new(StripeServiceSuite))
}

//* ------- Life Cycle ------------------

func (s *StripeServiceSuite) SetupSuite() {
	s.T().Log("----------- Init setup! -----------")

	config.MustLoadConfig("./../../../config")
	cfg := config.Get()

	userRepo := user.NewMemoryUserRepository(utils.UserAdmin, utils.UserExp1)
	s.service = NewStripePaymentService(cfg.Stripe.SecretKey, userRepo)

	utils.PrintBlueTesting(s.T(), "Init")
}

func (s *StripeServiceSuite) TearDownSuite() {
	utils.PrettyPrintTesting(s.T(), "Bye")
}

//* Tests

func (s *StripeServiceSuite) TestCreateAndDeleteCustomerID() {
	cusID, err := s.service.GetOrCreateCustomerID(utils.UserAdmin.ID)

	s.NoError(err, "should not be error")

	s.T().Logf("\n\n Customer ID: %s\n\n", cusID)

	err = s.service.DeleteCustomer(cusID)

	s.NoError(err, "should not be error")
}

func (s *StripeServiceSuite) TestAttachAndDetachCard() {
	pm := s.createPM()

	err := s.service.AttachCardToCustomer(
		pm,
		utils.UserExp1.CustomerID,
	)

	s.NoError(err, "should not be error")

	err = s.service.DetachCardFromCustomer(pm, utils.UserExp1.CustomerID)

	s.Require().NoError(err, "should not be error")
}

func (s *StripeServiceSuite) TestGetCards() {
	cards, err := s.service.GetCustomerCards("cus_NokAmwAreSjg1Y")

	s.NoError(err, "should not be error")

	utils.PrettyPrintTesting(s.T(), cards)
}

func (s *StripeServiceSuite) TestMakePayment() {
	pm := s.createPM()

	testCases := []struct {
		desc    string
		wantErr bool
		cusID   string
		pmID    string
	}{
		{
			desc:    "invalid payment method",
			cusID:   utils.UserExp1.CustomerID,
			wantErr: true,
			pmID:    "asd",
		},
		{
			desc:    "invalid customer id",
			cusID:   "asd",
			wantErr: true,
			pmID:    pm,
		},
		{
			desc:    "customer using card attached to other customer",
			cusID:   "cus_NomqrSHuyzac8E",
			wantErr: true,
			pmID:    "pm_1NKP27G8UXDxPRbaNZRE6Ajd",
		},
		{
			desc:    "proper work with new card",
			cusID:   utils.UserExp1.CustomerID,
			wantErr: false,
			pmID:    pm,
		},
		{
			desc:    "proper work with attached card",
			cusID:   utils.UserExp1.CustomerID,
			wantErr: false,
			pmID:    "pm_1NKP27G8UXDxPRbaNZRE6Ajd",
		},
	}
	for i, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.MakePayment(tC.cusID, tC.pmID, 4599+int64(i+1*10))

			s.Equal((err != nil), tC.wantErr, "expect err fail: %v", err)

			utils.PrintBlueTesting(s.T(), err)
		})
	}
}

func (s *StripeServiceSuite) TestDetachCard() {
	err := s.service.DetachCardFromCustomer(
		"pm_1NKPiEG8UXDxPRbaEDuh6BrU",
		utils.UserExp1.CustomerID,
	)

	s.Require().Error(err, "should be error because not exists")
}

func (s *StripeServiceSuite) createPM() string {
	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String("4242424242424242"),
			ExpMonth: stripe.Int64(4),
			ExpYear:  stripe.Int64(2025),
			CVC:      stripe.String("314"),
		},
		BillingDetails: &stripe.PaymentMethodBillingDetailsParams{
			Name: stripe.String("OMG CARD OMG"),
		},
		Type: stripe.String("card"),
	}

	pm, err := paymentmethod.New(params)

	s.NoError(err, "should not be error")

	return pm.ID
}
