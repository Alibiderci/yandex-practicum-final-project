package db

import (
	"database/sql"
	"fmt"
	"time"
	"strconv"
)

type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment,omitempty" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

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

func Tasks(limit int, search string) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	t, err := time.Parse("02.01.2006", search)

	if err == nil {
		date := t.Format("20060102")

		query := `SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, date, limit)
	} else {
		query := `SELECT * FROM scheduler`
		args := []any{}

		if search != "" {
			query += ` WHERE title LIKE ? OR comment LIKE ?`

			likeSearch := fmt.Sprintf("%%%s%%", search)
			args = append(args, likeSearch, likeSearch)
		}

		query += ` ORDER BY date LIMIT ?`
		args = append(args, limit)

		rows, err = DB.Query(query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func GetTask(id string) (*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	var dbID int64
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

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача с id=%s не найдена", task.ID)
	}

	return nil
}
