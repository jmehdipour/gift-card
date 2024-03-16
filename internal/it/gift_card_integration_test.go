//go:build integration
// +build integration

package it

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/interface/http/handlers"
)

func makeCreateGiftCardRequest(token, requestBody string) (string, int, error) {
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/gift-cards", bytes.NewReader([]byte(requestBody)))
	if err != nil {
		return "", 0, err
	}

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, token)

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

func makeUpdateGiftCardRequest(giftCardID int, token, requestBody string) (string, int, error) {
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080/gift-cards/%d/status", giftCardID), bytes.NewReader([]byte(requestBody)))
	if err != nil {
		return "", 0, err
	}

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, token)

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

func makeGetReceivedGiftCardsRequest(token string, status int) (string, int, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/gift-cards/received?status=%d", status), nil)
	if err != nil {
		return "", 0, err
	}

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, token)

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

func makeGetSentGiftCardsRequest(token string, status int) (string, int, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/gift-cards/sent?status=%d", status), nil)
	if err != nil {
		return "", 0, err
	}

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, token)

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

func loginUser(email, password string) (string, error) {
	jsonString, statusCode, err := makeLoginRequest(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		fmt.Println("status code: ", statusCode)
		return "", errors.New("login was not ok")
	}

	var loginResponse handlers.LoginHandlerResponse
	err = json.Unmarshal([]byte(jsonString), &loginResponse)
	if err != nil {
		return "", err
	}

	return loginResponse.Token, nil
}

type GiftCardsIntegrationTestSuite struct {
	suite.Suite
	Token string
}

func (suite *GiftCardsIntegrationTestSuite) SetupSuite() {
	require := suite.Require()

	token, err := loginUser("test0@example.com", "password")
	require.NoError(err)

	suite.Token = token
}

func (suite *GiftCardsIntegrationTestSuite) TestCreateGiftCard_Success() {
	require := suite.Require()
	requestBody := `{"amount": 100, "giftee_id": 20}`
	expectedResponse := `{"id":3,"amount":100,"status":2,"gifter_id":1,"giftee_id":20}`

	response, statusCode, err := makeCreateGiftCardRequest(suite.Token, requestBody)

	require.NoError(err)
	require.JSONEq(expectedResponse, response)
	require.Equal(http.StatusCreated, statusCode)
}

func (suite *GiftCardsIntegrationTestSuite) TestCreateGiftCard_InvalidRequestBody_Failure() {
	require := suite.Require()
	requestBody := `{"amount":, "giftee_id": 20}`
	expectedResponse := `{"message":"code=400, message=Syntax error: offset=11, error=invalid character ',' looking for beginning of value, internal=invalid character ',' looking for beginning of value"}`

	response, statusCode, err := makeCreateGiftCardRequest(suite.Token, requestBody)

	require.NoError(err)
	require.JSONEq(expectedResponse, response)
	require.Equal(http.StatusBadRequest, statusCode)
}

func (suite *GiftCardsIntegrationTestSuite) TestUpdateGiftCard_Success() {
	require := suite.Require()
	requestBody := `{"status": 1}`

	response, statusCode, err := makeUpdateGiftCardRequest(1, suite.Token, requestBody)

	require.NoError(err)
	require.Equal(http.StatusOK, statusCode)
	require.Empty(response)
}

func (suite *GiftCardsIntegrationTestSuite) TestUpdateGiftCard_InvalidRequestBody_Failure() {
	require := suite.Require()
	requestBody := `{"status": "foo"}`
	expectedResponse := `{"message":"Invalid request body"}`

	response, statusCode, err := makeUpdateGiftCardRequest(1, suite.Token, requestBody)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestUpdateGiftCard_GiftCardNotFound_Failure() {
	require := suite.Require()
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	expectedResponse := `{"message": "Gift card not found"}`

	response, statusCode, err := makeUpdateGiftCardRequest(101, suite.Token, requestBody)

	require.NoError(err)
	require.Equal(http.StatusNotFound, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestUpdateGiftCard_UnauthorizedUser_Failure() {
	require := suite.Require()
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	expectedResponse := `{"message": "forbidden: user 2 is not the receiver of gift card 1"}`

	token, err := loginUser("test1@example.com", "password")
	require.NoError(err)

	response, statusCode, err := makeUpdateGiftCardRequest(1, token, requestBody)

	require.NoError(err)
	require.Equal(http.StatusForbidden, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestGetReceivedGiftCards_AcceptedStatus_Success() {
	require := suite.Require()
	expectedResponse := `{"gift_cards":[{"id":1,"amount":100,"status":0,"gifter_id":1,"giftee_id":1}, {"id":2,"amount":100,"status":0,"gifter_id":1,"giftee_id":1}],"total":2,"page":1}`

	response, statusCode, err := makeGetReceivedGiftCardsRequest(suite.Token, int(domain.GCSAccepted))

	require.NoError(err)
	require.Equal(http.StatusOK, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestGetReceivedGiftCards_RejectedStatus_Success() {
	require := suite.Require()
	expectedResponse := `{"gift_cards":[{"id":3,"amount":100,"status":1,"gifter_id":1,"giftee_id":1}, {"id":4,"amount":100,"status":1,"gifter_id":1,"giftee_id":1}],"total":2,"page":1}`

	response, statusCode, err := makeGetReceivedGiftCardsRequest(suite.Token, int(domain.GCSRejected))

	require.NoError(err)
	require.Equal(http.StatusOK, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestGetReceivedGiftCards_InvalidGiftCardStatus_Failure() {
	require := suite.Require()
	expectedResponse := `{"message": "invalid gift-card status"}`

	response, statusCode, err := makeGetReceivedGiftCardsRequest(suite.Token, 5)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestGetSentGiftCards_AcceptedStatus_Success() {
	require := suite.Require()
	expectedResponse := `{"gift_cards":[{"id":1,"amount":100,"status":0,"gifter_id":1,"giftee_id":1}, {"id":2,"amount":100,"status":0,"gifter_id":1,"giftee_id":1}],"total":2,"page":1}`

	response, statusCode, err := makeGetSentGiftCardsRequest(suite.Token, int(domain.GCSAccepted))

	require.NoError(err)
	require.Equal(http.StatusOK, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestGetSentGiftCards_RejectedStatus_Success() {
	require := suite.Require()
	expectedResponse := `{"gift_cards":[{"id":3,"amount":100,"status":1,"gifter_id":1,"giftee_id":1}, {"id":4,"amount":100,"status":1,"gifter_id":1,"giftee_id":1}],"total":2,"page":1}`

	response, statusCode, err := makeGetSentGiftCardsRequest(suite.Token, int(domain.GCSRejected))

	require.NoError(err)
	require.Equal(http.StatusOK, statusCode)
	require.JSONEq(expectedResponse, response)
}

func (suite *GiftCardsIntegrationTestSuite) TestGetSentGiftCards_InvalidGiftCardStatus_Failure() {
	require := suite.Require()
	expectedResponse := `{"message": "invalid gift-card status"}`

	response, statusCode, err := makeGetSentGiftCardsRequest(suite.Token, 5)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, statusCode)
	require.JSONEq(expectedResponse, response)
}

func TestCreateGiftCard(t *testing.T) {
	suite.Run(t, new(GiftCardsIntegrationTestSuite))
}
