package main

import (
	"dev/go_final_project/database"
	"dev/go_final_project/handlers"
	"dev/go_final_project/tests"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Открываем лог файл
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ошибка открытия log-файла: %v", err) // Используем log.Fatalf() для выхода из программы при ошибке
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	// Подключаем БД
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("ошибка при подключении к базе данных: %v", err)
	}

	defer db.Close() // Закрываем соединение при завершении работы

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	// Создаём новый роутер chi
	r := chi.NewRouter()
	// Добавляем встроенные middleware для логирования и восстановления после паник
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	// обработчик API для вычисления следующей даты
	r.Get("/api/nextdate", handlers.HandleNextDate)
	// обработчик API для добавления новой задачи
	r.Post("/api/task", handlers.HandleAddTask)
	// обработчик API для вывода ближайших задач и поиска.
	r.Get("/api/tasks", handlers.HandleGetTasks)
	// обработчик API для обновления задачи по ID.
	r.Put("/api/task", handlers.HandleUpdateTask)
	// обработчик API для вывода задачи по ID.
	r.Get("/api/task", handlers.HandleGetTaskByID)
	// обработчик API для отметки о выполнении задачи по ID.
	r.Post("/api/task/done", handlers.HandleMarkTaskDone)
	// обработчик API для удаления задачи по ID.
	r.Delete("/api/task", handlers.HandleDeleteTask)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}
