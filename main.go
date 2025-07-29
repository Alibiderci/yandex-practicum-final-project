package main

import (
	"log"

	"yandex-practicum-final-project/pkg/server"
	"yandex-practicum-final-project/pkg/database"
)

func main() {

	err := database.Init()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.Close()

	err = server.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
}
