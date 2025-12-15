package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]any{"error": "id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, task)
}
