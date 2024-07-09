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
		http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", http.StatusBadRequest)
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
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON в случае успешного обновления
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}
