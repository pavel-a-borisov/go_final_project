package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

// GetTaskByID возвращает задачу по её идентификатору.
func GetTaskByID(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`
	row := db.QueryRow(query, sql.Named("id", id))

	var task Task
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
