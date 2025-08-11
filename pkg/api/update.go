package api

import (
	"encoding/json"
	"net/http"
	"yandex-practicum-final-project/pkg/db"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	task := db.Task{}

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "идентификатор не может быть пустым"})
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

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusCreated, map[string]string{})
}
