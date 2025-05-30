package shared

import (
	"testing"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MemoryRepoSuite struct {
	suite.Suite
	repo *MemoryRepo[domain.ExampleModel]
}

//* Main

func TestMemoryRepoSuite(t *testing.T) {
	suite.Run(t, new(MemoryRepoSuite))
}

//* ------- Life Cycle ------------------

func (s *MemoryRepoSuite) SetupSuite() {
	s.T().Log("----------- Init setup! -----------")

	s.repo = NewMemoryRepo(
		utils.NewSyncMap[uuid.UUID, domain.ExampleModel](),
	)
}

func (s *MemoryRepoSuite) SetupTest() {
	s.NoError(s.repo.Save(m1), "should create the model 1")
	s.NoError(s.repo.Save(m2), "should create the model 2")
}

func (s *MemoryRepoSuite) TearDownTest() {
	s.repo.clear()
}

func (s *MemoryRepoSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestMemoryRepo_Remove" {
		s.NoError(s.repo.Save(m3), "error creating model 3")
	}
}

//* ------ All tests --------------------

func (s *MemoryRepoSuite) TestMemoryRepo_Save() {
	testCases := []struct {
		desc      string
		expectErr bool
		input     *domain.ExampleModel
	}{
		{
			desc:      "normal saving",
			expectErr: false,
			input: &domain.ExampleModel{
				Model: domain.Model{
					ID:        uuid.New(),
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				},
				Name:  "model X      ",
				Tags:  []string{"D", "F", "G"},
				Check: false,
				Num:   63,
				Float: 250.11,
			},
		},
		{
			desc:      "model that already exists",
			expectErr: true,
			input:     m1,
		},
		{
			desc:      "nil model",
			expectErr: true,
			input:     nil,
		},
	}

	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			err := s.repo.Save(tC.input)

			s.Require().Equalf((err != nil), tC.expectErr, "expect error fail: %v", err)
		})
	}
}

func (s *MemoryRepoSuite) TestMemoryRepo_Find() {
	d, err := s.repo.Find()

	s.NoError(err, "should not be error")

	utils.PrettyPrintTesting(s.T(), d)
}

func (s *MemoryRepoSuite) TestMemoryRepo_FindByID() {
	testCases := []struct {
		desc     string
		id       uuid.UUID
		wantUser bool
		wantErr  bool
	}{
		{
			desc:     "model that doest't exist",
			id:       uuid.New(),
			wantUser: false,
			wantErr:  false,
		},
		{
			desc:     "model that exists",
			id:       m2.ID,
			wantErr:  false,
			wantUser: true,
		},
	}
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			u, err := s.repo.FindByID(tC.id)

			s.Require().Equal((err != nil), tC.wantErr, "expect err fail")

			s.Require().Equal((u != nil), tC.wantUser, "expect model fail")

			utils.PrettyPrintTesting(s.T(), u)
		})
	}
}

func (s *MemoryRepoSuite) TestMemoryRepo_FindByField() {
	testCases := []struct {
		desc      string
		wantErr   bool
		wantModel bool
		field     string
		val       any
	}{
		{
			desc:      "proper work",
			wantErr:   false,
			wantModel: true,
			field:     "Name",
			val:       "model 1",
		},
		{
			desc:      "model not found",
			wantErr:   false,
			wantModel: false,
			field:     "Name",
			val:       "tomas",
		},
		{
			desc:      "field that does'nt exist",
			wantErr:   true,
			wantModel: false,
			field:     "Email",
			val:       "tomas@gmail.com",
		},
		{
			desc:      "invalid field type",
			wantErr:   true,
			wantModel: false,
			field:     "Num",
			val:       "nopor",
		},
	}
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			model, err := s.repo.FindByField(tC.field, tC.val)

			if err != nil {
				s.T().Log(err.Error())
			}

			s.Require().Equal((err != nil), tC.wantErr, "expect error fail")

			s.Require().Equal((model != nil), tC.wantModel, "expect model fail")

			if model != nil {
				utils.PrettyPrintTesting(s.T(), model)
			}

		})
	}
}

