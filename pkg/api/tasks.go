package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

// tasksResp represents response with list of tasks.
// tasksResp представляет ответ со списком задач.
type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler handles GET request to retrieve list of tasks.
// tasksHandler обрабатывает GET запрос для получения списка задач.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJSON(w, tasksResp{Tasks: tasks}, http.StatusOK)
}
