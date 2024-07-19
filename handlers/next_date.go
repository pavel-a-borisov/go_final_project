package handlers

import (
	"dev/go_final_project/fns"
	"fmt"
	"net/http"
	"time"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(fns.FormatDate, nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Неверный формат текущей даты: %v", err), http.StatusBadRequest)
		return
	}

	result, err := fns.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}
