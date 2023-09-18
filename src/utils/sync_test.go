package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SyncMapSuite struct {
	suite.Suite
	syancMap *SyncMap[uuid.UUID, ExampleModel]
}

type ExampleModel struct {
	ID   uuid.UUID
	Name string
	Age  int
	Gay  bool
}

var m1 = ExampleModel{
	ID:   uuid.MustParse("1551f9f0-825a-438c-9307-90cbc0bd5d63"),
	Name: "model 1",
	Age:  24,
	Gay:  false,
}

var m2 = ExampleModel{
	ID:   uuid.MustParse("9f44a912-40f6-4ca6-b672-4911e3453443"),
	Name: "model 2",
	Age:  13,
	Gay:  true,
}

//* Main

func TestSyncMapSuite(t *testing.T) {
	suite.Run(t, new(SyncMapSuite))
}

//* Life cycle

func (s *SyncMapSuite) SetupSuite() {
	s.syancMap = NewSyncMap[uuid.UUID, ExampleModel]()
}

func (s *SyncMapSuite) SetupTest() {
	s.NoError(s.syancMap.Set(m1.ID, m1), "I should not crash")
	s.NoError(s.syancMap.Set(m2.ID, m2), "I should not crash")
}

func (s *SyncMapSuite) TearDownTest() {
	s.syancMap.Clear()
	s.Equal(s.syancMap.Count(), 0, "should be 0 because it was cleared")
}

//* Tests

func (s *SyncMapSuite) TestSet() {
	testCases := []struct {
		desc    string
		intput  ExampleModel
		wantErr bool
	}{
		{
			desc: "Save properly",
			intput: ExampleModel{
				ID:   uuid.New(),
				Name: "model 3",
				Age:  41,
				Gay:  false,
			},
			wantErr: false,
		},
		{
			desc:    "Error saving because already exists",
			intput:  m1,
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.syancMap.Set(tC.intput.ID, tC.intput)
			s.Equal(tC.wantErr, (err != nil), "wrong result")
		})
	}
}

func (s *SyncMapSuite) TestUpdate() {
	testCases := []struct {
		desc    string
		id      uuid.UUID
		intput  ExampleModel
		wantErr bool
	}{
		{
			desc: "Update properly",
			id:   m1.ID,
			intput: ExampleModel{
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
			intput:  m2,
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.syancMap.Update(tC.id, tC.intput)
			s.Equal(tC.wantErr, (err != nil), "wrong result")
		})
	}
}

func (s *SyncMapSuite) TestGet() {
	testCases := []struct {
		desc       string
		intput     uuid.UUID
		wantOutput ExampleModel
		wantErr    bool
	}{
		{
			desc:       "Get properly",
			intput:     m1.ID,
			wantOutput: m1,
			wantErr:    false,
		},
		{
			desc:       "Error getting because does'nt exist",
			intput:     uuid.New(),
			wantOutput: ExampleModel{},
			wantErr:    true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			got, err := s.syancMap.Get(tC.intput)
			s.Equal(tC.wantErr, (err != nil), "wrong result")

			s.Equal(tC.wantOutput, got, "should be the same")
		})
	}
}

func (s *SyncMapSuite) TestCount() {
	s.Equal(2, s.syancMap.Count(), "should be 2 because two models were added")
}

func (s *SyncMapSuite) TestExists() {
	s.False(s.syancMap.Exists(uuid.New()), "it should not exists")

	s.True(s.syancMap.Exists(m1.ID), "should exists")
}

func (s *SyncMapSuite) TestGetAll() {
	ds, err := s.syancMap.GetAll()

	s.Require().NoError(err, "should not be error")

	s.Equal([]ExampleModel{m1, m2}, ds, "should be the same")
}

func (s *SyncMapSuite) TestRemove() {
	testCases := []struct {
		desc    string
		intput  uuid.UUID
		wantErr bool
	}{
		{
			desc:    "Delete properly",
			intput:  m1.ID,
			wantErr: false,
		},
		{
			desc:    "Error deleting because does'nt exist",
			intput:  uuid.New(),
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.syancMap.Remove(tC.intput)
			s.Equal(tC.wantErr, (err != nil), "wrong result")
		})
	}
}
