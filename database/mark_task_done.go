package database

import (
	"dev/go_final_project/fns"
	"fmt"
	"time"
)

// MarkTaskDone отмечает задачу как выполненную.
func MarkTaskDone(id int) error {
	task, err := GetTaskByID(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		// Одноразовая задача, удаляем её
		return DeleteTask(id)
	}

	// Периодическая задача, обновляем дату следующего выполнения
	now := time.Now()
	nextDate, err := fns.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("не удалось рассчитать следующую дату выполнения: %v", err)
	}

	task.Date = nextDate
	return UpdateTask(*task)
}
