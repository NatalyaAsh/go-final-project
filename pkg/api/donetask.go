package api

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"go1f/pkg/dbase"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Error("doneTaskHandler: Не указан идентификатор")
		//http.Error(w, "не указан идентификатор", http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: "не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := dbase.GetTask(id)
	if err != nil {
		slog.Error("doneTaskHandler:", "", err.Error())
		//http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		err = dbase.DeleteTask(id)
		if err != nil {
			slog.Error("doneTaskHandler:", "", err.Error())
			//http.Error(w, err.Error(), http.StatusBadRequest)
			writeJson(w, ResponseErr{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		io.Writer.Write(w, []byte("{}"))
		return
	}

	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		slog.Error("doneTaskHandler:", "", err.Error())
		//http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	err = dbase.UpdateTask(task)
	if err != nil {
		slog.Error("doneTaskHandler:", "", err.Error())
		//http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	io.Writer.Write(w, []byte("{}"))
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Error("deleteTaskHandler: Не указан идентификатор")
		writeJson(w, ResponseErr{Error: "не указан идентификатор"}, http.StatusOK)
		return
	}
	err := dbase.DeleteTask(id)
	if err != nil {
		slog.Error("deleteTaskHandler:", "", err.Error())
		writeJson(w, ResponseErr{Error: err.Error()}, http.StatusOK)
		return
	}
	io.Writer.Write(w, []byte("{}"))
}
