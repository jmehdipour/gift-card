package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashed)

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}

	return true
}
