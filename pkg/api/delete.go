package api

import (
	"net/http"
	"strings"
	"yandex-practicum-final-project/pkg/db"
)

// deleteHandler обрабатывает DELETE-запросы к эндпоинту /api/task.
// Он предназначен для удаления задачи по её идентификатору.
// Идентификатор задачи должен быть передан как GET-параметр 'id' в URL.
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение параметра 'id' из URL запроса (например, из /api/task?id=123).
	id := r.URL.Query().Get("id")
	// Проверяем, был ли вообще предоставлен идентификатор.
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Вызываем функцию из пакета db для непосредственного удаления задачи из базы данных.
	err := db.DeleteTask(id)
	if err != nil {
		// Проверяем, является ли ошибка ошибкой "задача не найдена".
		// Это позволяет нам вернуть более точный и семантически верный HTTP-статус.
		if strings.Contains(err.Error(), "не найдена") {
			writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})	
		} else {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})	
		}	
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
