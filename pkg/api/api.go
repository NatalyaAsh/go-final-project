package api

import (
	"net/http"
)

const DateFormat = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", auth(addTaskHandler))
	mux.HandleFunc("GET /api/tasks", auth(TasksHandler))
	mux.HandleFunc("GET /api/task", auth(getTaskHandler))
	mux.HandleFunc("PUT /api/task", auth(updateTaskHandler))
	mux.HandleFunc("POST /api/task/done", auth(doneTaskHandler))
	mux.HandleFunc("DELETE /api/task", auth(deleteTaskHandler))
	mux.HandleFunc("POST /api/signin", signinHandler)
}
