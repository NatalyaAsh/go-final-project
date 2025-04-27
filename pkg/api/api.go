package api

import (
	"net/http"

	"go1f/pkg/conf"
)

const DateFormat = "20060102"

func Init(mux *http.ServeMux, cfg *conf.Configuration) {
	mux.HandleFunc("/api/nextdate", auth(nextDayHandler, cfg))
	mux.HandleFunc("POST /api/task", auth(addTaskHandler, cfg))
	mux.HandleFunc("GET /api/tasks", auth(TasksHandler, cfg))
	mux.HandleFunc("GET /api/task", auth(getTaskHandler, cfg))
	mux.HandleFunc("PUT /api/task", auth(updateTaskHandler, cfg))
	mux.HandleFunc("POST /api/task/done", auth(doneTaskHandler, cfg))
	mux.HandleFunc("DELETE /api/task", auth(deleteTaskHandler, cfg))
	mux.HandleFunc("POST /api/signin", signinWrapper(signinHandler, cfg))
}
