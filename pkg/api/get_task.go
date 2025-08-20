package api

import (
	"net/http"
	"strings"
	"yandex-practicum-final-project/pkg/db"
)

// getTaskHandler обрабатывает GET-запросы к эндпоинту /api/task.
// Его задача — вернуть полную информацию о одной задаче по её ID.
// Идентификатор задачи должен быть передан как GET-параметр 'id'.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получение и валидация ID задачи
	id := r.URL.Query().Get("id")
	if id == "" {
		// Если ID отсутствует, запрос некорректен.
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})	
		return
	}

	// Вызов бизнес-логики (получение задачи из БД)
	task, err := db.GetTask(id)
	if err != nil {
		// Обрабатываем возможные ошибки от слоя базы данных.
		if strings.Contains(err.Error(), "не найдена") {
			writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		} else {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return
	}
	
	writeJson(w, http.StatusOK, task)	
}
