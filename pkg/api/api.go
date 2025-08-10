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
		AddTaskHandler(w, r)	
	default: 
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func Init() {

	if webDir == "" {
		webDir = "web"
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
}

