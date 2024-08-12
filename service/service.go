package service

import (
	"errors"
	"fmt"
	"time"

	"dev/go_final_project/db"
	"dev/go_final_project/model"
)

type TaskServiceInterface interface {
	NextDate(now time.Time, startDate, repeat string) (string, error)
	AddTask(task model.Task) (string, error)
	GetTasks(search string, limit int) ([]model.Task, error)
	UpdateTask(task model.Task) error
	GetTaskByID(idStr string) (*model.Task, error)
	MarkTaskDone(id string) error
	DeleteTask(idStr string) error
}

type TaskService struct {
	DB *db.DB
}

func (s *TaskService) NextDate(now time.Time, startDate, repeat string) (string, error) {
	return NextDate(now, startDate, repeat)
}

func NewTaskService(db *db.DB) *TaskService {
	return &TaskService{DB: db}
}

func (s *TaskService) AddTask(task model.Task) (string, error) {
	// Проверка наличия заголовка задачи
	if task.Title == "" {
		return "", errors.New("не указан заголовок задачи")
	}

	now := time.Now()

	// Проверка даты и правила повторения
	if task.Date == "" {
		task.Date = now.Format(FormatDate)
	} else {
		_, err := time.Parse(FormatDate, task.Date)
		if err != nil {
			return "", errors.New("дата представлена в формате, отличном от 20060102")
		}
	}

	if task.Date < now.Format(FormatDate) {
		// задача без повторения
		if task.Repeat == "" {
			task.Date = now.Format(FormatDate)
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return "", errors.New("правило повторения указано в неправильном формате")
			}
			parsedDate, err := time.Parse(FormatDate, nextDate)
			if err != nil {
				return "", errors.New("ошибка получения следующей даты")
			}
			task.Date = parsedDate.Format(FormatDate)
		}
	}

	return s.DB.AddTask(task)
}

func (s *TaskService) GetTasks(search string, limit int) ([]model.Task, error) {
	return s.DB.GetTasks(search, limit)
}

func (s *TaskService) UpdateTask(task model.Task) error {
	// Проверка наличия заголовка задачи
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	now := time.Now()

	// Проверка даты и правила повторения
	if task.Date != "" {
		parsedDate, err := time.Parse(FormatDate, task.Date)
		if err != nil {
			return fmt.Errorf("дата представлена в формате, отличном от 20060102")
		}

		if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(FormatDate)
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return fmt.Errorf("правило повторения указано в неправильном формате")
				}
				task.Date = nextDate
			}
		}
	}
	return s.DB.UpdateTask(task)
}

func (s *TaskService) GetTaskByID(idStr string) (*model.Task, error) {
	return s.DB.GetTaskByID(idStr)
}

func (s *TaskService) MarkTaskDone(id string) error {
	task, err := s.DB.GetTaskByID(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return s.DB.DeleteTask(id)
	}

	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("не удалось рассчитать следующую дату выполнения: %v", err)
	}

	task.Date = nextDate
	return s.DB.UpdateTask(*task)
}

func (s *TaskService) DeleteTask(idStr string) error {
	return s.DB.DeleteTask(idStr)
}
