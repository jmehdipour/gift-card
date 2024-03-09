package service

import (
	"github.com/jmehdipour/gift-card/internal/domain"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/repository"
)

type UserService interface {
	CreateUser(email, password string) (*domain.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepository: userRepo}
}

func (s *userService) CreateUser(email, password string) (*domain.User, error) {
	user := &domain.User{Email: email}
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
