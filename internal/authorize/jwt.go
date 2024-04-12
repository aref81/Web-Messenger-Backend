package authorize

import (
	"backend/internal/configs"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	tokenValidationTime = time.Hour * 12
)

func ValidateJWT(tokenString string) (float64, error) {
	conf, _ := configs.LoadConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.ErrUnauthorized
		}
		return []byte(conf.Server.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, echo.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, echo.ErrUnauthorized
	}
	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, echo.ErrUnauthorized
	}
	return userID, nil
}

func JwtToken(id uint64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id
	claims["exp"] = time.Now().Add(tokenValidationTime).Unix()

	conf, _ := configs.LoadConfig()

	tokenString, err := token.SignedString([]byte(conf.Server.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
