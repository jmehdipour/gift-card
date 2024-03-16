package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jmehdipour/gift-card/internal/config"
	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/database"
)

var seedDatabaseCMD = &cobra.Command{
	Use:   "seed",
	Short: "seed database with data",
	Run: func(cmd *cobra.Command, args []string) {
		seedMainDB()
	},
}

func seedMainDB() {
	db, err := database.CreateDatabase(config.C.Database.String())
	if err != nil {
		log.Fatalf("Cannot open database: %s", err)
	}

	for i := 0; i < 2; i++ {
		u := domain.User{Email: fmt.Sprintf("test%d@example.com", i)}
		_ = u.SetPassword("password")
		insertUserQuery := `INSERT INTO users(email, password, created_at, updated_at) VALUES(?, ?, NOW(), NOW())`
		_, err = db.Exec(insertUserQuery, u.Email, u.Password)
		if err != nil {
			log.Fatal("database seed (insert user) failed: ", err)
		}
	}

	for i := 0; i < 2; i++ {
		giftCard := domain.GiftCard{Amount: 100, GifterID: 1, GifteeID: 1}
		insertGiftCardQuery := `INSERT INTO gift_cards (amount, sender_id, receiver_id, status, updated_at, created_at) VALUES (?, ?, ?, 0, NOW(), NOW())`
		_, err = db.Exec(insertGiftCardQuery, giftCard.Amount, giftCard.GifterID, giftCard.GifteeID)
		if err != nil {
			log.Fatal("database seed (insert gift-card) failed: ", err)
		}
	}

	for i := 0; i < 2; i++ {
		giftCard := domain.GiftCard{Amount: 100, GifterID: 1, GifteeID: 1}
		insertGiftCardQuery := `INSERT INTO gift_cards (amount, sender_id, receiver_id, status, updated_at, created_at) VALUES (?, ?, ?, 1, NOW(), NOW())`
		_, err = db.Exec(insertGiftCardQuery, giftCard.Amount, giftCard.GifterID, giftCard.GifteeID)
		if err != nil {
			log.Fatal("database seed (insert gift-card) failed: ", err)
		}
	}

	log.Info("database seed was successful")
}
