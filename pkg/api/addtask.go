package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go1f/pkg/dbase"
)

type ResponseId struct {
	ID int64 `json:"id"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task dbase.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: "ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	ok, msgErr := checkTask(&task)
	if !ok {
		writeJson(w, ResponseErr{Error: msgErr}, http.StatusOK) ////////////
		slog.Error("addTaskHandler:", "checkTask", msgErr)
		return
	}
	id, err := dbase.AddTask(&task)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	slog.Info("", "INSERT ", task, "id", id)
	writeJson(w, ResponseId{ID: id}, http.StatusOK) ///////////
}

func checkTask(task *dbase.Task) (bool, string) {
	if task.Title == "" {
		return false, "task title not specified"
	}

	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	}

	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return false, "wrog Date format"
	}

	ok, msg := checkRepeat(task.Repeat)
	if !ok {
		return false, msg
	}

	if !afterNow(date, time.Now()) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(DateFormat)
		} else {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err == nil {
				task.Date = nextDate
			}
		}
	}
	return true, ""
}

func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	msg, _ := json.Marshal(data)
	io.Writer.Write(w, msg)
}
