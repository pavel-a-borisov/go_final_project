package db

import (
	"database/sql"
	"dev/go_final_project/model"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	FormatDate      = "20060102"
	ShortFormatDate = "02.01.2006"
)

type DB struct {
	Conn *sql.DB
}

func NewDB(conn *sql.DB) *DB {
	return &DB{Conn: conn}
}

func ConnectDB(dsn string) (*sql.DB, error) {

	install := false
	if _, err := os.Stat(dsn); err != nil {
		install = true
	}

	// Открываем базу данных.
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// Проверяем соединение.
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// Создаём таблицу, если это необходимо.
	if install {
		createTableSQL := `
        CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT '',
			title VARCHAR(256) NOT NULL DEFAULT '',
			comment TEXT NOT NULL DEFAULT '',
			repeat VARCHAR(128) NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_datе ON scheduler (date);
        `
		_, err = db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании таблицы или индекса: %w", err)
		} else {
			log.Println("таблица и индекс успешно созданы.")
		}
	}

	return db, nil
}

// AddTask добавляет новую задачу в базу данных и возвращает ID добавленной задачи и ошибку, если она возникла.
func (db *DB) AddTask(task model.Task) (string, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	result, err := db.Conn.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return "", fmt.Errorf("не удалось добавить задачу: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("не удалось получить ID добавленной задачи: %w", err)
	}

	return strconv.FormatInt(id, 10), nil
}

// GetTasks возвращает список задач, отсортированных по дате в сторону увеличения и опциональным поиском по заголовку и комментарию.
func (db *DB) GetTasks(search string, limit int) ([]model.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	var args []interface{}

	if search != "" {
		// Попробуем сначала парсить поисковую строку как дату
		if parsedDate, err := time.Parse(ShortFormatDate, search); err == nil {
			// Если удалось распарсить как дату, добавляем условие поиска по дате
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date ASC LIMIT :limit`
			args = append(args, sql.Named("date", parsedDate.Format(FormatDate)), sql.Named("limit", limit))
		} else {
			// Иначе добавляем условия поиска по заголовку и комментарию
			// Из-за того, что SQLite не умеет обрабатывать кириллицу
			// с помощью встроенных функций LOWER() и UPPER()
			// я отказался от перевода значений полей в верхний/нижний регистр.
			// Сравниваю, как есть.

			//query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE :search OR LOWER(comment) LIKE :search ORDER BY date ASC LIMIT :limit`
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date ASC LIMIT :limit`
			//searchPattern := "%" + strings.ToLower(search) + "%"
			searchPattern := "%" + search + "%"
			args = append(args, sql.Named("search", searchPattern), sql.Named("limit", limit))
		}
	}

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи: %w", err)
	}
	defer rows.Close()

	tasks := []model.Task{} // присваиваем пустой массив задач
	for rows.Next() {
		var task model.Task
		var id int64
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("не удалось прочитать задачу: %w", err)
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %w", err)
	}

	return tasks, nil
}

// UpdateTask обновляет существующую задачу в базе данных.
func (db *DB) UpdateTask(task model.Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	result, err := db.Conn.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return fmt.Errorf("не удалось обновить задачу: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество затронутых строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// GetTaskByID возвращает задачу по её идентификатору.
func (db *DB) GetTaskByID(id string) (*model.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`
	row := db.Conn.QueryRow(query, sql.Named("id", id))

	var task model.Task
	var intID int64
	err := row.Scan(&intID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка при получении задачи: %v", err)
	}
	task.ID = strconv.FormatInt(intID, 10)

	return &task, nil
}

// DeleteTask удаляет задачу из базы данных.
func (db *DB) DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	result, err := db.Conn.Exec(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("не удалось удалить задачу: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество затронутых строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
