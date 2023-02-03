package shared

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MemoryRepoSuite struct {
	suite.Suite
	repo *MemoryRepo[ExampleModel]
}

//* Main

func TestMemoryRepoSuite(t *testing.T) {
	suite.Run(t, new(MemoryRepoSuite))
}

//* ------- Life Cycle ------------------

func (s *MemoryRepoSuite) SetupSuite() {
	s.T().Log("----------- Init setup! -----------")

}

func (s *MemoryRepoSuite) TearDownSuite() {

}

func (s *MemoryRepoSuite) SetupTest() {
}

func (s *MemoryRepoSuite) TearDownTest() {

}

func (s *MemoryRepoSuite) BeforeTest(suiteName, testName string) {

}
