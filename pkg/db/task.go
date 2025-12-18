package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Task represents a scheduler record.
// Task представляет запись планировщика.
type Task struct {
	ID      int64  `db:"id" json:"id,string"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

// AddTask inserts a new task into scheduler table and returns its ID.
// AddTask вставляет новую задачу в таблицу планировщика и возвращает её ID.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get last inserted id: %w", err)
	}
	return id, nil
}

// Tasks returns up to limit tasks ordered by date ascending.
// If search is provided, filters by date (dd.mm.yyyy) or substring in title/comment.
// Tasks возвращает до limit задач, отсортированных по дате по возрастанию.
// Если указан search, фильтрует по дате (dd.mm.yyyy) или подстроке в title/comment.
func Tasks(limit int, search string) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}

	search = strings.TrimSpace(search)

	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`, limit)
	} else {
		if d, derr := time.Parse("02.01.2006", search); derr == nil {
			searchDate := d.Format("20060102")
			rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`, searchDate, limit)
		} else {
			pattern := "%" + search + "%"
			rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`, pattern, pattern, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

// GetTask returns task by id.
// GetTask возвращает задачу по id.
func GetTask(id int) (*Task, error) {
	var t Task
	err := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).
		Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTask updates task fields by id.
// UpdateTask обновляет поля задачи по id.
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
		return fmt.Errorf("task not found")
	}
	return nil
}
