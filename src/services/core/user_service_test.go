package core

import (
	"github.com/ZaphCode/clean-arch/src/repositories/user"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserServiceSuite struct {
	suite.Suite
	service *userService
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) SetupSuite() {
	s.T().Logf("\n-------------- init ---------------")

	usrRepo := user.NewMemoryUserRepository(utils.UserExp2)

	s.service = &userService{
		usrRepo: usrRepo,
	}
}

func (s *UserServiceSuite) TestUserService_GetAll() {
	users, err := s.service.GetAll()

	s.NoError(err)

	s.T().Logf("%+v", users)
}
