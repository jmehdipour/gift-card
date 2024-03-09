package repository

import (
	"database/sql"
	"time"

	"github.com/jmehdipour/gift-card/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type UserEntity struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u UserEntity) ToAggregate() domain.User {
	return domain.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users(email, password, created_at, updated_at) VALUES(?, ?, NOW(), NOW())`
	result, err := r.db.Exec(query, user.Email, user.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)

	return nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var e UserEntity
	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).
		Scan(&e.ID, &e.Email, &e.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	domainUser := e.ToAggregate()

	return &domainUser, nil
}
