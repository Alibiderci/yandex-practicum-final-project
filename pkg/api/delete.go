package api

import (
	"net/http"
	"strings"
	"yandex-practicum-final-project/pkg/db"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})	
		} else {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})	
		}	
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
