package utils

import (
	"context"
	"testing"

	firebase "firebase.google.com/go/v4"
	"github.com/ZaphCode/clean-arch/config"
	"github.com/stretchr/testify/suite"
)

type FirebaseSuite struct {
	suite.Suite
	app *firebase.App
}

//* Main

func TestFirebaseSuite(t *testing.T) {
	suite.Run(t, new(FirebaseSuite))
}

//* Life Cycle

func (s *FirebaseSuite) SetupSuite() {
	s.T().Log("----------- Init setup! -----------")

	config.MustLoadConfig("./../../config")
	config.MustLoadFirebaseConfig("./../../config")

	s.app = config.GetFirebaseApp()
}

func (s *FirebaseSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestDeleteFirestoreCollection" {
		client := GetFirestoreClient(s.app)

		ctx := context.TODO()

		if _, err := client.Collection("objects").Doc("A1").Create(ctx, struct {
			Name string
			Age  int
		}{"A", 1}); err != nil {
			s.T().Fatal(err.Error())
		}

		if _, err := client.Collection("objects").Doc("B2").Create(ctx, struct {
			Name string
			Age  int
		}{"B", 2}); err != nil {
			s.T().Fatal(err.Error())
		}
	}
}

//* Tests

func (s *FirebaseSuite) TestGetFirestoreClient() {
	client := GetFirestoreClient(s.app)

	s.NotNil(client, "Should be a real client")

	s.Panics(func() { GetFirestoreClient(nil) }, "I should panic")
}

func (s *FirebaseSuite) TestGetStorageClient() {
	client := GetStorageClient(s.app)

	s.NotNil(client, "Should be a real client")

	s.Panics(func() { GetStorageClient(nil) }, "I should panic")
}

func (s *FirebaseSuite) TestDeleteFirestoreCollection() {
	client := GetFirestoreClient(config.GetFirebaseApp())

	testCases := []struct {
		desc     string
		collName string
		bs       int
	}{
		{
			desc:     "random coll, 1 batch size",
			collName: "random",
			bs:       1,
		},
		{
			desc:     "coll with objects",
			collName: "objects",
			bs:       10,
		},
		{
			desc:     "testing coll, 10 batch size",
			collName: "testing",
			bs:       10,
		},
		{
			desc:     "sexo, 69 batch size",
			collName: "sexo",
			bs:       69,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			s.Require().NoError(DeleteFirestoreCollection(client, tC.collName, tC.bs), "error deleting coll")

			ss, err := client.Collection(tC.collName).Documents(context.TODO()).GetAll()

			s.Require().NoError(err, "error")

			s.Require().Len(ss, 0, "collection should be empty")
		})
	}
}
