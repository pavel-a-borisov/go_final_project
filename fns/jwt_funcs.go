package fns

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Переменная, содержащая секретный ключ для подписи JWT
var jwtKey = []byte(os.Getenv("TODO_PASSWORD"))

// Структура для JWT-клейма
type Claims struct {
	Hash string `json:"hash"`
	jwt.StandardClaims
}

// Функция для создания JWT
func CreateJWT(hash string) (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		Hash: hash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Функция для валидации и парсинга JWT
func ValidateJWT(tokenStr string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Hash, nil
}
