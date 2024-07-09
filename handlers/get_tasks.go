package handlers

import (
	"dev/go_final_project/database"
	"encoding/json"
	"fmt"
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
		http.Error(w, fmt.Sprintf("Ошибка получения задач: %v", err), http.StatusInternalServerError)
		return
	}

	response := database.TasksResponse{Tasks: tasks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
