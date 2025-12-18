package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

// editTaskHandler handles PUT request to update an existing task.
// editTaskHandler обрабатывает PUT запрос для обновления существующей задачи.
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	task.Comment = strings.TrimSpace(task.Comment)
	task.Repeat = strings.TrimSpace(task.Repeat)

	if task.Title == "" {
		writeJSON(w, map[string]any{"error": "title is required"}, http.StatusBadRequest)
		return
	}

	now := normalizeDate(time.Now())

	if err := checkDate(&task, now); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if task.ID == 0 {
		writeJSON(w, map[string]any{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]any{}, http.StatusOK)
}
