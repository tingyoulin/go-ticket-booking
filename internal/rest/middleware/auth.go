package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/tingyoulin/go-ticket-booking/domain"
	"github.com/tingyoulin/go-ticket-booking/internal/repository/redis"
	"github.com/tingyoulin/go-ticket-booking/internal/utils"
)

func Auth(tokenRedisRepo *redis.TokenRedisRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, utils.ResponseError{Message: "Authorization header is required"})
			}

			// check if token is blacklisted
			if isBlacklisted, err := tokenRedisRepo.IsBlacklisted(c.Request().Context(), tokenString); err != nil || isBlacklisted {
				return c.JSON(http.StatusUnauthorized, utils.ResponseError{Message: domain.ErrUnauthorized.Error()})
			}

			// parse token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, utils.ResponseError{Message: domain.ErrUnauthorized.Error()})
			}

			// get claims
			claims, ok := token.Claims.(jwt.MapClaims)
			passengerID, claimsOk := claims["passenger_id"].(float64)
			if !ok || !token.Valid || !claimsOk {
				return c.JSON(http.StatusUnauthorized, utils.ResponseError{Message: domain.ErrUnauthorized.Error()})
			}
			c.Set("passenger_id", int64(passengerID))

			return next(c)
		}
	}
}
