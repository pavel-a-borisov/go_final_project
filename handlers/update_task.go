package handlers

import (
	"dev/go_final_project/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Обработчик PUT-запроса для обновления задачи
func HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task database.Task

	// Декодирование JSON-запроса в структуру Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response := database.Response{ID: "error", Error: "ошибка десериализации JSON"}
		returnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", err)
		return
	}

	// Обновление задачи в базе данных
	err := database.UpdateTask(database.Task{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	})
	if err != nil {
		response := database.Response{ID: "error", Error: fmt.Sprintf("ошибка при обновлении данных: %v", err)}
		returnJSON(w, response, http.StatusNotFound)
		log.Printf("ошибка при обновлении данных: %v", err)
		return
	}

	// Возвращаем пустой JSON в случае успешного обновления
	response := database.Response{ID: "", Error: ""}
	returnJSON(w, response, http.StatusOK)
}
