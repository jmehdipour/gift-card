package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/jmehdipour/gift-card/internal/domain"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *userRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.db, suite.mock, _ = sqlmock.New()
	suite.repo = &userRepository{
		db: suite.db,
	}
}

func (suite *UserRepositoryTestSuite) TeardownTest() {
	_ = suite.db.Close()
}

func (suite *UserRepositoryTestSuite) TestNewGiftCardRepository() {
	require := suite.Require()

	db, _, _ := sqlmock.New()
	repo := NewUserRepository(db)

	require.NotNil(repo)
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	require := suite.Require()
	id := uint(101)
	u := &domain.User{
		Email:    "foo@example.com",
		Password: "securePassword",
	}

	suite.mock.ExpectExec("^INSERT INTO users").
		WithArgs(u.Email, u.Password).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))

	err := suite.repo.Create(u)

	require.NoError(err)
	require.Equal(id, u.ID)
}

func (suite *UserRepositoryTestSuite) TestCreate_Failure() {
	require := suite.Require()
	u := &domain.User{
		Email:    "foo@example.com",
		Password: "securePassword",
	}
	expectedError := errors.New("error in inserting to users table")

	suite.mock.ExpectExec("^INSERT INTO users").
		WithArgs(u.Email, u.Password).
		WillReturnError(expectedError)

	err := suite.repo.Create(u)

	require.EqualError(err, expectedError.Error())
}

func (suite *UserRepositoryTestSuite) TestCreate_LastInsertIdError_Failure() {
	require := suite.Require()
	u := &domain.User{
		Email:    "foo@example.com",
		Password: "securePassword",
	}
	expectedError := errors.New("LastInsertId error")

	suite.mock.ExpectExec("^INSERT INTO users").
		WithArgs(u.Email, u.Password).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId error")))

	err := suite.repo.Create(u)

	require.Error(err)
	require.EqualError(err, expectedError.Error())
}

func (suite *UserRepositoryTestSuite) TestFindByID_DBError_Failure() {
	require := suite.Require()
	expectedError := "database failure"
	email := "foo@example.com"

	suite.mock.ExpectQuery("^SELECT .+ FROM users").
		WithArgs(email).
		WillReturnError(errors.New("database failure"))

	result, err := suite.repo.FindByEmail(email)

	require.Error(err)
	require.EqualError(err, expectedError)
	require.Empty(result)
}

func (suite *UserRepositoryTestSuite) TestFindByID_NotFound() {
	require := suite.Require()
	email := "foo@example.com"

	suite.mock.ExpectQuery("^SELECT .+ FROM users").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	result, err := suite.repo.FindByEmail(email)
	require.NoError(err)
	require.Empty(result)
}

func (suite *UserRepositoryTestSuite) TestFindByID_Success() {
	require := suite.Require()
	email := "foo@example.com"
	expectedResult := &domain.User{
		ID:       10,
		Email:    email,
		Password: "securePassword",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(expectedResult.ID, expectedResult.Email, expectedResult.Password)
	suite.mock.ExpectQuery("^SELECT .+ FROM users").
		WithArgs(email).
		WillReturnRows(rows)

	result, err := suite.repo.FindByEmail(email)
	require.NoError(err)
	require.Equal(expectedResult, result)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
