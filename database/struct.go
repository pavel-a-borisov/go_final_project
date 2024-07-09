package database

// Структура для задачи
type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Структура для ответа
type Response struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// Структура для списка задач
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
