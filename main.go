// Пакет main является точкой входа в приложение.
// Он отвечает за инициализацию зависимостей и запуск сервера.
package main

import (
	"log"

	"yandex-practicum-final-project/pkg/server"
	"yandex-practicum-final-project/pkg/db"
)

// main - главная функция программы
// Она последовательно инициализирует соединение с базой данных,
// а затем запускает веб-сервер
func main() {
	// Инициализация соединения с базой данных.
	err := db.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Отложенное закрытие соединения с БД при завершении работы программы.
	defer db.Close()

	// Запуск веб-сервера.
	err = server.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
}
