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
	var responseId ResponseId
	var responseErr ResponseErr

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		responseErr.Error = "ошибка десериализации JSON"
		msg, _ := json.Marshal(responseErr)
		io.Writer.Write(w, msg)
		return
	}

	slog.Info("Recive: ", "task", task)
	ok, msgErr := checkTask(&task)
	if !ok {
		responseErr.Error = msgErr
		msg, _ := json.Marshal(responseErr)
		io.Writer.Write(w, msg)
		slog.Error(responseErr.Error)
		return

	}
	id, err := dbase.AddTask(&task)

	if err != nil {
		responseErr.Error = "task title not specified"
		slog.Error("Result query", "err", err)
	} else {
		responseId.ID = id
		slog.Info("Result query", "id", id)
	}
	msg, _ := json.Marshal(responseId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	io.Writer.Write(w, msg)
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
		slog.Error("Wrog Date format", "Было Date", task.Date)
		return false, "wrog Date format"
	}

	ok, msg := checkRepeat(task.Repeat)
	if !ok {
		return false, msg
	}

	if !afterNow(date, time.Now()) {
		slog.Info("Date < Now ->")
		if task.Repeat == "" {
			task.Date = time.Now().Format(DateFormat)
			slog.Info("Repeat Empty", "Date", task.Date)
		} else {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err == nil {
				task.Date = nextDate
				slog.Info("Repeat Exist", "nextDate", task.Date)
			}
		}
	}

	return true, ""
}
