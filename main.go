package main

import (
	"dev/go_final_project/database"
	"dev/go_final_project/fns"
	"dev/go_final_project/handlers"
	"dev/go_final_project/tests"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

// Реализация middleware для проверки аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// получение пароля из переменной окружения
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получение куки
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			var valid bool
			// Код для валидации и проверки JWT-токена
			if jwt != "" {
				hash, err := fns.ValidateJWT(jwt)
				if err == nil && hash == "someHashBasedOnPassword" { // Используйте реальную проверку хеша
					valid = true
				} else {
					// Лог для отладки ошибок валидации токена
					if err != nil {
						log.Printf("ошибка при валидации токена: %v", err)
					}
					if hash != "someHashBasedOnPassword" {
						log.Printf("неправильный хеш токена: %s", hash)
					}
				}
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				//http.Error(w, "Authentification required", http.StatusUnauthorized)
				response := database.Response{ID: "error", Error: "аутентификация требуется"}
				handlers.ReturnJSON(w, response, http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func main() {
	// Открываем лог файл
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ошибка открытия log-файла: %v", err) // Используем log.Fatalf() для выхода из программы при ошибке
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	// Загрузка переменных окружения из файла .env
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки файла .env: %v", err)
	}

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

	// Регистрация маршрутов
	r.Post("/api/signin", handlers.HandleSignIn)
	// обработчик API для вычисления следующей даты
	r.Get("/api/nextdate", auth(handlers.HandleNextDate))
	// обработчик API для добавления новой задачи
	r.Post("/api/task", auth(handlers.HandleAddTask))
	// обработчик API для вывода ближайших задач и поиска.
	r.Get("/api/tasks", auth(handlers.HandleGetTasks))
	// обработчик API для обновления задачи по ID.
	r.Put("/api/task", auth(handlers.HandleUpdateTask))
	// обработчик API для вывода задачи по ID.
	r.Get("/api/task", auth(handlers.HandleGetTaskByID))
	// обработчик API для отметки о выполнении задачи по ID.
	r.Post("/api/task/done", auth(handlers.HandleMarkTaskDone))
	// обработчик API для удаления задачи по ID.
	r.Delete("/api/task", auth(handlers.HandleDeleteTask))

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}
