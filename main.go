package main

import (
	"dev/go_final_project/db"
	"dev/go_final_project/handlers"
	"dev/go_final_project/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	server_port = "7540"
	filepath_db = "scheduler.db"
)

// Настройка логирования
func setupLogging() (*os.File, error) {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия log-файла: %v", err)
	}
	log.SetOutput(logFile)
	return logFile, nil
}

// Загрузка переменных окружения из файла .env
func loadEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла .env: %v", err)
	}
	return nil
}

// Получение порта из переменных окружения
func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = server_port
	}
	return port
}

// Получение пути к файлу базы данных из переменной окружения или директории приложения
func getFileDB() (string, error) {
	fileDB := os.Getenv("TODO_DBFILE")
	if fileDB == "" {
		appPath, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("не удалось получить путь к приложению: %v", err)
		}
		fileDB = filepath.Join(filepath.Dir(appPath), filepath_db)
	}
	return fileDB, nil
}

// Инициализация приложения и подключение к базе данных
func initializeApp(connectionString string) (*handlers.App, error) {
	// Подключаем БД
	conn, err := db.ConnectDB(connectionString)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err.Error())
	}

	// Инициализация объекта DB
	db := db.NewDB(conn)

	// Инициализация объекта TaskService
	taskService := service.NewTaskService(db)

	app := &handlers.App{
		TaskService: taskService,
	}

	return app, nil
}

func main() {
	// Открываем лог файл
	logFile, err := setupLogging()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Загрузка переменных окружения из файла .env
	err = loadEnvVariables()
	if err != nil {
		log.Fatalf(err.Error())
	}
	// получаем номер порта
	port := getPort()

	// получение пароля из переменной окружения
	pass := os.Getenv("TODO_PASSWORD")

	file_db, err := getFileDB()
	if err != nil {
		log.Fatalf(err.Error())
	}

	app, err := initializeApp(file_db)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer app.TaskService.(*service.TaskService).DB.Conn.Close() // Закрытие соединения с базой данных

	// Создаём новый роутер chi
	r := chi.NewRouter()
	// Добавляем встроенные middleware для логирования и восстановления после паник
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	// Регистрация маршрутов
	// Маршруты, требующие аутентификации
	r.Group(func(r chi.Router) {
		r.Use(handlers.Auth(pass)) // Применение middleware для аутентификации

		// обработчик API для добавления новой задачи
		r.Post("/api/task", app.HandleAddTask)

		// обработчик API для вывода ближайших задач и поиска.
		r.Get("/api/tasks", app.HandleGetTasks)

		// обработчик API для обновления задачи по ID.
		r.Put("/api/task", app.HandleUpdateTask)

		// обработчик API для вывода задачи по ID.
		r.Get("/api/task", app.HandleGetTaskByID)

		// обработчик API для отметки о выполнении задачи по ID.
		r.Post("/api/task/done", app.HandleMarkTaskDone)
		// обработчик API для удаления задачи по ID.
		r.Delete("/api/task", app.HandleDeleteTask)
	})

	// Открытые маршруты
	// обработчик API для вычисления следующей даты
	r.Get("/api/nextdate", app.HandleNextDate)

	// маршрут для аутентификации
	r.Post("/api/signin", handlers.HandleSignIn)

	fmt.Println("The server started on port: " + port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
