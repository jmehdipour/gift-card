package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func CreateDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Infof("Successfully connected  to the database!")

	return db, nil
}
