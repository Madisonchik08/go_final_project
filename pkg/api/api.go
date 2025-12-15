package api

import "net/http"

// DateLayout is the expected date format for parameters.
const DateLayout = "20060102"

// Init registers all API handlers on provided mux.
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
}
