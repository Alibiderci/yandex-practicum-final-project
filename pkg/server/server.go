package server

import (
	"net/http"
	"os"
)

func Start() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir := os.Getenv("TODO_WEBDIR")
	if webDir == "" {
		webDir = "./web"
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	return http.ListenAndServe(":" + port, mux)
}
