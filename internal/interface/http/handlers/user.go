package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/jmehdipour/gift-card/internal/service"
	"github.com/jmehdipour/gift-card/utils"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func (r CreateUserRequest) Validate() error {
	if !utils.ValidateEmail(r.Email) {
		return errors.New("invalid email")
	}

	if strings.TrimSpace(r.Password) == "" {
		return errors.New("invalid password")
	}

	return nil
}

func CreateUserHandler(userService service.UserService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := new(CreateUserRequest)
		err := ctx.Bind(request)
		if err != nil {
			return err
		}

		err = request.Validate()
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		}

		user, err := userService.CreateUser(request.Email, request.Password)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		}

		return ctx.JSON(http.StatusCreated, CreateUserResponse{ID: user.ID, Email: user.Email})
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginHandlerResponse struct {
	Token string `json:"token"`
}

func LoginHandler(authService service.AuthService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := new(LoginRequest)
		err := ctx.Bind(request)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid request body"})
		}

		token, err := authService.Login(request.Email, request.Password)
		if err != nil || token == "" {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		return ctx.JSON(http.StatusOK, LoginHandlerResponse{Token: token})
	}
}
