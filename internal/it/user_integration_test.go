//go:build integration
// +build integration

package it

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/jmehdipour/gift-card/internal/domain"
)

func makeCreateUserRequest(requestBody string) (string, int, error) {
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/users/register", bytes.NewReader([]byte(requestBody)))
	if err != nil {
		return "", 0, err
	}

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", 0, err
	}

	defer response.Body.Close()
	var responseBody bytes.Buffer
	if _, err := io.Copy(&responseBody, response.Body); err != nil {
		return "", 0, err
	}

	return responseBody.String(), response.StatusCode, nil
}

func makeLoginRequest(requestBody string) (string, int, error) {
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/users/login", bytes.NewReader([]byte(requestBody)))
	if err != nil {
		return "", 0, err
	}

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", 0, err
	}

	defer response.Body.Close()
	var responseBody bytes.Buffer
	if _, err := io.Copy(&responseBody, response.Body); err != nil {
		return "", 0, err
	}

	return responseBody.String(), response.StatusCode, nil
}

type CreateUserIntegrationTestSuite struct {
	suite.Suite
}

func (suite *CreateUserIntegrationTestSuite) TestCreateUser_Success() {
	require := suite.Require()
	user := domain.User{
		Email:    "user@example.com",
		Password: "examplePassword",
	}
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password)

	_, statusCode, err := makeCreateUserRequest(requestBody)

	require.NoError(err)
	require.Equal(http.StatusCreated, statusCode)
}

func (suite *CreateUserIntegrationTestSuite) TestCreateUser_InvalidEmail_Failure() {
	require := suite.Require()
	email := "example.com"
	password := "examplePassword"
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	expectedResponse := `{"message":"invalid email"}`

	response, statusCode, err := makeCreateUserRequest(requestBody)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *CreateUserIntegrationTestSuite) TestCreateUser_InvalidPassword_Failure() {
	require := suite.Require()
	email := "foo@bar.com"
	password := ""
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	expectedResponse := `{"message":"invalid password"}`

	response, statusCode, err := makeCreateUserRequest(requestBody)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, statusCode)
	require.JSONEq(expectedResponse, response)
}

type LoginIntegrationTestSuite struct {
	suite.Suite
}

func (suite *LoginIntegrationTestSuite) TestLogin_Success() {
	require := suite.Require()
	email := "test1@example.com"
	password := "password"
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)

	response, statusCode, err := makeLoginRequest(requestBody)

	require.NoError(err)
	require.Equal(http.StatusOK, statusCode)
	require.NotEmpty(response)
}

func (suite *LoginIntegrationTestSuite) TestLogin_InvalidRequestBody_Failure() {
	require := suite.Require()
	password := "examplePassword"
	requestBody := fmt.Sprintf(`{"email": 10, "password": "%s"}`, password)
	expectedResponse := `{"message": "invalid request body"}`

	response, statusCode, err := makeLoginRequest(requestBody)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *LoginIntegrationTestSuite) TestLogin_UnauthorizedUser_Success() {
	require := suite.Require()
	email := "foo1@example.com"
	password := "examplePassword"
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)

	response, statusCode, err := makeLoginRequest(requestBody)

	require.NoError(err)
	require.Equal(http.StatusUnauthorized, statusCode)
	require.Empty(response)
}

func TestCreateUserIntegration(t *testing.T) {
	suite.Run(t, new(CreateUserIntegrationTestSuite))
}

func TestLoginIntegration(t *testing.T) {
	suite.Run(t, new(LoginIntegrationTestSuite))
}
