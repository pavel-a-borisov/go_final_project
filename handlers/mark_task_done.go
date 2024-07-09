package handlers

import (
	"dev/go_final_project/database"
	"fmt"
	"net/http"
	"strconv"
)

// Обработчик POST-запроса для отметки задачи как выполненной
func HandleMarkTaskDone(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Неправильный формат идентификатора"}`, http.StatusBadRequest)
		return
	}

	// Отмечаем задачу как выполненную в базе данных
	err = database.MarkTaskDone(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON в случае успешного выполнения
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}
