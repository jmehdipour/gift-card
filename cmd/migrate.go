package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jmehdipour/gift-card/internal/config"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/database"
)

var migrateDatabaseCMD = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrateGiftCardDB()
	},
}

func migrateGiftCardDB() {
	db, err := database.CreateDatabase(config.C.Database.String())
	if err != nil {
		log.Fatalf("Cannot open database: %s", err)
	}

	drop := `DROP TABLE IF EXISTS gift_cards;
DROP TABLE IF EXISTS users;`

	_, err = db.Exec(drop)
	if err != nil {
		log.Fatal("database migration (drop tables) failed: ", err)
	}

	schema := `CREATE TABLE gift_cards (
    id INT AUTO_INCREMENT,
    amount DECIMAL(10, 2),
    sender_id INT,
    receiver_id INT,
    created_at DATETIME,
    updated_at DATETIME,
    status TINYINT,
    PRIMARY KEY (id)
);
CREATE TABLE users (
    id INT AUTO_INCREMENT,
    username VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id),
    UNIQUE (username),
    UNIQUE (email)
);
`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal("database migration (create tables) failed: ", err)
	}
}
