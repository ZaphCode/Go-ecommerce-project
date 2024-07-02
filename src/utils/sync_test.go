package utils

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SyncMapSuite struct {
	suite.Suite
	syncMap *SyncMap[uuid.UUID, ExampleModel]
}

type ExampleModel struct {
	ID   uuid.UUID
	Name string
	Age  int
	Gay  bool
}

var (
	m1 ExampleModel
	m2 ExampleModel
)

//* Main

func TestSyncMapSuite(t *testing.T) {
	suite.Run(t, new(SyncMapSuite))
}

//* Life cycle

func (s *SyncMapSuite) SetupSuite() {
	m1 = ExampleModel{
		// ID:   uuid.MustParse("1551f9f0-825a-438c-9307-90cbc0bd5d63"),
		ID:   uuid.New(),
		Name: "model 1",
		Age:  24,
		Gay:  false,
	}
	m2 = ExampleModel{
		// ID:   uuid.MustParse("9f44a912-40f6-4ca6-b672-4911e3453443"),
		ID:   uuid.New(),
		Name: "model 2",
		Age:  13,
		Gay:  true,
	}

	s.syncMap = NewSyncMap[uuid.UUID, ExampleModel]("example.json")
}

func (s *SyncMapSuite) TearDownSuite() {
	if _, err := os.Stat(s.syncMap.datafile); err == nil {
		err := os.RemoveAll(s.syncMap.datafile)
		s.NoError(err, "I should not crash")
	}
}

func (s *SyncMapSuite) SetupTest() {
	s.NoError(s.syncMap.Set(m1.ID, m1), "I should not crash")
	s.NoError(s.syncMap.Set(m2.ID, m2), "I should not crash")
}

func (s *SyncMapSuite) TearDownTest() {
	s.syncMap.Clear()
	s.Equal(s.syncMap.Count(), 0, "should be 0 because it was cleared")
}

//* Tests

func (s *SyncMapSuite) TestSet() {
	testCases := []struct {
		desc    string
		input   ExampleModel
		wantErr bool
	}{
		{
			desc: "Save properly",
			input: ExampleModel{
				ID:   uuid.New(),
				Name: "model 3",
				Age:  41,
				Gay:  false,
			},
			wantErr: false,
		},
		{
			desc:    "Error saving because already exists",
			input:   m1,
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.syncMap.Set(tC.input.ID, tC.input)
			s.Equal(tC.wantErr, (err != nil), "wrong result")
		})
	}
}

func (s *SyncMapSuite) TestUpdate() {
	testCases := []struct {
		desc    string
		id      uuid.UUID
		input   ExampleModel
		wantErr bool
	}{
		{
			desc: "Update properly",
			id:   m1.ID,
			input: ExampleModel{
				ID:   m1.ID,
				Name: "model 1 updated",
				Age:  41,
				Gay:  true,
			},
			wantErr: false,
		},
		{
			desc:    "Error updating because does'nt exist",
			id:      uuid.New(),
			input:   m2,
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.syncMap.Update(tC.id, tC.input)
			s.Equal(tC.wantErr, (err != nil), "wrong result")
		})
	}
}

func (s *SyncMapSuite) TestGet() {
	testCases := []struct {
		desc       string
		input      uuid.UUID
		wantOutput ExampleModel
		wantErr    bool
	}{
		{
			desc:       "Get properly",
			input:      m1.ID,
			wantOutput: m1,
			wantErr:    false,
		},
		{
			desc:       "Error getting because does'nt exist",
			input:      uuid.New(),
			wantOutput: ExampleModel{},
			wantErr:    true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			got, err := s.syncMap.Get(tC.input)
			s.Equal(tC.wantErr, err != nil, "wrong result")

			s.Equal(tC.wantOutput, got, "should be the same")
		})
	}
}

func (s *SyncMapSuite) TestCount() {
	s.Equal(2, s.syncMap.Count(), "should be 2 because two models were added")
}

func (s *SyncMapSuite) TestExists() {
	s.False(s.syncMap.exists(uuid.New()), "it should not exists")

	s.True(s.syncMap.exists(m1.ID), "should exists")
}

func (s *SyncMapSuite) TestGetAll() {
	ds, err := s.syncMap.GetAll()

	s.Require().NoError(err, "should not be error")

	s.Equal([]ExampleModel{m1, m2}, ds, "should be the same")
}

func (s *SyncMapSuite) TestRemove() {
	testCases := []struct {
		desc    string
		input   uuid.UUID
		wantErr bool
	}{
		{
			desc:    "Delete properly",
			input:   m1.ID,
			wantErr: false,
		},
		{
			desc:    "Error deleting because does'nt exist",
			input:   uuid.New(),
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.syncMap.Remove(tC.input)
			s.Equal(tC.wantErr, err != nil, "wrong result")
		})
	}
}
