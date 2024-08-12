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

// App - структура приложения с базой данных и сервисом задач
type App struct {
	DB          *db.DB
	TaskService service.TaskServiceInterface
}

// Настройка логирования
func setupLogging() *os.File {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ошибка открытия log-файла: %v", err)
	}
	log.SetOutput(logFile)
	return logFile
}

// Загрузка переменных окружения из файла .env
func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки файла .env: %v", err)
	}
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
func getFileDB() string {
	fileDB := os.Getenv("TODO_DBFILE")
	if fileDB == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatalf("не удалось получить путь к приложению: %v", err)
		}
		fileDB = filepath.Join(filepath.Dir(appPath), filepath_db)
	}
	return fileDB
}

// Инициализация приложения и подключение к базе данных
func initializeApp(connectionString string) (*App, error) {
	// Подключаем БД
	conn, err := db.ConnectDB(connectionString)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	// Инициализация объекта DB
	db := db.NewDB(conn)

	// Инициализация объекта TaskService
	taskService := service.NewTaskService(db)

	app := &App{
		DB:          db,
		TaskService: taskService,
	}

	return app, nil
}

func main() {
	// Открываем лог файл
	logFile := setupLogging()
	defer logFile.Close()

	log.SetOutput(logFile)

	// Загрузка переменных окружения из файла .env
	loadEnvVariables()

	// получаем номер порта
	port := getPort()

	// получение пароля из переменной окружения
	pass := os.Getenv("TODO_PASSWORD")

	app, err := initializeApp(getFileDB())
	if err != nil {
		log.Fatalf("ошибка при подключении к базе данных: %v", err)
	}
	defer app.DB.Conn.Close() // Закрытие соединения с базой данных

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
		r.Post("/api/task", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleAddTask(app.TaskService, w, r)
		})
		// обработчик API для вывода ближайших задач и поиска.
		r.Get("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleGetTasks(app.TaskService, w, r)
		})

		// обработчик API для обновления задачи по ID.
		r.Put("/api/task", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleUpdateTask(app.TaskService, w, r)
		})

		// обработчик API для вывода задачи по ID.
		r.Get("/api/task", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleGetTaskByID(app.TaskService, w, r)
		})
		// обработчик API для отметки о выполнении задачи по ID.
		r.Post("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleMarkTaskDone(app.TaskService, w, r)
		})
		// обработчик API для удаления задачи по ID.
		r.Delete("/api/task", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleDeleteTask(app.TaskService, w, r)
		})
	})

	// Открытые маршруты
	// обработчик API для вычисления следующей даты
	r.Get("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleNextDate(app.TaskService, w, r)
	})

	// маршрут для аутентификации
	r.Post("/api/signin", handlers.HandleSignIn)

	fmt.Println("The server started on port: " + port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
