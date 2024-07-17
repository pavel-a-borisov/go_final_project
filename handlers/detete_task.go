package handlers

import (
	"dev/go_final_project/database"
	"log"
	"net/http"
	"strconv"
)

// Обработчик DELETE-запроса для удаления задачи по идентификатору
func HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := database.Response{ID: "error", Error: "не указан идентификатор для удаления задачи из базы"}
		returnJSON(w, response, http.StatusBadRequest)
		log.Printf("не указан идентификатор для удаления задачи из базы: %v", http.StatusBadRequest)
		return
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		response := database.Response{ID: "error", Error: "неправильный формат идентификатора для удаления задачи из базы"}
		returnJSON(w, response, http.StatusBadRequest)
		log.Printf("неправильный формат идентификатора для удаления задачи из базы: %v", err)
		return
	}

	// Удаляем задачу из базы данных
	err = database.DeleteTask(idStr)
	if err != nil {
		response := database.Response{ID: "error", Error: "ошибка при удалении задачи из базы"}
		returnJSON(w, response, http.StatusInternalServerError)
		log.Printf("ошибка при удалении задачи из базы: %v", err)
		return

	}

	// Возвращаем пустой JSON в случае успешного удаления
	response := database.Response{ID: "", Error: ""}
	returnJSON(w, response, http.StatusOK)
}
