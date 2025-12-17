package api

import (
	"net/http"
	"strconv"

	"go_final_project/pkg/db"
)

// deleteTaskHandler handles DELETE request to remove a task.
// deleteTaskHandler обрабатывает DELETE запрос для удаления задачи.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, map[string]any{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, map[string]any{"error": "id is invalid"})
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]any{"error": "failed to delete task"})
		return
	}

	writeJSON(w, map[string]any{})
}
