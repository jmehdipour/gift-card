package domain

import (
	"time"
)

type GiftCardStatus int

var validGiftCard = map[GiftCardStatus]struct{}{
	GCSAccepted: {},
	GCSRejected: {},
	GCSPending:  {},
}

func (s GiftCardStatus) IsValid() bool {
	_, ok := validGiftCard[s]

	return ok
}

const (
	GCSAccepted GiftCardStatus = iota
	GCSRejected
	GCSPending
)

type GiftCard struct {
	ID           uint
	Amount       float64
	Status       GiftCardStatus
	GifterID     uint
	GifteeID     uint
	CreationDate time.Time
}

func (c *GiftCard) CanUpdateStatus() bool {
	if c.Status == GCSPending {
		return true
	}

	return false
}
