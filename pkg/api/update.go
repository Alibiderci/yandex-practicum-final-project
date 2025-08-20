package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"yandex-practicum-final-project/pkg/db"
)

// updateHandler обрабатывает PUT-запросы к эндпоинту /api/task.
// Он предназначен для полного обновления существующей задачи.
// Ожидает получить в теле запроса JSON со всеми полями задачи, включая ID.
func updateHandler(w http.ResponseWriter, r *http.Request) {
	task := db.Task{}

	// Декодирование и базовая валидация
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Проверяем, что обязательные поля не пустые.
	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "идентификатор не может быть пустым"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// Комплексная валидация даты и правила повторения
	// Используем ту же самую функцию валидации, что и при создании задачи
	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return

	}

	// Вызов бизнес-логики (обновление задачи в БД)
	// После того как все данные проверены, пытаемся обновить запись в базе данных.
	err = db.UpdateTask(&task)
	// Обрабатываем ошибки от слоя БД.
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		} else {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
