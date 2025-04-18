package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"go1f/pkg/dbase"
)

type TasksResp struct {
	Tasks []dbase.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var tasks []dbase.Task
	var err error
	if search == "" {
		tasks, err = dbase.SelectTask()
	} else {
		date, notDate := time.Parse("02.01.2006", search)
		if notDate != nil {
			slog.Info("TasksHandler: GetTaskTitle", "search", search)
			tasks, err = dbase.GetTaskTitle(search)
		} else {
			slog.Info("TasksHandler: GetTaskDate", "date", date)
			tasks, err = dbase.GetTaskDate(date)
		}
	}
	if err != nil {
		slog.Error("TasksHandler:", "", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: err.Error()})
		return
	}
	writeJson(w, TasksResp{Tasks: tasks})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Error("getTaskHandler: Не указан идентификатор")
		writeJson(w, ResponseErr{Error: "не указан идентификатор"})
		return
	}
	task, err := dbase.GetTask(id)
	if err != nil {
		slog.Error("getTaskHandler:", "", err.Error())
		writeJson(w, ResponseErr{Error: err.Error()})
		return
	}
	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task dbase.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: err.Error()})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: "ошибка десериализации JSON"})
		return
	}
	ok, msgErr := checkTask(&task)
	if !ok {
		writeJson(w, ResponseErr{Error: msgErr})
		return
	}

	if task.ID == "" {
		slog.Error("updateTaskHandler: Не указан идентификатор")
		writeJson(w, ResponseErr{Error: "не указан идентификатор"})
		return
	}
	err = dbase.UpdateTask(&task)
	if err != nil {
		slog.Error("updateTaskHandler", "", err.Error())
		writeJson(w, ResponseErr{Error: err.Error()})
		return
	}
	writeJson(w, task)
}
