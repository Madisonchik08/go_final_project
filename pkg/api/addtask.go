package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
)

// addTaskHandler handles POST request to create a new task.
// addTaskHandler обрабатывает POST запрос для создания новой задачи.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	task.Comment = strings.TrimSpace(task.Comment)
	task.Repeat = strings.TrimSpace(task.Repeat)

	if task.Title == "" {
		writeJSON(w, map[string]any{"error": "title is required"})
		return
	}

	now := normalizeDate(time.Now())
	if err := checkDate(&task, now); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]any{"id": strconv.FormatInt(id, 10)})
}

// checkDate validates and adjusts task date based on repeat rules.
// checkDate проверяет и корректирует дату задачи на основе правил повторения.
func checkDate(task *db.Task, now time.Time) error {
	if task.Date == "" {
		task.Date = now.Format(DateLayout)
	}

	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return err
	}
	t = normalizeDate(t)

	var next string
	if task.Repeat != "" {
		next, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// If task date is before today, adjust based on repeat
	// Если дата задачи раньше сегодняшней, корректируем на основе повторения
	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(DateLayout)
		} else {
			task.Date = next
		}
	}
	return nil
}

// normalizeDate removes time component, keeping only date.
// normalizeDate удаляет компонент времени, оставляя только дату.
func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
