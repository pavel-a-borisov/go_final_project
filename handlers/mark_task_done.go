package handlers

import (
	"dev/go_final_project/database"
	"dev/go_final_project/fns"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Обработчик POST-запроса для отметки задачи как выполненной
func HandleMarkTaskDone(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := database.Response{ID: "error", Error: "Не указан идентификатор при отметке задачи"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("не указан идентификатор при отметке задачи: %v", http.StatusBadRequest)
		return
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		response := database.Response{ID: "error", Error: "неправильный формат идентификатора при отметке задачи"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("неправильный формат идентификатора при отметке задачи: %v", err)
		return
	}

	// Находим задачу
	task, err := database.GetTaskByID(idStr)
	if err != nil {
		response := database.Response{ID: "error", Error: "задача не найдена"}
		ReturnJSON(w, response, http.StatusInternalServerError)
		log.Printf("задача не найдена: %v", err)
		return
	}

	// Если задача одноразовая, удаляем её
	if task.Repeat == "" {
		err = database.DeleteTask(idStr)
		if err != nil {
			http.Error(w, `{"error":"ошибка при удалении одноразовой задачи"}`, http.StatusInternalServerError)
			log.Printf("%v", err)
			return
		}

	} else { // Периодическая задача, обновляем дату следующего выполнения
		now := time.Now()
		nextDate, err := fns.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"не удалось рассчитать следующую дату выполнения"}`, http.StatusInternalServerError)
			log.Printf("не удалось рассчитать следующую дату выполнения: %v", err)
			return
		}

		// Обновляем задачу с новой датой
		task.Date = nextDate
		err = database.UpdateTask(*task)
		if err != nil {
			http.Error(w, `{"error":"не удалось обновить задачу с новой датаой"}`, http.StatusInternalServerError)
			log.Printf("%v", err)
			return
		}
	}

	// Возвращаем пустой JSON в случае успешного обновления
	response := database.Response{ID: "", Error: ""}
	ReturnJSON(w, response, http.StatusOK)
}
