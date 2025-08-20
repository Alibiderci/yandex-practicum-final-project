package api

import (
	"net/http"
	"strings"
	"time"
	"yandex-practicum-final-project/pkg/db"
)

// doneHandler обрабатывает POST-запросы к эндпоинту /api/task/done.
// Этот хендлер отмечает задачу как "выполненную".
// Логика зависит от типа задачи:
// - Одноразовые задачи (без правила `repeat`) удаляются.
// - Повторяющиеся задачи получают новую, следующую дату выполнения.
func doneHandler(w http.ResponseWriter, r *http.Request) {
	// Получение и валидация ID задачи
	id := r.URL.Query().Get("id")	
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})	
		return
	}

	// Получение данных о задаче 
	// Нам нужна полная информация о задаче, чтобы понять, есть ли у нее правило повторения.
	task, err := db.GetTask(id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		} else {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return
	}

	// Основная бизнес-логика: удалить или обновить
	if task.Repeat == "" {
		// Сценарий A: Одноразовая задача
		// Если правила повторения нет, просто удаляем задачу.
		err := db.DeleteTask(id)	

		if err != nil {
			if strings.Contains(err.Error(), "не найдена") {
				writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			} else {
				writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			return
		}
	} else {
		// Сценарий B: Повторяющаяся задача
		now := time.Now()

		// Вычисляем следующую дату выполнения, используя нашу основную функцию.
		next, err := nextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// Обновляем дату задачи в базе данных на вычисленную следующую дату.
		err = db.UpdateDate(next, id)
		if err != nil {
			if strings.Contains(err.Error(), "не найдена") {
				writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			} else {
				writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			return
		}
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
