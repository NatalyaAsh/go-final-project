package api

import (
	"net/http"
)

const DateFormat = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", addTaskHandler)
	mux.HandleFunc("GET /api/task", getTaskHandler)
	mux.HandleFunc("GET /api/tasks", TasksHandler)
}

// func taskHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodPost:
// 		addTaskHandler(w, r)
// 	case http.MethodGet:
// 		getTasksHandler(w, r)
// 	}
// }
