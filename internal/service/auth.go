package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jmehdipour/gift-card/internal/config"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/repository"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepository: userRepo}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil || !user.CheckPassword(password) {
		return "", nil
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	signedToken, err := token.SignedString([]byte(config.C.User.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
