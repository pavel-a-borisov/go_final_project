package handlers

import (
	"dev/go_final_project/database"
	"dev/go_final_project/fns"
	"fmt"
	"log"
	"net/http"
	"time"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(fns.FormatDate, nowStr)
	if err != nil {
		response := database.Response{ID: "error", Error: fmt.Sprintf("неверный формат текущей даты: %v", err)}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("неверный формат текущей даты: %v", err)
		return
	}

	result, err := fns.NextDate(now, dateStr, repeat)
	if err != nil {
		response := database.Response{ID: "error", Error: fmt.Sprintf("Ошибка: %v", err)}
		ReturnJSON(w, response, http.StatusBadRequest)
		log.Printf("ошибка при расчете следующей даты: %v", err)
		return
	}

	response := map[string]interface{}{"next_date": result}
	ReturnJSON(w, response, http.StatusOK)
}
