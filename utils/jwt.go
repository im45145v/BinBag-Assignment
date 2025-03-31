package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/im45145v/BinBag-Assignment/config"
)

func GenerateToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecretKey))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.JWTSecretKey), nil
	})
}
