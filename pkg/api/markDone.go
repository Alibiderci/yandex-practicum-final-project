package api

import (
	"net/http"
	"strings"
	"time"
	"yandex-practicum-final-project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")	
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})	
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		} else {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)	

		if err != nil {
			if strings.Contains(err.Error(), "не найдена") {
				writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			} else {
				writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			return
		}
	} else {
		now := time.Now()

		next, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			if strings.Contains(err.Error(), "не найдена") {
				writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			} else {
				writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			return
		}
	}

	writeJson(w, http.StatusCreated, map[string]string{})
}
