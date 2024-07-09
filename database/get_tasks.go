package database

import (
	"database/sql"
	"dev/go_final_project/fns"
	"fmt"
	"strings"
	"time"
)

// GetTasks возвращает список задач, отсортированных по дате в сторону увеличения и опциональным поиском по заголовку и комментарию.
func GetTasks(search string, limit int) ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	var args []interface{}

	if search != "" {
		// Попробуем сначала парсить поисковую строку как дату
		if parsedDate, err := time.Parse(fns.ShortFormatDate, search); err == nil {
			// Если удалось распарсить как дату, добавляем условие поиска по дате
			//query += ` WHERE date = ?`
			//args = append(args, parsedDate.Format(fns.FormatDate))
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date ASC LIMIT :limit`
			args = append(args, sql.Named("date", parsedDate.Format(fns.FormatDate)), sql.Named("limit", limit))
		} else {
			// Иначе добавляем условия поиска по заголовку и комментарию
			//query += ` WHERE title LIKE ? OR comment LIKE ?`
			//searchPattern := "%" + search + "%"
			//args = append(args, searchPattern, searchPattern)
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE :search OR LOWER(comment) LIKE :search ORDER BY date ASC LIMIT :limit`
			searchPattern := "%" + strings.ToLower(search) + "%"
			args = append(args, sql.Named("search", searchPattern), sql.Named("limit", limit))
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("не удалось прочитать задачу: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %v", err)
	}

	return tasks, nil
}
