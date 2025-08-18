package db

import (
	"database/sql"
	"fmt"
	"time"
	"strconv"
)

// Task представляет собой одну задачу в планировщике.
// Структура соответствует полям таблицы 'scheduler' в базе данных.
type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment,omitempty" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// AddTask добавляет новую задачу в базу данных.
// Возвращает ID созданной записи или ошибку.
func AddTask(task *Task) (int64, error) {

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Tasks получает список задач из БД с возможностью поиска.
func Tasks(limit int, search string) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	// Пробуем распознать строку поиска как дату.
	t, err := time.Parse("02.01.2006", search)

	if err == nil {
		// Поиск по дате.
		date := t.Format("20060102")

		query := `SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, date, limit)
	} else {
		// Поиск по ключевому слову или получение всех задач.
		// Создаём запрос к базе данных и список аргументов динамически
		query := `SELECT * FROM scheduler`
		args := []any{}

		// Если ищем задачу по ключевому слову, добавляем нужный отрывок запроса
		if search != "" {
			query += ` WHERE title LIKE ? OR comment LIKE ?`

			// % нужны для поиска подстроки. Они должны быть частью параметра, а не запроса.
			likeSearch := fmt.Sprintf("%%%s%%", search)
			args = append(args, likeSearch, likeSearch)
		}

		// Сортируем по дате вне зависимости, ищем мы по ключевому слову или нет
		query += ` ORDER BY date LIMIT ?`
		args = append(args, limit)

		rows, err = DB.Query(query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Инициализируем пустой срез, чтобы в JSON всегда был массив "tasks": [], а не null.
	tasks := make([]*Task, 0)

	for rows.Next() {
		task := Task{}
		var id int64

		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		task.ID = strconv.FormatInt(id, 10)

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask получает одну задачу из БД по её идентификатору.
func GetTask(id string) (*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	var dbID int64 // Временная переменная для сканирования числового ID из БД.
	task := Task{}

	err := DB.QueryRow(query, id).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Задача с id=%s не найдена", id)
		}
		return nil, err
	}

	task.ID = strconv.FormatInt(dbID, 10)
	
	return &task, nil	
}

// UpdateTask обновляет все поля существующей задачи в БД.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	// Проверяем, была ли действительно обновлена хотя бы одна строка.
	// Если нет, значит задачи с таким ID не существует.
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача с id=%s не найдена", task.ID)
	}

	return nil
}

// DeleteTask удаляет задачу из БД по её идентификатору.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	// Проверка на существование задачи, аналогично UpdateTask.
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача с id=%s не найдена", id)
	}

	return nil
}

// UpdateDate обновляет только дату у существующей задачи.
// Используется при выполнении повторяющейся задачи.
func UpdateDate(next, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача с id=%s не найдена", id)
	}

	return nil
}
