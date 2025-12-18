package api

import (
	"net/http"
	"strconv"

	"go_final_project/pkg/db"
)

// getTaskHandler handles GET request to retrieve a task by ID.
// getTaskHandler обрабатывает GET запрос для получения задачи по ID.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusNotFound)
		return
	}
	writeJSON(w, task, http.StatusOK)
}
