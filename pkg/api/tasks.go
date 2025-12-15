package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJSON(w, tasksResp{Tasks: tasks})
}
