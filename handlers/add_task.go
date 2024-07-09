package handlers

import (
	"dev/go_final_project/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Обработчик POST-запроса для добавления задачи
func HandleAddTask(w http.ResponseWriter, r *http.Request) {
	var task database.Task

	// Декодирование JSON-запроса в структуру Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", http.StatusBadRequest)
		return
	}

	// Добавление задачи в базу данных
	id, err := database.AddTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка добавления задачи в базу: %v", err), http.StatusInternalServerError)
		log.Printf("ошибка добавления задачи в базу: %v", err)
		return
	}

	// Возвращаем ID задачи
	response := database.Response{ID: id}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}
