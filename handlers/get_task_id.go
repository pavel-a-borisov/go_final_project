package handlers

import (
	"dev/go_final_project/database"
	"log"
	"net/http"
	"strconv"
)

// Обработчик GET-запроса для получения задачи по идентификатору
func HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := database.Response{ID: "error", Error: "не указан идентификатор"}
		returnJSON(w, response, http.StatusBadRequest)
		log.Printf("не указан идентификатор")
		return
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		response := database.Response{ID: "error", Error: "неправильный формат идентификатора"}
		returnJSON(w, response, http.StatusBadRequest)
		log.Printf("неправильный формат идентификатора: %v", err)
		return
	}

	task, err := database.GetTaskByID(idStr)
	if err != nil {
		response := database.Response{ID: "error", Error: "Ошибка при получении данных из базы"}
		returnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка при получении данных из базы: %v", err)
		return
	}

	// возвращаем задачу
	returnJSON(w, task, http.StatusOK)

}
