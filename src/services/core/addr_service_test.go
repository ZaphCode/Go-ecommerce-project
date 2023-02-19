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
				Name:   "Test address",
				UserID: uuid.New(),
				// rest...
			},
			wantErr: true,
		},
		{
			desc: "proper work",
			input: domain.Address{
				Name:       "Main house",
				City:       "CMDX",
				UserID:     utils.UserExp1.ID,
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
		desc         string
		id           uuid.UUID
		wantErr      bool
		wantCategory bool
	}{
		{
			desc:         "proper work",
			id:           utils.AddrExp2.ID,
			wantErr:      false,
			wantCategory: true,
		},
		{
			desc:         "not found",
			id:           uuid.New(),
			wantErr:      false,
			wantCategory: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			got, err := s.service.GetByID(tC.id)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			s.Equal(tC.wantCategory, (got != nil), "expect category fail")
		})
	}
}

func (s *AddressServiceSuite) TestAddressService_GetAllByUserID() {
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
		addrID  uuid.UUID
		usrID   uuid.UUID
		wantErr bool
	}{
		{
			desc:    "card not found",
			addrID:  uuid.New(),
			wantErr: true,
		},
		{
			desc:    "not owner",
			addrID:  utils.AddrExp1.ID,
			usrID:   utils.UserExp2.ID,
			wantErr: true,
		},
		{
			desc:    "proper work: card found",
			addrID:  utils.AddrExp1.ID,
			usrID:   utils.UserExp1.ID,
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Delete(tC.addrID, tC.usrID)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
			}
		})
	}
}

func (s *AddressServiceSuite) TestAddressService_Update() {
	type args struct {
		uf     domain.UpdateFields
		addrID uuid.UUID
		usrID  uuid.UUID
	}
	testCases := []struct {
		desc    string
		args    args
		wantErr bool
	}{
		{
			desc: "address not found",
			args: args{
				usrID:  uuid.New(),
				addrID: uuid.New(),
				uf:     domain.UpdateFields{},
			},
			wantErr: true,
		},
		{
			desc: "Invalid owner",
			args: args{
				addrID: utils.AddrExp1.ID,
				usrID:  uuid.New(),
				uf:     domain.UpdateFields{},
			},
			wantErr: true,
		},
		{
			desc: "proper work: card found",
			args: args{
				addrID: utils.AddrExp1.ID,
				usrID:  utils.AddrExp1.UserID,
				uf: domain.UpdateFields{
					"Name": "Random address",
				},
			},
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Update(tC.args.addrID, tC.args.usrID, tC.args.uf)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
			}
		})
	}
}
