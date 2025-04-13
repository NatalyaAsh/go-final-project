package api

import (
	"net/http"
)

const DateFormat = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", taskHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}
