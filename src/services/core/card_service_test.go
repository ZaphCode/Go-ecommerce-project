package core

import (
	"testing"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/card"
	"github.com/ZaphCode/clean-arch/src/repositories/user"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CardServiceSuite struct {
	suite.Suite
	service *cardService
}

func TestCardServiceSuite(t *testing.T) {
	suite.Run(t, new(CardServiceSuite))
}

func (s *CardServiceSuite) SetupSuite() {
	s.T().Logf("\n-------------- init ---------------")

	cardRepo := card.NewMemoryCardRepository(
		utils.CardExp1,
		utils.CardExp2,
	)
	usrRepo := user.NewMemoryUserRepository(
		utils.UserExp1,
		utils.UserExp2,
	)

	s.service = &cardService{
		cardRepo: cardRepo,
		usrRepo:  usrRepo,
	}
}

func (s *CardServiceSuite) TestCardService_Create() {
	testCases := []struct {
		desc    string
		input   domain.Card
		wantErr bool
	}{
		{
			desc: "user that does not exist",
			input: domain.Card{
				UserID:   uuid.New(),
				Name:     "random card",
				ExpMonth: 4,
				// rest...
			},
			wantErr: true,
		},
		{
			desc: "proper work",
			input: domain.Card{
				UserID:    utils.UserExp1.ID,
				Name:      "secundary card",
				ExpMonth:  4,
				Country:   "Mexico",
				ExpYear:   2032,
				Brand:     "martercard",
				Last4:     "2445",
				PaymentID: "cd_dfkaldfkjl;ak'dit",
			},
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Create(&tC.input)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
				return
			}

			s.NotZero(tC.input.ID, "should not be zero")

			utils.PrettyPrintTesting(s.T(), tC.input)
		})
	}
}

func (s *CardServiceSuite) TestCardService_GetByID() {
	testCases := []struct {
		desc      string
		id        uuid.UUID
		wantErr   bool
		wantCards bool
	}{
		{
			desc:      "not found",
			id:        uuid.New(),
			wantErr:   true,
			wantCards: false,
		},
		{
			desc:      "proper work: user with no cards",
			id:        utils.UserExp2.ID,
			wantErr:   false,
			wantCards: false,
		},
		{
			desc:      "proper work: user with cards",
			id:        utils.UserExp1.ID,
			wantErr:   false,
			wantCards: true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			cs, err := s.service.GetAllByUserID(tC.id)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
				return
			}

			s.Equal(tC.wantCards, len(cs) > 0, "should contain cards")

			utils.PrettyPrintTesting(s.T(), cs)
		})
	}
}

func (s *CardServiceSuite) TestCardService_GetAll() {
	cs, err := s.service.GetAll()

	s.NoError(err, "should not be error")

	s.Greater(len(cs), 0, "should has cards")

	utils.PrettyPrintTesting(s.T(), cs)
}

func (s *CardServiceSuite) TestCardService_Delete() {
	testCases := []struct {
		desc    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			desc:    "card not found",
			id:      uuid.New(),
			wantErr: true,
		},
		{
			desc:    "proper work: card found",
			id:      utils.CardExp1.ID,
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Delete(tC.id)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
			}
		})
	}
}