func (s *MemoryRepoSuite) TestMemoryRepo_FindWhere() {
	testCases := []struct {
		desc       string
		field      string
		cond       string
		val        any
		wantErr    bool
		wantModels bool
	}{
		{
			desc:       "Proper work: single model",
			field:      "Name",
			cond:       "==",
			val:        "model 1",
			wantErr:    false,
			wantModels: true, // [m1]
		},
		{
			desc:       "Proper work: both models",
			field:      "Num",
			cond:       "==",
			val:        143,
			wantErr:    false,
			wantModels: true, // [m1, m2]
		},
		{
			desc:       "Unexisting field",
			field:      "Email",
			cond:       "==",
			val:        "test@test.com",
			wantErr:    true,
			wantModels: false, // nil
		},
		{
			desc:       "Invalid condition",
			field:      "Name",
			cond:       "fsdfadf",
			val:        "UwU",
			wantErr:    true,
			wantModels: false, // nil
		},
		{
			desc:       "Differend type",
			field:      "Name",
			cond:       "==",
			val:        []bool{true, false, true},
			wantErr:    false,
			wantModels: true, // []
		},
	}
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			model, err := s.repo.FindWhere(tC.field, tC.cond, tC.val)

			s.Require().Equal((err != nil), tC.wantErr, "expect error fail")

			s.Require().Equal((model != nil), tC.wantModels, "expect model fail")

			if model != nil {
				utils.PrettyPrintTesting(s.T(), model)
			}
		})
	}
}

func (s *MemoryRepoSuite) TestMemoryRepo_Update() {
	testCases := []struct {
		desc      string
		id        uuid.UUID
		expectErr bool
		input     domain.UpdateFields
	}{
		{
			desc:      "model that doest't exist",
			id:        uuid.New(),
			expectErr: true,
			input:     domain.UpdateFields{"Name": "model updated"},
		},
		{
			desc:      "nil model",
			id:        m1.ID,
			expectErr: true,
			input:     nil,
		},
	}
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			err := s.repo.Update(tC.id, tC.input)

			s.Require().Equal((err != nil), tC.expectErr, "expect error fail")
		})
	}
}

func (s *MemoryRepoSuite) TestMemoryRepo_UpdateField() {
	testCases := []struct {
		desc      string
		id        uuid.UUID
		expectErr bool
		field     string
		val       any
	}{
		{
			desc:      "proper work",
			id:        m1.ID,
			expectErr: false,
			field:     "Name",
			val:       "tomas",
		},
		{
			desc:      "model that doest't exist",
			id:        uuid.New(),
			expectErr: true,
			field:     "Name",
			val:       "tomas",
		},
		{
			desc:      "field that does'nt exist",
			id:        m2.ID,
			expectErr: true,
			field:     "Email",
			val:       "tomas@gmail.com",
		},
		{
			desc:      "invalid val type",
			id:        m1.ID,
			expectErr: true,
			field:     "Num",
			val:       "random mesage",
		},
	}
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			err := s.repo.UpdateField(tC.id, tC.field, tC.val)

			if err != nil {
				s.T().Log(err.Error())
			}

			s.Require().Equal((err != nil), tC.expectErr, "expect error fail")
		})
	}
}

func (s *MemoryRepoSuite) TestMemoryRepo_Remove() {
	testCases := []struct {
		desc    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			desc:    "remove object that does'nt exist",
			id:      uuid.New(),
			wantErr: true,
		},
		{
			desc:    "remove existing object",
			id:      m3.ID,
			wantErr: false,
		},
	}
	for i, tC := range testCases {
		tC = testCases[i]
		s.Run(tC.desc, func() {
			err := s.repo.Remove(tC.id)

			if err != nil {
				s.T().Log(err.Error())
			}

			s.Require().Equal((err != nil), tC.wantErr, "expect error fail")
		})
	}
}
