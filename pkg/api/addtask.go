package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
)

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

	// if task date is before today, adjust based on repeat
	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(DateLayout)
		} else {
			task.Date = next
		}
	}
	return nil
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
