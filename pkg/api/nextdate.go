package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/nextdate"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			writeJSON(w, map[string]any{"error": "invalid now parameter"})
			return
		}
	}

	res, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	_, _ = w.Write([]byte(res))
}
