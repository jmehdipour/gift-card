package repository

import (
	"github.com/stretchr/testify/mock"

	"github.com/jmehdipour/gift-card/internal/domain"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) Create(user *domain.User) error {
	args := u.Called(user)
	user.ID = 15

	return args.Error(0)
}

func (u *UserRepositoryMock) FindByEmail(email string) (*domain.User, error) {
	args := u.Called(email)

	var r0 *domain.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.User)
	}

	return r0, args.Error(1)
}

type GiftCardRepositoryMock struct {
	mock.Mock
}

func (r *GiftCardRepositoryMock) Create(giftCard *domain.GiftCard) error {
	args := r.Called(giftCard)
	giftCard.ID = 15

	return args.Error(0)
}

func (r *GiftCardRepositoryMock) FindByID(id uint) (*domain.GiftCard, error) {
	args := r.Called(id)

	var r0 *domain.GiftCard
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.GiftCard)
	}

	return r0, args.Error(1)
}

func (r *GiftCardRepositoryMock) UpdateStatus(id uint, status domain.GiftCardStatus) error {
	args := r.Called(id, status)

	return args.Error(0)
}

func (r *GiftCardRepositoryMock) FindReceivedGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	args := r.Called(userID, status, pageSize, pageNumber)

	var r0 []domain.GiftCard
	if args.Get(0) != nil {
		r0 = args.Get(0).([]domain.GiftCard)
	}

	return r0, args.Int(1), args.Error(2)
}

func (r *GiftCardRepositoryMock) FindSentGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	args := r.Called(userID, status, pageSize, pageNumber)

	var r0 []domain.GiftCard
	if args.Get(0) != nil {
		r0 = args.Get(0).([]domain.GiftCard)
	}

	return r0, args.Int(1), args.Error(2)
}
