package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func ConnectDB() (*sql.DB, error) {
	// Получение пути к файлу базы данных из переменной окружения или директории приложения
	dsn := os.Getenv("TODO_DBFILE")
	if dsn == "" {
		appPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("не удалось получить путь к приложению: %v", err)
		}
		dsn = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	install := false
	if _, err := os.Stat(dsn); err != nil {
		install = true
	}

	// Открываем базу данных.
	var err error
	db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}

	// Проверяем соединение.
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
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
			return nil, fmt.Errorf("ошибка при создании таблицы или индекса: %v", err)
		} else {
			log.Println("таблица и индекс успешно созданы.")
		}
	}

	return db, nil
}
