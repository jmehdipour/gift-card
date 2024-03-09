package service

import (
	"github.com/stretchr/testify/mock"

	"github.com/jmehdipour/gift-card/internal/domain"
)

type UserServiceMock struct {
	mock.Mock
}

func (s *UserServiceMock) CreateUser(email, password string) (*domain.User, error) {
	args := s.Called(email, password)

	var r0 *domain.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.User)
	}

	return r0, args.Error(1)
}

type GiftCardServiceMock struct {
	mock.Mock
}

func (s *GiftCardServiceMock) CreateGiftCard(amount float64, gifterID, gifteeID uint) (*domain.GiftCard, error) {
	args := s.Called(amount, gifterID, gifteeID)

	var r0 *domain.GiftCard
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.GiftCard)
	}

	return r0, args.Error(1)
}

func (s *GiftCardServiceMock) FindGiftCard(id uint) (*domain.GiftCard, error) {
	args := s.Called(id)

	var r0 *domain.GiftCard
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.GiftCard)
	}

	return r0, args.Error(1)
}

func (s *GiftCardServiceMock) UpdateStatus(giftCardID uint, status domain.GiftCardStatus) error {
	args := s.Called(giftCardID, status)

	return args.Error(0)
}

func (s *GiftCardServiceMock) GetReceivedGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	args := s.Called(userID, status, pageSize, pageNumber)

	var r0 []domain.GiftCard
	if args.Get(0) != nil {
		r0 = args.Get(0).([]domain.GiftCard)
	}

	return r0, args.Int(1), args.Error(2)
}

func (s *GiftCardServiceMock) GetSentGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	args := s.Called(userID, status, pageSize, pageNumber)

	var r0 []domain.GiftCard
	if args.Get(0) != nil {
		r0 = args.Get(0).([]domain.GiftCard)
	}

	return r0, args.Int(1), args.Error(2)
}

type AuthServiceMock struct {
	mock.Mock
}

func (s *AuthServiceMock) Login(email, password string) (string, error) {
	args := s.Called(email, password)

	return args.String(0), args.Error(1)
}
