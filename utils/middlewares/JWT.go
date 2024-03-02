package middlewares

import (
	"emailnotifl3n/app/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(config.JWT_SECRET),
		SigningMethod: "HS256",
	})
}

// Generate token jwt
func CreateToken(userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT_SECRET))

}

// extract token jwt
func ExtractTokenUserId(e echo.Context) int {
	header := e.Request().Header.Get("Authorization")
	headerToken := strings.Split(header, " ")
	token := headerToken[len(headerToken)-1]
	tokenJWT, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET), nil
	})

	if tokenJWT.Valid {
		claims := tokenJWT.Claims.(jwt.MapClaims)
		userId, isValidUserId := claims["userId"].(float64)
		if !isValidUserId {
			return 0
		}
		return int(userId)
	}
	return 0
}
