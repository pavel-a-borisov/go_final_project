package database

import (
	"database/sql"
	"dev/go_final_project/fns"
	"fmt"
	"strconv"
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
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date ASC LIMIT :limit`
			args = append(args, sql.Named("date", parsedDate.Format(fns.FormatDate)), sql.Named("limit", limit))
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

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("не удалось прочитать задачу: %v", err)
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %v", err)
	}

	// Возвращаем пустой массив задач, если ни одной задачи не найдено
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}
