package main

import (
	"log"

	"yandex-practicum-final-project/pkg/server"
	"yandex-practicum-final-project/pkg/db"
)

func main() {

	err := db.Init()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	err = server.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
}
