package database

import (
	"database/sql"
	"dev/go_final_project/fns"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// AddTask добавляет новую задачу в базу данных и возвращает ID добавленной задачи и ошибку, если она возникла.
func AddTask(task Task) (int, error) {
	// Проверка наличия заголовка задачи
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	now := time.Now()
	date := now

	// Проверка даты и правила повторения
	if task.Date != "" {
		parsedDate, err := time.Parse(fns.FormatDate, task.Date)
		if err != nil {
			return 0, errors.New("дата представлена в формате, отличном от 20060102")
		}

		if parsedDate.Before(now) {
			if task.Repeat == "" {
				date = now
			} else {
				nextDate, err := fns.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return 0, errors.New("правило повторения указано в неправильном формате")
				}
				parsedDate, err = time.Parse(fns.FormatDate, nextDate)
				if err != nil {
					return 0, errors.New("правило повторения указано в неправильном формате")
				}
				date = parsedDate
			}
		} else {
			date = parsedDate
		}
		//} else {
		//		task.Date =
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	result, err := db.Exec(query,
		sql.Named("date", date.Format(fns.FormatDate)),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return 0, fmt.Errorf("не удалось добавить задачу: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить ID добавленной задачи: %v", err)
	}

	return int(id), nil
}
