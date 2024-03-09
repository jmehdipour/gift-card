package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/jmehdipour/gift-card/internal/config"
)

func ValidateUser() echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			token := ctx.Request().Header.Get("Authorization")
			claims, err := parseAndValidateToken(token, config.C.User.Secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			userID := claims["user_id"].(float64)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user_id not found in token claims")
			}

			ctx.Set("user_id", uint(userID))

			return handler(ctx)
		}
	}
}

func parseAndValidateToken(tokenString string, signingKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token.Claims.(jwt.MapClaims), nil
	}

	return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims type")
}
