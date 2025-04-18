package api

import (
	"net/http"
)

const DateFormat = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", addTaskHandler)
	mux.HandleFunc("GET /api/tasks", TasksHandler)
	mux.HandleFunc("GET /api/task", getTaskHandler)
	mux.HandleFunc("PUT /api/task", updateTaskHandler)
	mux.HandleFunc("POST /api/task/done", doneTaskHandler)
	mux.HandleFunc("DELETE /api/task", deleteTaskHandler)
}

// func taskHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodPost:
// 		addTaskHandler(w, r)
// 	case http.MethodGet:
// 		getTasksHandler(w, r)
// 	}
// }
