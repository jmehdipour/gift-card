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

func createGiftCardNewEchoContext(body string, userID uint) (echo.Context, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(http.MethodPost, "/gift-cards", bytes.NewReader([]byte(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(request, response)
	ctx.Set("user_id", userID)

	return ctx, response
}

func updateGiftCardNewEchoContext(body string, userID uint, giftCardID uint) (echo.Context, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/gift-cards/%d/status", giftCardID),
		bytes.NewReader([]byte(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(request, response)
	ctx.Set("user_id", userID)
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprintf("%d", giftCardID))

	return ctx, response
}

func getReceivedGiftCardsNewEchoContext(userID uint, status int) (echo.Context, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/gift-cards/received?status=%d", status),
		nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(request, response)
	ctx.Set("user_id", userID)

	return ctx, response
}

func getSentGiftCardsNewEchoContext(userID uint, status int) (echo.Context, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/gift-cards/sent?status=%d", status),
		nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(request, response)
	ctx.Set("user_id", userID)

	return ctx, response
}

type CreateGiftCardsHandlerTestSuite struct {
	suite.Suite
	giftCardService *service.GiftCardServiceMock
}

func (suite *CreateGiftCardsHandlerTestSuite) SetupSuite() {
	suite.giftCardService = new(service.GiftCardServiceMock)
}

func (suite *CreateGiftCardsHandlerTestSuite) TestCreateGiftCardHandler_Success() {
	require := suite.Require()
	userID := uint(10)
	giftCard := domain.GiftCard{
		ID:       15,
		Amount:   100,
		GifterID: userID,
		GifteeID: 20,
	}
	requestBody := fmt.Sprintf(`{"amount": %v, "giftee_id": %d}`, giftCard.Amount, giftCard.GifteeID)
	expectedResponse := `{"id":15,"amount":100,"status":2,"gifter_id":10,"giftee_id":20}`

	defer suite.giftCardService.On("CreateGiftCard", giftCard.Amount, userID, giftCard.GifteeID).Return(&giftCard, nil).Unset()

	ctx, response := createGiftCardNewEchoContext(requestBody, userID)
	err := CreateGiftCardHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.JSONEq(expectedResponse, response.Body.String())
	require.Equal(http.StatusCreated, response.Code)
}

func (suite *CreateGiftCardsHandlerTestSuite) TestCreateGiftCardHandler_InvalidRequestBody_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCard := domain.GiftCard{
		ID:       15,
		Amount:   100,
		GifterID: userID,
		GifteeID: 20,
	}
	requestBody := fmt.Sprintf(`{"amount":, "giftee_id": %d}`, giftCard.GifteeID)
	expectedResponse := `{"message":"code=400, message=Syntax error: offset=11, error=invalid character ',' looking for beginning of value, internal=invalid character ',' looking for beginning of value"}`

	ctx, response := createGiftCardNewEchoContext(requestBody, userID)
	err := CreateGiftCardHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.JSONEq(expectedResponse, response.Body.String())
	require.Equal(http.StatusBadRequest, response.Code)
}

func (suite *CreateGiftCardsHandlerTestSuite) TestCreateGiftCardHandler_ServiceError_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCard := domain.GiftCard{
		Amount:   100,
		GifterID: userID,
		GifteeID: 20,
	}
	requestBody := fmt.Sprintf(`{"amount": %v, "giftee_id": %d}`, giftCard.Amount, giftCard.GifteeID)
	expectedResponse := `{"message":"service layer error"}`
	expectedError := errors.New("service layer error")

	defer suite.giftCardService.On("CreateGiftCard", giftCard.Amount, userID, giftCard.GifteeID).Return(nil, expectedError).Unset()

	ctx, response := createGiftCardNewEchoContext(requestBody, userID)
	err := CreateGiftCardHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.JSONEq(expectedResponse, response.Body.String())
	require.Equal(http.StatusInternalServerError, response.Code)
}

type UpdateGiftCardStatusHandlerTestSuite struct {
	suite.Suite
	giftCardService *service.GiftCardServiceMock
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) SetupSuite() {
	suite.giftCardService = new(service.GiftCardServiceMock)
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) TestUpdateGiftCardHandler_Success() {
	require := suite.Require()
	userID := uint(10)
	giftCardID := uint(101)
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	giftCard := domain.GiftCard{
		Amount:   100,
		GifterID: 20,
		GifteeID: userID,
	}

	defer suite.giftCardService.On("FindGiftCard", giftCardID).Return(&giftCard, nil).Unset()
	defer suite.giftCardService.On("UpdateStatus", giftCardID, status).Return(nil).Unset()

	ctx, response := updateGiftCardNewEchoContext(requestBody, userID, giftCardID)
	err := UpdateGiftCardStatusHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusOK, response.Code)
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) TestUpdateGiftCardHandler_InvalidRequestBody_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCardID := uint(101)
	requestBody := `{"status": "foo"}`
	expectedResponse := `{"message":"Invalid request body"}`

	ctx, response := updateGiftCardNewEchoContext(requestBody, userID, giftCardID)
	err := UpdateGiftCardStatusHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) TestUpdateGiftCardHandler_FindGiftCardError_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCardID := uint(101)
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	expectedResponse := `{"message": "Failed to get gift card"}`

	defer suite.giftCardService.On("FindGiftCard", giftCardID).Return(nil, errors.New("service error")).Unset()

	ctx, response := updateGiftCardNewEchoContext(requestBody, userID, giftCardID)
	err := UpdateGiftCardStatusHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusInternalServerError, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) TestUpdateGiftCardHandler_GiftCardNotFound_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCardID := uint(101)
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	expectedResponse := `{"message": "Gift card not found"}`

	defer suite.giftCardService.On("FindGiftCard", giftCardID).Return(nil, nil).Unset()

	ctx, response := updateGiftCardNewEchoContext(requestBody, userID, giftCardID)
	err := UpdateGiftCardStatusHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusNotFound, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) TestUpdateGiftCardHandler_InvalidGiftCard_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCardID := uint(101)
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	expectedResponse := `{"message": "Invalid gift card ID"}`

	defer suite.giftCardService.On("FindGiftCard", giftCardID).Return(nil, nil).Unset()

	ctx, response := updateGiftCardNewEchoContext(requestBody, userID, giftCardID)
	ctx.SetParamValues("foo")
	err := UpdateGiftCardStatusHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) TestUpdateGiftCardHandler_UnauthorizedUser_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCardID := uint(101)
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	expectedResponse := fmt.Sprintf(`{"message": "forbidden: user %d is not the receiver of gift card %d"}`, userID, giftCardID)
	giftCard := domain.GiftCard{
		ID:       giftCardID,
		Amount:   100,
		GifterID: 20,
		GifteeID: 30,
	}

	defer suite.giftCardService.On("FindGiftCard", giftCardID).Return(&giftCard, nil).Unset()

	ctx, response := updateGiftCardNewEchoContext(requestBody, userID, giftCardID)
	err := UpdateGiftCardStatusHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusForbidden, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *UpdateGiftCardStatusHandlerTestSuite) TestUpdateGiftCardHandler_UpdateStatusError_Failure() {
	require := suite.Require()
	userID := uint(10)
	giftCardID := uint(101)
	status := domain.GCSRejected
	requestBody := fmt.Sprintf(`{"status": %v}`, status)
	giftCard := domain.GiftCard{
		ID:       giftCardID,
		Amount:   100,
		GifterID: 20,
		GifteeID: userID,
	}

	defer suite.giftCardService.On("FindGiftCard", giftCardID).Return(&giftCard, nil).Unset()
	defer suite.giftCardService.On("UpdateStatus", giftCardID, status).Return(errors.New("update error")).Unset()

	ctx, response := updateGiftCardNewEchoContext(requestBody, userID, giftCardID)
	err := UpdateGiftCardStatusHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusInternalServerError, response.Code)
}

