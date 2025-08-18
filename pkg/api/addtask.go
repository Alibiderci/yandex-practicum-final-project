package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"yandex-practicum-final-project/pkg/db"
)

// writeJson — это вспомогательная функция для отправки JSON-ответов.
// Она устанавливает правильный заголовок Content-Type, статус-код
// и кодирует переданные данные в JSON.
func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// checkDate проверяет и корректирует дату задачи в соответствии с бизнес-правилами.
// Функция гарантирует, что у задачи всегда будет актуальная дата.
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(layout)
	}

	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате отличном от 20060102")
	}


	// Проверяем, что дата задачи находится строго в прошлом (т.е. вчера или ранее).
	if !afterNow(t, now) && t.Format(layout) != now.Format(layout) {
		if task.Repeat == "" {
			task.Date = now.Format(layout) 
		} else {
			next, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	} else {
		// Дата сегодня или в будущем. Дату не меняем, но правило повтора нужно проверить.
		if task.Repeat != "" {
			_, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err 

			}
		}

	}
	return nil
}


// AddTaskHandler обрабатывает POST-запросы к эндпоинту /api/task.
// Он декодирует JSON из тела запроса, валидирует данные с помощью checkDate
// и сохраняет новую задачу в базу данных.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {	
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}


	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return

	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "ошибка при добавлении задачи в базу данных"})
		return
	}

	writeJson(w, http.StatusCreated, map[string]int64{"id": id})
}
