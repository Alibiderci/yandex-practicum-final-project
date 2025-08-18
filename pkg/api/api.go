package api

import (
	"net/http"
	"os"
)

var Port string = os.Getenv("TODO_PORT")
var webDir string = os.Getenv("TODO_WEBDIR")

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)	
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	default: 
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func Init() {

	if webDir == "" {
		webDir = "web"
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
	http.HandleFunc("/api/signin", loginHandler)
}

