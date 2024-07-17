package handlers

import (
	"dev/go_final_project/database"
	"fmt"
	"log"
	"net/http"
)

// Обработчик GET-запроса для получения списка задач
func HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем ограничение на количество задач
	const limit = 50
	// Получаем параметр search из строки запроса
	search := r.URL.Query().Get("search")

	tasks, err := database.GetTasks(search, limit)
	if err != nil {
		response := database.Response{ID: "error", Error: fmt.Sprintf("ошибка получения задач: %v", err)}
		returnJSON(w, response, http.StatusInternalServerError)
		log.Printf("oшибка при получении задач: %v", err)
		return
	}

	// выводим задачи
	response := map[string]interface{}{"tasks": tasks}
	returnJSON(w, response, http.StatusOK)
}
