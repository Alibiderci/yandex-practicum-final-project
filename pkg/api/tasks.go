package api

import (
	"net/http"
	"yandex-practicum-final-project/pkg/db"
)

type TaskResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, searchQuery)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := TaskResponse{Tasks: tasks}
	writeJson(w, http.StatusOK, response)
}
