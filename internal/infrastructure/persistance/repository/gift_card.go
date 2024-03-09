package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmehdipour/gift-card/internal/domain"
)

type GiftCardRepository interface {
	Create(giftCard *domain.GiftCard) error
	FindByID(id uint) (*domain.GiftCard, error)
	UpdateStatus(id uint, status domain.GiftCardStatus) error
	FindReceivedGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error)
	FindSentGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error)
}

type GiftCardEntity struct {
	ID         uint
	Amount     float64
	SenderID   uint
	ReceiverID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     int
}

func (g GiftCardEntity) ToAggregate() domain.GiftCard {
	return domain.GiftCard{
		ID:           g.ID,
		Status:       domain.GiftCardStatus(g.Status),
		GifterID:     g.SenderID,
		GifteeID:     g.ReceiverID,
		Amount:       g.Amount,
		CreationDate: g.CreatedAt,
	}
}

type giftCardRepository struct {
	db *sql.DB
}

func NewGiftCardRepository(db *sql.DB) GiftCardRepository {
	return &giftCardRepository{db: db}
}

func (r *giftCardRepository) Create(giftCard *domain.GiftCard) error {
	query := `INSERT INTO gift_cards (amount, sender_id, receiver_id, status, updated_at, created_at) VALUES (?, ?, ?, 2, NOW(), NOW())`
	res, err := r.db.Exec(query, giftCard.Amount, giftCard.GifterID, giftCard.GifteeID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	giftCard.ID = uint(id)

	return nil
}

func (r *giftCardRepository) FindByID(id uint) (*domain.GiftCard, error) {
	e := new(GiftCardEntity)
	err := r.db.
		QueryRow("SELECT id, sender_id, receiver_id, amount, status FROM gift_cards WHERE id = ?", id).
		Scan(&e.ID, &e.SenderID, &e.ReceiverID, &e.Amount, &e.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	domainGiftCard := e.ToAggregate()

	return &domainGiftCard, nil
}

func (r *giftCardRepository) UpdateStatus(id uint, status domain.GiftCardStatus) error {
	query := "UPDATE gift_cards SET status = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.Exec(query, uint(status), id)

	return err
}

func (r *giftCardRepository) FindReceivedGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	offset := (pageNumber - 1) * pageSize
	query := "SELECT id, status, sender_id, receiver_id, amount FROM gift_cards WHERE receiver_id = ?"
	if status != nil {
		query += fmt.Sprintf(" AND status = %d", *status)
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var giftCards []domain.GiftCard
	for rows.Next() {
		var g GiftCardEntity
		err := rows.Scan(&g.ID, &g.Status, &g.SenderID, &g.ReceiverID, &g.Amount)
		if err != nil {
			return nil, 0, err
		}

		giftCards = append(giftCards, g.ToAggregate())
	}

	// Fetch total count
	var totalCount int
	totalCountQuery := "SELECT COUNT(*) FROM gift_cards WHERE receiver_id = ?"
	if status != nil {
		totalCountQuery += fmt.Sprintf(" AND status = %d", *status)
	}

	err = r.db.QueryRow(totalCountQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return giftCards, totalCount, nil
}

func (r *giftCardRepository) FindSentGiftCardsByUserID(userID uint, status *domain.GiftCardStatus, pageSize int, pageNumber int) ([]domain.GiftCard, int, error) {
	offset := (pageNumber - 1) * pageSize
	query := "SELECT id, status, sender_id, receiver_id, amount FROM gift_cards WHERE sender_id = ?"
	if status != nil {
		query += fmt.Sprintf(" AND status = %d", *status)
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var giftCards []domain.GiftCard
	for rows.Next() {
		var g GiftCardEntity
		err := rows.Scan(&g.ID, &g.Status, &g.SenderID, &g.ReceiverID, &g.Amount)
		if err != nil {
			return nil, 0, err
		}

		giftCards = append(giftCards, g.ToAggregate())
	}

	// Fetch total count
	var totalCount int
	totalCountQuery := "SELECT COUNT(*) FROM gift_cards WHERE sender_id = ?"
	if status != nil {
		totalCountQuery += fmt.Sprintf(" AND status = %d", *status)
	}
	err = r.db.QueryRow(totalCountQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return giftCards, totalCount, nil
}
