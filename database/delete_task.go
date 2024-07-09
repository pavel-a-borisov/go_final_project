package database

import (
	"database/sql"
	"fmt"
)

// DeleteTask удаляет задачу из базы данных.
func DeleteTask(id int) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	result, err := db.Exec(query, sql.Named("id", id))
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
