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
			slog.Info("GetTaskTitle", "search", search)
			tasks, err = dbase.GetTaskTitle(search)
		} else {
			slog.Info("GetTaskDate", "date", date)
			tasks, err = dbase.GetTaskDate(date)
		}
	}
	if err != nil {
		slog.Error(err.Error())
		var resErr ResponseErr
		resErr.Error = err.Error()
		msg, err := json.Marshal(resErr)
		if err != nil {
			return
		}
		io.Writer.Write(w, msg)
		return
	}
	//slog.Info("Select", "tasks", tasks)
	var tasksResp TasksResp
	tasksResp.Tasks = tasks
	msg, err := json.Marshal(tasksResp)
	if err != nil {
		return
	}
	io.Writer.Write(w, msg)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		task, err := dbase.GetTask(id)
		if err != nil {
			slog.Error(err.Error())
			var resErr ResponseErr
			resErr.Error = err.Error()
			msg, err := json.Marshal(resErr)
			if err != nil {
				return
			}
			io.Writer.Write(w, msg)
			return
		}
		msg, err := json.Marshal(task)
		if err != nil {
			return
		}
		io.Writer.Write(w, msg)
		return

	} else {
		slog.Error("Не указан идентификатор")
		var resErr ResponseErr
		resErr.Error = "не указан идентификатор"
		msg, err := json.Marshal(resErr)
		if err != nil {
			return
		}
		io.Writer.Write(w, msg)
		return
	}

}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task dbase.Task
	//var responseId ResponseId
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
	ok, msgErr := checkTask(&task)
	if !ok {
		responseErr.Error = msgErr
		msg, _ := json.Marshal(responseErr)
		io.Writer.Write(w, msg)
		slog.Error(responseErr.Error)
		return
	}

	if task.ID != "" {
		err := dbase.UpdateTask(&task)
		if err != nil {
			slog.Error(err.Error())
			var resErr ResponseErr
			resErr.Error = err.Error()
			msg, err := json.Marshal(resErr)
			if err != nil {
				return
			}
			io.Writer.Write(w, msg)
			return
		}
		msg, err := json.Marshal(task)
		if err != nil {
			return
		}
		io.Writer.Write(w, msg)
		return

	} else {
		slog.Error("Не указан идентификатор")
		var resErr ResponseErr
		resErr.Error = "не указан идентификатор"
		msg, err := json.Marshal(resErr)
		if err != nil {
			return
		}
		io.Writer.Write(w, msg)
		return
	}

}
