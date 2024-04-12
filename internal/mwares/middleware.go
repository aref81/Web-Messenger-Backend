package mwares

import (
	"backend/internal/authorize"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication token not found")
		}

		authSplitToken := strings.Split(authHeader, "Bearer ")
		if len(authSplitToken) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication token format is not valid")
		}

		claims, err := authorize.ValidateJWT(authSplitToken[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication token is not valid")
		}

		c.Set("userID", claims)
		return next(c)
	}
}
