package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
)

// doneTaskHandler handles POST request to mark a task as done.
// If task has no repeat, it deletes the task.
// If task has repeat, it updates the date to the next occurrence.
// doneTaskHandler обрабатывает POST запрос для отметки задачи как выполненной.
// Если у задачи нет повторения, она удаляется.
// Если у задачи есть повторение, дата обновляется на следующее вхождение.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, map[string]any{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, map[string]any{"error": "invalid id"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, map[string]any{"error": "task not found"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]any{"error": "database error"}, http.StatusInternalServerError)
		}
		return
	}

	date := task.Date
	repeat := task.Repeat

	if repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]any{"error": "failed to delete task"}, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := nextdate.NextDate(now, date, repeat)
		if err != nil {
			writeJSON(w, map[string]any{"error": "failed to calculate next date"}, http.StatusBadRequest)
			return
		}
		err = db.UpdateDate(id, nextDate)
		if err != nil {
			writeJSON(w, map[string]any{"error": "failed to update task date"}, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]any{}, http.StatusOK)
}
