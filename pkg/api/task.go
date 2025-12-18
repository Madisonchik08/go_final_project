package api

import (
	"encoding/json"
	"net/http"
)

// taskHandler routes task requests based on HTTP method.
// taskHandler маршрутизирует запросы к задачам в зависимости от HTTP метода.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// writeJSON writes JSON response to the client with specified status code.
// writeJSON записывает JSON ответ клиенту с указанным кодом статуса.
func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log encoding error if it occurs
		// Логируем ошибку кодирования, если она возникает
	}
}
