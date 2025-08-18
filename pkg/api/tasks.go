package api

import (
	"net/http"
	"yandex-practicum-final-project/pkg/db"
)

// tasksResponse — это структура-обертка для JSON-ответа со списком задач.
// Использование такой структуры делает JSON более читаемым и расширяемым в будущем.
type TaskResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы к эндпоинту /api/tasks.
// Он возвращает список всех задач, отсортированных по дате.
// Также поддерживает опциональный GET-параметр 'search' для фильтрации
// задач по дате или ключевому слову.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра для поиска
	// Извлекаем значение параметра 'search'. Если его нет, вернется пустая строка "".
	searchQuery := r.URL.Query().Get("search")

	// Вызов бизнес-логики (получение задачи из БД) 
	// Вызываем функцию из пакета db, передавая ей лимит и строку поиска.
	// Слой БД сам разберется, как выполнить запрос (поиск по дате, по слову или все задачи).
	tasks, err := db.Tasks(50, searchQuery)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Создаем объект ответа, помещая в него полученный срез задач.
	// (Функция db.Tasks уже гарантирует, что tasks не будет nil, а будет пустым срезом).
	response := TaskResponse{Tasks: tasks}
	writeJson(w, http.StatusOK, response)
}
