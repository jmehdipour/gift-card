package service

import (
	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/repository"
)

type GiftCardService interface {
	CreateGiftCard(amount float64, gifterID, gifteeID uint) (*domain.GiftCard, error)
	FindGiftCard(id uint) (*domain.GiftCard, error)
	UpdateStatus(giftCardID uint, status domain.GiftCardStatus) error
	GetReceivedGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error)
	GetSentGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error)
}

type giftCardService struct {
	giftCardRepository repository.GiftCardRepository
}

func NewGiftCardService(giftCardRepo repository.GiftCardRepository) GiftCardService {
	return &giftCardService{giftCardRepository: giftCardRepo}
}

func (s *giftCardService) CreateGiftCard(amount float64, gifterID, gifteeID uint) (*domain.GiftCard, error) {
	giftCard := domain.GiftCard{
		Amount:   amount,
		GifterID: gifterID,
		GifteeID: gifteeID,
	}
	err := s.giftCardRepository.Create(&giftCard)
	if err != nil {
		return nil, err
	}

	return &giftCard, nil
}

func (s *giftCardService) FindGiftCard(id uint) (*domain.GiftCard, error) {
	return s.giftCardRepository.FindByID(id)
}

func (s *giftCardService) UpdateStatus(giftCardID uint, status domain.GiftCardStatus) error {
	return s.giftCardRepository.UpdateStatus(giftCardID, status)
}

func (s *giftCardService) GetReceivedGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	return s.giftCardRepository.FindReceivedGiftCardsByUserID(userID, status, pageSize, pageNumber)
}

func (s *giftCardService) GetSentGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	return s.giftCardRepository.FindSentGiftCardsByUserID(userID, status, pageSize, pageNumber)
}
