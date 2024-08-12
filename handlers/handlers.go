package handlers

import (
	"dev/go_final_project/model"
	"dev/go_final_project/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// ограничение на количество задач при выводе списка задач
const (
	limit      = 50
	FormatDate = "20060102"
)

// Функция для централизованного обработки JSON ответов
func ReturnJSON(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
	}
}

// HandleNextDate - обработчик GET-запроса для вычисления следующей даты для повторяемых задач
func HandleNextDate(taskService service.TaskServiceInterface, w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(FormatDate, nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Неверный формат текущей даты: %v", err), http.StatusBadRequest)
		return
	}

	result, err := taskService.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}

// HandleAddTask - обработчик POST-запроса для добавления задачи
func HandleAddTask(taskService service.TaskServiceInterface, w http.ResponseWriter, r *http.Request) {
	var task model.Task

	// Декодирование JSON-запроса в структуру Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response := model.Response{ID: "error", Error: "Ошибка десериализации JSON"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", err)
		return
	}

	// Добавление задачи в базу данных
	id, err := taskService.AddTask(task)
	if err != nil {
		response := model.Response{ID: "error", Error: fmt.Sprintf("ошибка добавления задачи в базу: %v", err)}
		ReturnJSON(w, response, http.StatusInternalServerError)
		log.Printf("ошибка добавления задачи в базу: %v", err)
		return
	}

	// Используем returnJSON для возврата ID задачи
	response := map[string]interface{}{"id": id}
	ReturnJSON(w, response, http.StatusOK)
}

// HandleGetTasks - обработчик GET-запроса для получения списка задач
func HandleGetTasks(taskService service.TaskServiceInterface, w http.ResponseWriter, r *http.Request) {
	// Получаем параметр search из строки запроса
	search := r.URL.Query().Get("search")

	tasks, err := taskService.GetTasks(search, limit)
	if err != nil {
		response := model.Response{ID: "error", Error: fmt.Sprintf("ошибка получения задач: %v", err)}
		ReturnJSON(w, response, http.StatusInternalServerError)
		log.Printf("oшибка при получении задач: %v", err)
		return
	}

	// выводим задачи
	response := map[string]interface{}{"tasks": tasks}
	ReturnJSON(w, response, http.StatusOK)
}

// HandleUpdateTask - обработчик PUT-запроса для обновления задачи
func HandleUpdateTask(taskService service.TaskServiceInterface, w http.ResponseWriter, r *http.Request) {
	var task model.Task

	// Декодирование JSON-запроса в структуру Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response := model.Response{ID: "error", Error: "ошибка десериализации JSON"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", err)
		return
	}

	// Обновление задачи в базе данных
	err := taskService.UpdateTask(model.Task{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	})
	if err != nil {
		response := model.Response{ID: "error", Error: fmt.Sprintf("ошибка при обновлении данных: %v", err)}
		ReturnJSON(w, response, http.StatusNotFound)
		log.Printf("ошибка при обновлении данных: %v", err)
		return
	}

	// Возвращаем пустой JSON в случае успешного обновления
	response := model.Response{ID: "", Error: ""}
	ReturnJSON(w, response, http.StatusOK)
}

