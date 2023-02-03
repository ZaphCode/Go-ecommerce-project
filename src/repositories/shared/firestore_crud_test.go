package shared

import (
	"context"
	"testing"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type FirestoreRepoSuite struct {
	suite.Suite
	repo *FirestoreRepo[ExampleModel]
}

//* Main

func TestFirestoreRepoSuite(t *testing.T) {
	suite.Run(t, new(FirestoreRepoSuite))
}

//* ------- Life Cycle ------------------

func (s *FirestoreRepoSuite) SetupSuite() {
	s.T().Log("----------- Init setup! -----------")

	config.MustLoadConfig("./../../../config")
	config.MustLoadFirebaseConfig("./../../../config")

	client := utils.GetFirestoreClient(config.GetFirebaseApp())

	s.repo = NewFirestoreRepo[ExampleModel](client, "example", "model")
}

func (s *FirestoreRepoSuite) TearDownSuite() {
	s.T().Log("----------- Clean up suite! -----------")

	if err := utils.DeleteFirestoreCollection(
		s.repo.Client,
		"example",
		10,
	); err != nil {
		s.FailNowf("Something went wrong: %s", err.Error())
	}
}

func (s *FirestoreRepoSuite) SetupTest() {
	s.Require().NoError(s.repo.Save(m1), "saving model 1 error")
	s.Require().NoError(s.repo.Save(m2), "saving model 2 error")
}

func (s *FirestoreRepoSuite) TearDownTest() {
	s.NoError(utils.DeleteFirestoreCollection(
		s.repo.Client,
		"example",
		4,
	), "error deleting all this")
}

func (s *FirestoreRepoSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestRemove" {
		s.Require().NoError(s.repo.Save(m3), "error creating model 3")
	}

	if testName == "TestFindByID" {

		s.T().Log("------ Before -----", testName)

		ref := s.repo.Client.Collection(s.repo.CollName).Doc(m3.ID.String())

		_, err := ref.Create(context.TODO(), struct {
			CreatedAt string
			Counter   uint16
			Name      bool
		}{"Maritin tin tin", 25, true})

		if err != nil {
			s.FailNow(err.Error())
		}
	}
}

//* -------------- Actual Test ---------------------

func (s *FirestoreRepoSuite) TestSave() {
	testCases := []struct {
		desc      string
		expectErr bool
		input     *ExampleModel
	}{
		{
			desc:      "normal saving",
			expectErr: false,
			input: &ExampleModel{
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

func (s *FirestoreRepoSuite) TestFind() {
	d, err := s.repo.Find()

	s.NoError(err, "should not be error")

	utils.PrettyPrintTesting(s.T(), d)
}

func (s *FirestoreRepoSuite) TestFindByID() {
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
		{
			desc:     "corrupted model",
			id:       m3.ID,
			wantErr:  true,
			wantUser: false,
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

func (s *FirestoreRepoSuite) TestFindByField() {
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

func (s *FirestoreRepoSuite) TestFindWhere() {
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

func (s *FirestoreRepoSuite) TestUpdate() {
	testCases := []struct {
		desc      string
		id        uuid.UUID
		expectErr bool
		input     *ExampleModel
	}{
		{
			desc:      "model that doest't exist",
			id:        uuid.New(),
			expectErr: true,
			input:     &ExampleModel{Name: "model updated"},
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

func (s *FirestoreRepoSuite) TestUpdateField() {
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

func (s *FirestoreRepoSuite) TestRemove() {
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
