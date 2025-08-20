// Пакет server отвечает за конфигурацию и запуск HTTP-сервера.
package server

import (
	"net/http"

	"yandex-practicum-final-project/pkg/api"
)

// Start инициализирует маршрутизатор (mux), регистрирует все обработчики API
// и запускает прослушивание HTTP-запросов на заданном порту.
func Start() error {
	// Регистрация всех маршрутов API.
	api.Init()

	// Получение порта из переменной окружения или использование значения по умолчанию.
	port := api.Port
	if port == "" {
		port = "7540" // Порт по умолчанию для тестов и локального запуска.
	}

	return http.ListenAndServe(":" + port, nil)
}