// HandleGetTaskByID - обработчик GET-запроса для получения задачи по идентификатору
func HandleGetTaskByID(taskService service.TaskServiceInterface, w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := model.Response{ID: "error", Error: "не указан идентификатор"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("не указан идентификатор")
		return
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		response := model.Response{ID: "error", Error: "неправильный формат идентификатора"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("неправильный формат идентификатора: %v", err)
		return
	}

	task, err := taskService.GetTaskByID(idStr)
	if err != nil {
		response := model.Response{ID: "error", Error: "Ошибка при получении данных из базы"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка при получении данных из базы: %v", err)
		return
	}

	// возвращаем задачу
	ReturnJSON(w, task, http.StatusOK)
}

// HandleMarkTaskDone - обработчик POST-запроса для отметки задачи как выполненной
func HandleMarkTaskDone(taskService service.TaskServiceInterface, w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := model.Response{ID: "error", Error: "Не указан идентификатор при отметке задачи"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("не указан идентификатор при отметке задачи: %v", http.StatusBadRequest)
		return
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		response := model.Response{ID: "error", Error: "неправильный формат идентификатора при отметке задачи"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("неправильный формат идентификатора при отметке задачи: %v", err)
		return
	}

	// Находим задачу
	task, err := taskService.GetTaskByID(idStr)
	if err != nil {
		response := model.Response{ID: "error", Error: "задача не найдена"}
		ReturnJSON(w, response, http.StatusInternalServerError)
		log.Printf("задача не найдена: %v", err)
		return
	}

	// Если задача одноразовая, удаляем её
	if task.Repeat == "" {
		err = taskService.DeleteTask(idStr)
		if err != nil {
			http.Error(w, `{"error":"ошибка при удалении одноразовой задачи"}`, http.StatusInternalServerError)
			log.Printf("%v", err)
			return
		}

	} else { // Периодическая задача, обновляем дату следующего выполнения
		now := time.Now()
		nextDate, err := service.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"не удалось рассчитать следующую дату выполнения"}`, http.StatusInternalServerError)
			log.Printf("не удалось рассчитать следующую дату выполнения: %v", err)
			return
		}

		// Обновляем задачу с новой датой
		task.Date = nextDate
		err = taskService.UpdateTask(*task)
		if err != nil {
			http.Error(w, `{"error":"не удалось обновить задачу с новой датаой"}`, http.StatusInternalServerError)
			log.Printf("%v", err)
			return
		}
	}

	// Возвращаем пустой JSON в случае успешного обновления
	response := model.Response{ID: "", Error: ""}
	ReturnJSON(w, response, http.StatusOK)
}

// HandleDeleteTask - обработчик DELETE-запроса для удаления задачи по идентификатору
func HandleDeleteTask(taskService service.TaskServiceInterface, w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из строки запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := model.Response{ID: "error", Error: "не указан идентификатор для удаления задачи из базы"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("не указан идентификатор для удаления задачи из базы: %v", http.StatusBadRequest)
		return
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		response := model.Response{ID: "error", Error: "неправильный формат идентификатора для удаления задачи из базы"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("неправильный формат идентификатора для удаления задачи из базы: %v", err)
		return
	}

	// Удаляем задачу из базы данных
	err = taskService.DeleteTask(idStr)
	if err != nil {
		response := model.Response{ID: "error", Error: "ошибка при удалении задачи из базы"}
		ReturnJSON(w, response, http.StatusInternalServerError)
		log.Printf("ошибка при удалении задачи из базы: %v", err)
		return

	}

	// Возвращаем пустой JSON в случае успешного удаления
	response := model.Response{ID: "", Error: ""}
	ReturnJSON(w, response, http.StatusOK)
}

// HandleSignIn - обработчик POST-запроса для аутентификации пользователя
func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Password string `json:"password"`
	}

	// Декодирование JSON-запроса
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		response := model.Response{ID: "error", Error: "Ошибка десериализации JSON"}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка десериализации JSON: %v", err)
		return
	}

	password := os.Getenv("TODO_PASSWORD")

	// Проверка введенного пароля
	if creds.Password != password {
		response := model.Response{ID: "error", Error: "Неверный пароль"}
		ReturnJSON(w, response, http.StatusUnauthorized)
		log.Printf("неверный пароль: %v", creds.Password)
		return
	}

	// Создание JWT
	hash := "someHashBasedOnPassword" // Используйте реальную функцию генерации хеша
	token, err := CreateJWT(hash)
	if err != nil {
		response := model.Response{ID: "error", Error: "Не удалось создать токен"}
		ReturnJSON(w, response, http.StatusInternalServerError)
		log.Printf("ошибка создания токена: %v", err)
		return
	}

	// Установка куки с JWT-токеном
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(8 * time.Hour),
	})

	// Возвращение успешного ответа
	response := map[string]string{"token": token}
	ReturnJSON(w, response, http.StatusOK)
}
