package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/repository"
)

type GiftCardServiceTestSuite struct {
	suite.Suite
	giftCardRepo    *repository.GiftCardRepositoryMock
	giftCardService *giftCardService
}

func (suite *GiftCardServiceTestSuite) SetupTest() {
	suite.giftCardRepo = new(repository.GiftCardRepositoryMock)
	suite.giftCardService = &giftCardService{
		giftCardRepository: suite.giftCardRepo,
	}
}

func (suite *GiftCardServiceTestSuite) TestNewGiftCardRepository() {
	require := suite.Require()

	service := NewGiftCardService(suite.giftCardRepo)

	require.NotNil(service)
}

func (suite *GiftCardServiceTestSuite) TestCreateGiftCard_Success() {
	require := suite.Require()
	giftCard := domain.GiftCard{
		ID:       15,
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
	}

	defer suite.giftCardRepo.On("Create", mock.Anything).Return(nil).Unset()
	giftCardResult, err := suite.giftCardService.CreateGiftCard(giftCard.Amount, giftCard.GifterID, giftCard.GifteeID)

	require.NoError(err)
	require.Equal(giftCard.ID, giftCardResult.ID)
}

func (suite *GiftCardServiceTestSuite) TestCreateGiftCard_Failure() {
	require := suite.Require()
	expectedError := errors.New("repo error")
	giftCard := domain.GiftCard{
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
	}

	defer suite.giftCardRepo.On("Create", mock.Anything).Return(expectedError).Unset()
	giftCardResult, err := suite.giftCardService.CreateGiftCard(giftCard.Amount, giftCard.GifterID, giftCard.GifteeID)

	require.Error(err)
	require.EqualError(expectedError, err.Error())
	require.Empty(giftCardResult)
}

func (suite *GiftCardServiceTestSuite) TestFindGiftCard_Failure() {
	require := suite.Require()
	expectedError := errors.New("repo error")
	id := uint(10)

	defer suite.giftCardRepo.On("FindByID", mock.Anything).Return(nil, expectedError).Unset()
	giftCardResult, err := suite.giftCardService.FindGiftCard(id)

	require.Error(err)
	require.EqualError(expectedError, err.Error())
	require.Empty(giftCardResult)
}

func (suite *GiftCardServiceTestSuite) TestFindGiftCard_Success() {
	require := suite.Require()
	giftCard := domain.GiftCard{
		ID:       10,
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
	}

	defer suite.giftCardRepo.On("FindByID", mock.Anything).Return(&giftCard, nil).Unset()
	giftCardResult, err := suite.giftCardService.FindGiftCard(giftCard.ID)

	require.NoError(err)
	require.Equal(&giftCard, giftCardResult)
}

func (suite *GiftCardServiceTestSuite) TestUpdateStatus_Success() {
	require := suite.Require()
	id := uint(10)

	defer suite.giftCardRepo.On("UpdateStatus", id, domain.GCSAccepted).Return(nil).Unset()
	err := suite.giftCardService.UpdateStatus(id, domain.GCSAccepted)

	require.NoError(err)
}

func (suite *GiftCardServiceTestSuite) TestUpdateStatus_Failure() {
	require := suite.Require()
	expectedError := errors.New("repo error")
	id := uint(10)

	defer suite.giftCardRepo.On("UpdateStatus", id, domain.GCSAccepted).Return(expectedError).Unset()
	err := suite.giftCardService.UpdateStatus(id, domain.GCSAccepted)

	require.Error(err)
}

func (suite *GiftCardServiceTestSuite) TestFindReceivedGiftCardsByUserID_Failure() {
	require := suite.Require()
	expectedError := errors.New("repo error")
	userID := uint(10)

	defer suite.giftCardRepo.On("FindReceivedGiftCardsByUserID", userID, (*domain.GiftCardStatus)(nil), 10, 1).
		Return([]domain.GiftCard{}, 0, expectedError).Unset()
	giftCards, total, err := suite.giftCardService.GetReceivedGiftCardsByUserID(userID, nil, 10, 1)

	require.Error(err)
	require.Empty(giftCards)
	require.Zero(total)
}

func (suite *GiftCardServiceTestSuite) TestFindReceivedGiftCardsByUserID_Success() {
	require := suite.Require()
	userID := uint(20)

	giftCards := []domain.GiftCard{{ID: 10, Amount: 100, Status: domain.GCSAccepted, GifterID: 10, GifteeID: 20}}
	defer suite.giftCardRepo.On("FindReceivedGiftCardsByUserID", userID, (*domain.GiftCardStatus)(nil), 10, 1).
		Return(giftCards, 1, nil).Unset()
	giftCardsResult, total, err := suite.giftCardService.GetReceivedGiftCardsByUserID(userID, nil, 10, 1)

	require.NoError(err)
	require.Equal(giftCards, giftCardsResult)
	require.Equal(len(giftCards), total)
}

func (suite *GiftCardServiceTestSuite) TestFindSentGiftCardsByUserID_Failure() {
	require := suite.Require()
	expectedError := errors.New("repo error")
	userID := uint(10)

	defer suite.giftCardRepo.On("FindSentGiftCardsByUserID", userID, (*domain.GiftCardStatus)(nil), 10, 1).
		Return([]domain.GiftCard{}, 0, expectedError).Unset()
	giftCards, total, err := suite.giftCardService.GetSentGiftCardsByUserID(userID, nil, 10, 1)

	require.Error(err)
	require.Empty(giftCards)
	require.Zero(total)
}

func (suite *GiftCardServiceTestSuite) TestFindSendGiftCardsByUserID_Success() {
	require := suite.Require()
	userID := uint(10)

	giftCards := []domain.GiftCard{{ID: 10, Amount: 100, Status: domain.GCSAccepted, GifterID: 10, GifteeID: 20}}
	defer suite.giftCardRepo.On("FindSentGiftCardsByUserID", userID, (*domain.GiftCardStatus)(nil), 10, 1).
		Return(giftCards, 1, nil).Unset()
	giftCardsResult, total, err := suite.giftCardService.GetSentGiftCardsByUserID(userID, nil, 10, 1)

	require.NoError(err)
	require.Equal(giftCards, giftCardsResult)
	require.Equal(len(giftCards), total)
}

func TestGiftCardService(t *testing.T) {
	suite.Run(t, new(GiftCardServiceTestSuite))
}
