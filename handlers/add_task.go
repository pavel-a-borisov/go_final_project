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
		response := database.Response{ID: "error", Error: "Ошибка десериализации JSON"}
		returnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", err)
		return
	}

	// Добавление задачи в базу данных
	id, err := database.AddTask(task)
	if err != nil {
		response := database.Response{ID: "error", Error: fmt.Sprintf("ошибка добавления задачи в базу: %v", err)}
		returnJSON(w, response, http.StatusInternalServerError)
		log.Printf("ошибка добавления задачи в базу: %v", err)
		return
	}

	// Используем returnJSON для возврата ID задачи
	response := map[string]interface{}{"id": id}
	returnJSON(w, response, http.StatusOK)
}
