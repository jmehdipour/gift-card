package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/repository"
)

type UserServiceTestSuite struct {
	suite.Suite
	userRepo    *repository.UserRepositoryMock
	userService *userService
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.userRepo = new(repository.UserRepositoryMock)
	suite.userService = &userService{
		userRepository: suite.userRepo,
	}
}

func (suite *UserServiceTestSuite) TestNewGiftCardRepository() {
	require := suite.Require()

	repo := NewUserService(suite.userRepo)

	require.NotNil(repo)
}

func (suite *UserServiceTestSuite) TestCreateUser_Success() {
	require := suite.Require()
	user := domain.User{
		ID:       15,
		Email:    "foo@example.com",
		Password: "password",
	}

	defer suite.userRepo.On("Create", mock.Anything).Return(nil).Unset()
	userResult, err := suite.userService.CreateUser(user.Email, user.Password)

	require.NoError(err)
	require.Equal(user.ID, userResult.ID)
}

func (suite *UserServiceTestSuite) TestCreateUser_Failure() {
	require := suite.Require()
	expectedError := errors.New("repo error")
	user := domain.User{
		Email:    "foo@example.com",
		Password: "password",
	}

	defer suite.userRepo.On("Create", mock.Anything).Return(expectedError).Unset()

	userResult, err := suite.userService.CreateUser(user.Email, user.Password)

	require.Error(err)
	require.EqualError(expectedError, err.Error())
	require.Empty(userResult)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
