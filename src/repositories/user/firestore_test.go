package user

import (
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/domain/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

var repo domain.UserRepository

type FSUserRepoSuite struct {
	suite.Suite
	client *firestore.Client
}

// * Main test
func TestFirestoreUserRepoSuite(t *testing.T) {
	suite.Run(t, new(FSUserRepoSuite))
}

// Set up

func (s *FSUserRepoSuite) SetupSuite() {
	s.T().Log("Init setup!")

	config.MustLoadConfig("./../../../config")
	config.MustLoadFirebaseConfig("./../../../config")

	s.client = utils.GetFirestoreClient(config.GetFirebaseApp())

	repo = NewFirestoreUserRepository(s.client, "user_test")
}

// Shot down

func (s *FSUserRepoSuite) TearDownSuite() {
	s.T().Log("Clean up suite!")

	if err := utils.DeleteFirestoreCollection(
		s.client,
		"user_test",
		10,
	); err != nil {
		s.FailNowf("Something went wrong: %s", err.Error())
	}
}

func (s *FSUserRepoSuite) TestSaveUser() {

	u := &domain.User{
		DomainModel: shared.DomainModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		CustomerID:    "",
		Username:      "paul",
		Email:         "paul@test.com",
		Role:          utils.UserRole,
		Password:      "dafdfkdjsafdas",
		VerifiedEmail: false,
		ImageUrl:      "test",
		Age:           19,
	}

	s.Require().NoError(repo.Save(u), "Error saving user")

}

func (s *FSUserRepoSuite) TestFind() {
	_, err := repo.Find()

	s.Require().NoErrorf(err, "Should not be error: %s", err.Error())

}

func (s *FSUserRepoSuite) TestFindByID() {
	u, err := repo.FindByID(uuid.New())

	s.Require().NoErrorf(err, "Sould not be error %s", err)

	s.Require().Nil(u, "Sould be error: that user dont exists ")
}
