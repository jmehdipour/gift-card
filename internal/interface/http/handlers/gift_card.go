package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/service"
)

type CreateGiftCardRequest struct {
	Amount   float64 `json:"amount"`
	GifteeID uint    `json:"giftee_id"`
}

type GiftCardResponse struct {
	ID       uint    `json:"id"`
	Amount   float64 `json:"amount"`
	Status   int     `json:"status"`
	GifterID uint    `json:"gifter_id"`
	GifteeID uint    `json:"giftee_id"`
}

func CreateGiftCardHandler(giftCardService service.GiftCardService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := new(CreateGiftCardRequest)
		err := ctx.Bind(request)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		}

		userID := ctx.Get("user_id").(uint)
		giftCard, err := giftCardService.CreateGiftCard(request.Amount, userID, request.GifteeID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		}

		return ctx.JSON(http.StatusCreated, GiftCardResponse{
			ID:       giftCard.ID,
			Amount:   giftCard.Amount,
			GifterID: giftCard.GifterID,
			GifteeID: giftCard.GifteeID,
			Status:   int(domain.GCSPending),
		})
	}
}

type UpdateGiftCardStatusRequest struct {
	Status int `json:"status"`
}

func UpdateGiftCardStatusHandler(giftCardService service.GiftCardService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := new(UpdateGiftCardStatusRequest)
		err := ctx.Bind(request)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: "Invalid request body"})
		}

		giftCardID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: "Invalid gift card ID"})
		}

		giftCard, err := giftCardService.FindGiftCard(uint(giftCardID))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, MessageResponse{Message: "Failed to get gift card"})
		}

		if giftCard == nil {
			return ctx.JSON(http.StatusNotFound, MessageResponse{Message: "Gift card not found"})
		}

		userID := ctx.Get("user_id").(uint)
		if giftCard.GifteeID != userID {
			return ctx.JSON(http.StatusForbidden, MessageResponse{Message: fmt.Sprintf("forbidden: user %d is not the receiver of gift card %d", userID, giftCardID)})
		}

		status := domain.GiftCardStatus(request.Status)
		if !status.IsValid() {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: "Invalid gift card status for update"})
		}

		err = giftCardService.UpdateStatus(uint(giftCardID), status)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, MessageResponse{Message: "Failed to update gift card status"})
		}

		return ctx.NoContent(http.StatusOK)
	}
}

type GetGiftCards struct {
	GiftCards []GiftCardResponse `json:"gift_cards"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
}

func normalizeStatus(statusStr string) (*domain.GiftCardStatus, error) {
	requestedStatusInt, err := strconv.Atoi(statusStr)
	if err != nil {
		return nil, err
	}

	status := domain.GiftCardStatus(requestedStatusInt)
	if !status.IsValid() {
		return nil, errors.New("invalid gift-card status")
	}

	return &status, nil
}

func GetReceivedGiftCardsHandler(giftCardService service.GiftCardService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(uint)
		pageSize := 10
		status, err := normalizeStatus(ctx.QueryParam("status"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid gift-card status"})
		}

		pageNumberStr := ctx.QueryParam("page")
		pageNumberInt, _ := strconv.Atoi(pageNumberStr)
		if pageNumberInt < 1 {
			pageNumberInt = 1
		}

		giftCards, totalCount, err := giftCardService.GetReceivedGiftCardsByUserID(userID, status, pageSize, pageNumberInt)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, MessageResponse{Message: "Failed to get gift cards"})
		}

		giftCardsResponse := make([]GiftCardResponse, 0, len(giftCards))
		for _, g := range giftCards {
			giftCardsResponse = append(giftCardsResponse, GiftCardResponse{
				ID:       g.ID,
				GifteeID: g.GifteeID,
				GifterID: g.GifterID,
				Amount:   g.Amount,
				Status:   int(g.Status),
			})
		}

		return ctx.JSON(http.StatusOK, GetGiftCards{GiftCards: giftCardsResponse, Total: totalCount, Page: pageNumberInt})
	}
}

func GetSentGiftCardsHandler(giftCardService service.GiftCardService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(uint)
		pageSize := 10
		status, err := normalizeStatus(ctx.QueryParam("status"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid gift-card status"})
		}

		pageNumberStr := ctx.QueryParam("page")
		pageNumberInt, _ := strconv.Atoi(pageNumberStr)
		if pageNumberInt < 1 {
			pageNumberInt = 1
		}

		giftCards, totalCount, err := giftCardService.GetSentGiftCardsByUserID(userID, status, pageSize, pageNumberInt)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, MessageResponse{Message: "Failed to get gift cards"})
		}

		giftCardsResponse := make([]GiftCardResponse, 0, len(giftCards))
		for _, g := range giftCards {
			giftCardsResponse = append(giftCardsResponse, GiftCardResponse{
				ID:       g.ID,
				GifteeID: g.GifteeID,
				GifterID: g.GifterID,
				Amount:   g.Amount,
				Status:   int(g.Status),
			})
		}

		return ctx.JSON(http.StatusOK, GetGiftCards{GiftCards: giftCardsResponse, Total: totalCount, Page: pageNumberInt})
	}
}
