// Пакет api содержит все обработчики HTTP-запросов (хендлеры) для веб-сервера.
// Он отвечает за маршрутизацию, прием запросов, валидацию данных и отправку
// JSON-ответов клиенту.
package api

import (
	"net/http"
	"os"
)

// Port хранит порт, который будет прослушивать сервер.
// Значение берется из переменной окружения TODO_PORT.
var Port string = os.Getenv("TODO_PORT")

// webDir хранит путь к каталогу со статическими файлами фронтенда.
// Значение берется из переменной окружения TODO_WEBDIR.
var webDir string = os.Getenv("TODO_WEBDIR")

// taskDispatcher является диспетчером для маршрута /api/task.
// Он анализирует HTTP-метод запроса (GET, POST, PUT, DELETE) и вызывает
// соответствующий, более специфичный обработчик.
func taskDispatcher(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)	
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	default: 
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Init регистрирует все маршруты API в предоставленном ServeMux.
// Эта функция должна вызываться один раз при старте сервера.
func Init() {

	if webDir == "" {
		webDir = "web" // Путь к файлам фронтенда по умолчанию.
	}

	// Обработчик для статических файлов (HTML, CSS, JS).
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Обработчик для входа в программу
	http.HandleFunc("/api/signin", loginHandler)

	// Регистрация всех API-эндпоинтов.
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(taskDispatcher)) // Диспетчер для одиночных задач.
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
}

