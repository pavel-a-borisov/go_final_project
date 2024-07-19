package handlers

import (
	"dev/go_final_project/database"
	"dev/go_final_project/fns"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// HandleSignIn handles the signin POST request
func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Password string `json:"password"`
	}

	// Декодирование JSON-запроса
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		response := database.Response{ID: "error", Error: "Ошибка десериализации JSON"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", err)
		return
	}

	password := os.Getenv("TODO_PASSWORD")

	// Проверка введенного пароля
	if creds.Password != password {
		response := database.Response{ID: "error", Error: "Неверный пароль"}
		ReturnJSON(w, response, http.StatusUnauthorized)
		log.Printf("неверный пароль: %v", creds.Password)
		return
	}

	// Создание JWT
	hash := "someHashBasedOnPassword" // Используйте реальную функцию генерации хеша
	token, err := fns.CreateJWT(hash)
	if err != nil {
		response := database.Response{ID: "error", Error: "Не удалось создать токен"}
		ReturnJSON(w, response, http.StatusInternalServerError)
		log.Printf("ошибка создания токена: %v", err)
		return
	}

	// Установка куки с JWT-токеном
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(8 * time.Hour),
	})

	// Возвращение успешного ответа
	response := map[string]string{"token": token}
	ReturnJSON(w, response, http.StatusOK)
}
