package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/jmehdipour/gift-card/internal/config"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/database"
	"github.com/jmehdipour/gift-card/internal/infrastructure/persistance/repository"
	"github.com/jmehdipour/gift-card/internal/interface/http/handlers"
	"github.com/jmehdipour/gift-card/internal/interface/http/middleware"
	"github.com/jmehdipour/gift-card/internal/service"
)

// Package http is a package for defining and creating a http server
//
// Example of usage:
//
//	http.NewServer().Serve()
//
// Description of what package do:
// This package creates a http server and defines its routes.
// It also handles the server graceful shutdown.

var asciiArt = ` 
________  ___  ________ _________        ________  ________  ________  ________     
|\   ____\|\  \|\  _____\\___   ___\     |\   ____\|\   __  \|\   __  \|\   ___ \    
\ \  \___|\ \  \ \  \__/\|___ \  \_|     \ \  \___|\ \  \|\  \ \  \|\  \ \  \_|\ \   
 \ \  \  __\ \  \ \   __\    \ \  \       \ \  \    \ \   __  \ \   _  _\ \  \ \\ \  
  \ \  \|\  \ \  \ \  \_|     \ \  \       \ \  \____\ \  \ \  \ \  \\  \\ \  \_\\ \ 
   \ \_______\ \__\ \__\       \ \__\       \ \_______\ \__\ \__\ \__\\ _\\ \_______\
    \|_______|\|__|\|__|        \|__|        \|_______|\|__|\|__|\|__|\|__|\|_______|`

type Server interface {
	Serve()
}

// echoServer is the struct that holds the echo server
type echoServer struct {
	e *echo.Echo
}

// NewServer creates a new echo echoServer
func NewServer() Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(echomw.Logger())

	return &echoServer{
		e: e,
	}
}

// Serve starts the echo server and listens on the configured port
func (s *echoServer) Serve() {
	ctx := context.Background()
	db, err := database.CreateDatabase(config.C.Database.String())
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	giftCardRepo := repository.NewGiftCardRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	giftCardService := service.NewGiftCardService(giftCardRepo)

	s.e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, asciiArt)
	})

	s.e.POST("/users/register", handlers.CreateUserHandler(userService))
	s.e.POST("/users/login", handlers.LoginHandler(authService))

	s.e.POST("/gift-cards", handlers.CreateGiftCardHandler(giftCardService), middleware.ValidateUser())
	s.e.PUT("/gift-cards/:id/status", handlers.UpdateGiftCardStatusHandler(giftCardService), middleware.ValidateUser())
	s.e.GET("/gift-cards/received", handlers.GetReceivedGiftCardsHandler(giftCardService), middleware.ValidateUser())
	s.e.GET("/gift-cards/sent", handlers.GetSentGiftCardsHandler(giftCardService), middleware.ValidateUser())

	go func() {
		if err := s.e.Start(config.C.HTTPServer.Address); err != nil && err != http.ErrServerClosed {
			s.e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGTERM,
		syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		log.Fatalf("error in shutdown: %v", err)
	}
}
