package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"dev/go_final_project/fns"

	_ "github.com/mattn/go-sqlite3"
)

// AddTask добавляет новую задачу в базу данных и возвращает ID добавленной задачи и ошибку, если она возникла.
func AddTask(task Task) (string, error) {
	// Проверка наличия заголовка задачи
	if task.Title == "" {
		return "", errors.New("не указан заголовок задачи")
	}

	now := time.Now()

	// Проверка даты и правила повторения
	if task.Date == "" {
		task.Date = now.Format(fns.FormatDate)
	} else {
		_, err := time.Parse(fns.FormatDate, task.Date)
		if err != nil {
			return "", errors.New("дата представлена в формате, отличном от 20060102")
		}
	}

	if task.Date < now.Format(fns.FormatDate) {
		// задача без повторения
		if task.Repeat == "" {
			task.Date = now.Format(fns.FormatDate)
		} else {
			nextDate, err := fns.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return "", errors.New("правило повторения указано в неправильном формате")
			}
			parsedDate, err := time.Parse(fns.FormatDate, nextDate)
			if err != nil {
				return "", errors.New("ошибка получения следующей даты")
			}
			task.Date = parsedDate.Format(fns.FormatDate)
		}
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	result, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return "", fmt.Errorf("не удалось добавить задачу: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("не удалось получить ID добавленной задачи: %v", err)
	}

	return strconv.FormatInt(id, 10), nil
}
