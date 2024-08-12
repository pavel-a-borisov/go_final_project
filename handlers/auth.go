package handlers

import (
	"dev/go_final_project/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
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

// Реализация middleware для проверки аутентификации
func Auth(pass string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(pass) > 0 {
				var jwt string // JWT-токен из куки
				// получение куки
				cookie, err := r.Cookie("token")
				if err == nil {
					jwt = cookie.Value
				}

				var valid bool
				// Код для валидации и проверки JWT-токена
				if jwt != "" {
					hash, err := ValidateJWT(jwt)
					if err == nil && hash == "someHashBasedOnPassword" { // Используйте реальную проверку хеша
						valid = true
					} else {
						// Лог для отладки ошибок валидации токена
						if err != nil {
							log.Printf("ошибка при валидации токена: %v", err)
						}
						if hash != "someHashBasedOnPassword" {
							log.Printf("неправильный хеш токена: %s", hash)
						}
					}
				}

				if !valid {
					// возвращаем ошибку авторизации 401
					response := model.Response{ID: "error", Error: "требуется аутентификация"}
					ReturnJSON(w, response, http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
