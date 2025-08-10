package server

import (
	"net/http"

	"yandex-practicum-final-project/pkg/api"
)

func Start() error {
	api.Init()

	port := api.Port
	if port == "" {
		port = "7540"
	}

	return http.ListenAndServe(":" + port, nil)
}
