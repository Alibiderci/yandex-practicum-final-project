package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"yandex-practicum-final-project/pkg/db"
)

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(layout)
	}

	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате отличном от 20060102")
	}


	if !afterNow(t, now) && t.Format(layout) != now.Format(layout) {
		if task.Repeat == "" {
			task.Date = now.Format(layout) 
		} else {
			next, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	} else {
		if task.Repeat != "" {
			_, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err 

			}
		}

	}
	return nil
}


func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {	
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}


	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return

	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "ошибка при добавлении задачи в базу данных"})
		return
	}

	writeJson(w, http.StatusCreated, map[string]int64{"id": id})
}