type GetReceivedGiftCardsHandlerTestSuite struct {
	suite.Suite
	giftCardService *service.GiftCardServiceMock
}

func (suite *GetReceivedGiftCardsHandlerTestSuite) SetupSuite() {
	suite.giftCardService = new(service.GiftCardServiceMock)
}

func (suite *GetReceivedGiftCardsHandlerTestSuite) TestGetReceivedGiftCardsHandler_Success() {
	require := suite.Require()
	userID := uint(10)
	status := domain.GCSAccepted
	giftCards := []domain.GiftCard{{ID: 10, Amount: 100, Status: status, GifterID: 10, GifteeID: userID}}
	expectedResponse := `{"gift_cards":[{"id":10,"amount":100,"status":0,"gifter_id":10,"giftee_id":10}],"total":1,"page":1}`

	defer suite.giftCardService.On("GetReceivedGiftCardsByUserID", userID, &status, 10, 1).Return(giftCards, len(giftCards), nil).Unset()

	ctx, response := getReceivedGiftCardsNewEchoContext(userID, int(status))
	err := GetReceivedGiftCardsHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusOK, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *GetReceivedGiftCardsHandlerTestSuite) TestGetReceivedGiftCardsHandler_InvalidGiftCardStatus_Failure() {
	require := suite.Require()
	userID := uint(10)
	status := 5
	expectedResponse := `{"message": "invalid gift-card status"}`

	ctx, response := getReceivedGiftCardsNewEchoContext(userID, status)
	err := GetReceivedGiftCardsHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *GetReceivedGiftCardsHandlerTestSuite) TestGetReceivedGiftCardsHandler_ServiceError_Failure() {
	require := suite.Require()
	userID := uint(10)
	status := domain.GCSRejected
	expectedResponse := `{"message": "Failed to get gift cards"}`

	defer suite.giftCardService.On("GetReceivedGiftCardsByUserID", userID, &status, 10, 1).
		Return(nil, 0, errors.New("service layer error")).Unset()

	ctx, response := getReceivedGiftCardsNewEchoContext(userID, int(status))
	err := GetReceivedGiftCardsHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusInternalServerError, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

type GetGetGiftCardsHandlerTestSuite struct {
	suite.Suite
	giftCardService *service.GiftCardServiceMock
}

func (suite *GetGetGiftCardsHandlerTestSuite) SetupSuite() {
	suite.giftCardService = new(service.GiftCardServiceMock)
}

func (suite *GetGetGiftCardsHandlerTestSuite) TestGetSentGiftCardsHandler_Success() {
	require := suite.Require()
	userID := uint(10)
	status := domain.GCSAccepted
	giftCards := []domain.GiftCard{{ID: 10, Amount: 100, Status: status, GifterID: 10, GifteeID: userID}}
	expectedResponse := `{"gift_cards":[{"id":10,"amount":100,"status":0,"gifter_id":10,"giftee_id":10}],"total":1,"page":1}`

	defer suite.giftCardService.On("GetSentGiftCardsByUserID", userID, &status, 10, 1).Return(giftCards, len(giftCards), nil).Unset()

	ctx, response := getSentGiftCardsNewEchoContext(userID, int(status))
	err := GetSentGiftCardsHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusOK, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *GetGetGiftCardsHandlerTestSuite) TestGetSentGiftCardsHandler_InvalidGiftCardStatus_Failure() {
	require := suite.Require()
	userID := uint(10)
	status := 5
	expectedResponse := `{"message": "invalid gift-card status"}`

	ctx, response := getSentGiftCardsNewEchoContext(userID, status)
	err := GetSentGiftCardsHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusBadRequest, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func (suite *GetGetGiftCardsHandlerTestSuite) TestGetSentGiftCardsHandler_ServiceError_Failure() {
	require := suite.Require()
	userID := uint(10)
	status := domain.GCSRejected
	expectedResponse := `{"message": "Failed to get gift cards"}`

	defer suite.giftCardService.On("GetSentGiftCardsByUserID", userID, &status, 10, 1).
		Return(nil, 0, errors.New("service layer error")).Unset()

	ctx, response := getSentGiftCardsNewEchoContext(userID, int(status))
	err := GetSentGiftCardsHandler(suite.giftCardService)(ctx)

	require.NoError(err)
	require.Equal(http.StatusInternalServerError, response.Code)
	require.JSONEq(expectedResponse, response.Body.String())
}

func TestCreateGiftCardHandler(t *testing.T) {
	suite.Run(t, new(CreateGiftCardsHandlerTestSuite))
}

func TestUpdateGiftCardStatusHandler(t *testing.T) {
	suite.Run(t, new(UpdateGiftCardStatusHandlerTestSuite))
}

func TestGetReceivedGiftCardsHandler(t *testing.T) {
	suite.Run(t, new(GetReceivedGiftCardsHandlerTestSuite))
}

func TestGetSentGiftCardsHandler(t *testing.T) {
	suite.Run(t, new(GetGetGiftCardsHandlerTestSuite))
}
