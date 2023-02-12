package core

import (
	"testing"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/address"
	"github.com/ZaphCode/clean-arch/src/repositories/user"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AddressServiceSuite struct {
	suite.Suite
	service *addressService
}

func TestAddressServiceSuite(t *testing.T) {
	suite.Run(t, new(AddressServiceSuite))
}

func (s *AddressServiceSuite) SetupSuite() {
	s.T().Logf("\n-------------- init ---------------")

	addrRepo := address.NewMemoryAddressRepository(
		utils.AddrExp1,
		utils.AddrExp2,
	)

	usrRepo := user.NewMemoryUserRepository(
		utils.UserExp1,
		utils.UserExp2,
	)

	s.service = &addressService{
		addrRepo: addrRepo,
		usrRepo:  usrRepo,
	}
}

func (s *AddressServiceSuite) TestAddressService_Create() {
	testCases := []struct {
		desc    string
		input   domain.Address
		wantErr bool
	}{
		{
			desc: "user that does not exist",
			input: domain.Address{
				UserID: uuid.New(),
				Name:   "Test address",
				// rest...
			},
			wantErr: true,
		},
		{
			// TODO: FILL THIS
			desc: "proper work",
			input: domain.Address{
				Name:       "Main house",
				UserID:     utils.UserExp1.ID,
				City:       "CMDX",
				Country:    "Mexico",
				PostalCode: "13513",
				Line1:      "Pato 24",
				Line2:      "Luis Echeverria",
				State:      "DF",
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

func (s *AddressServiceSuite) TestAddressService_GetByID() {
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

func (s *AddressServiceSuite) TestAddressService_GetAll() {
	cs, err := s.service.GetAll()

	s.NoError(err, "should not be error")

	s.Greater(len(cs), 0, "should has cards")

	utils.PrettyPrintTesting(s.T(), cs)
}

func (s *AddressServiceSuite) TestAddressService_Delete() {
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
			id:      utils.AddrExp1.ID,
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
