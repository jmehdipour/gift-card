package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/service"
)

func createUserNewEchoContext(body string) (echo.Context, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader([]byte(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(request, response)

	return ctx, response
}

func loginUserNewEchoContext(body string) (echo.Context, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader([]byte(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(request, response)

	return ctx, response
}

type CreateUserHandlerTestSuite struct {
	suite.Suite
	userService *service.UserServiceMock
}

func (suite *CreateUserHandlerTestSuite) SetupSuite() {
	suite.userService = new(service.UserServiceMock)
}

func (suite *CreateUserHandlerTestSuite) TestCreateUserHandler_Success() {
	require := suite.Require()
	user := domain.User{
		ID:       15,
		Email:    "foo@example.com",
		Password: "examplePassword",
	}
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password)

	defer suite.userService.On("CreateUser", user.Email, user.Password).Return(&user, nil).Unset()

	ctx, response := createUserNewEchoContext(requestBody)

	err := CreateUserHandler(suite.userService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusCreated, response.Code)
}

func (suite *CreateUserHandlerTestSuite) TestCreateUserHandler_InvalidEmail_Failure() {
	require := suite.Require()
	email := "example.com"
	password := "examplePassword"
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	expectedResponse := `{"message":"invalid email"}`

	ctx, response := createUserNewEchoContext(requestBody)

	err := CreateUserHandler(suite.userService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *CreateUserHandlerTestSuite) TestCreateUserHandler_InvalidPassword_Failure() {
	require := suite.Require()
	email := "foo@bar.com"
	password := ""
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	expectedResponse := `{"message":"invalid password"}`

	ctx, response := createUserNewEchoContext(requestBody)

	err := CreateUserHandler(suite.userService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *CreateUserHandlerTestSuite) TestCreateUserHandler_UserServiceError_Failure() {
	require := suite.Require()
	email := "foo@bar.com"
	password := "examplePassword"
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	expectedResponse := `{"message":"service layer error"}`
	expectedError := errors.New("service layer error")

	defer suite.userService.On("CreateUser", email, password).Return(nil, expectedError).Unset()

	ctx, response := createUserNewEchoContext(requestBody)

	err := CreateUserHandler(suite.userService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusInternalServerError, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

type LoginHandlerTestSuite struct {
	suite.Suite
	authService *service.AuthServiceMock
}

func (suite *LoginHandlerTestSuite) SetupSuite() {
	suite.authService = new(service.AuthServiceMock)
}

func (suite *LoginHandlerTestSuite) TestLoginHandler_Success() {
	require := suite.Require()
	email := "foo@example.com"
	password := "examplePassword"
	token := "exampleToken"
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)

	defer suite.authService.On("Login", email, password).Return(token, nil).Unset()

	ctx, response := loginUserNewEchoContext(requestBody)

	err := LoginHandler(suite.authService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusOK, response.Code)
}

func (suite *LoginHandlerTestSuite) TestLoginHandler_InvalidRequestBody_Failure() {
	require := suite.Require()
	password := "examplePassword"
	requestBody := fmt.Sprintf(`{"email": 10, "password": "%s"}`, password)
	expectedResponse := `{"message": "invalid request body"}`

	ctx, response := loginUserNewEchoContext(requestBody)
	err := LoginHandler(suite.authService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *LoginHandlerTestSuite) TestLoginHandler_UnauthorizedUser_Success() {
	require := suite.Require()
	email := "foo@example.com"
	password := "examplePassword"
	requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)

	defer suite.authService.On("Login", email, password).Return("", nil).Unset()

	ctx, response := loginUserNewEchoContext(requestBody)
	err := LoginHandler(suite.authService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusUnauthorized, response.Code)
}

func TestCreateUserHandler(t *testing.T) {
	suite.Run(t, new(CreateUserHandlerTestSuite))
}

func TestLoginHandlerHandler(t *testing.T) {
	suite.Run(t, new(LoginHandlerTestSuite))
}
