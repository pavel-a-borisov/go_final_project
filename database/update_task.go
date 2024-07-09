package database

import (
	"database/sql"
	"dev/go_final_project/fns"
	"fmt"
	"time"
)

// UpdateTask обновляет существующую задачу в базе данных.
func UpdateTask(task Task) error {
	// Проверка наличия заголовка задачи
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	now := time.Now()

	// Проверка даты и правила повторения
	if task.Date != "" {
		parsedDate, err := time.Parse(fns.FormatDate, task.Date)
		if err != nil {
			return fmt.Errorf("дата представлена в формате, отличном от 20060102")
		}

		if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(fns.FormatDate)
			} else {
				nextDate, err := fns.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return fmt.Errorf("правило повторения указано в неправильном формате")
				}
				task.Date = nextDate
			}
		}
	}

	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	result, err := db.Exec(query,
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
