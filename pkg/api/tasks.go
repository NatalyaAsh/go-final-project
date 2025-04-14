package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
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
	// Поиск
	// search := r.URL.Query().Get("search")
	// if search != "" {
	// 	date, err := time.Parse("02.01.2006", search)

	// 	var tasks []dbase.Task
	// 	if err != nil {
	// 		slog.Info("GetTaskTitle")
	// 		tasks, err = dbase.GetTaskTitle(search)
	// 	} else {
	// 		slog.Info("GetTaskDate")
	// 		tasks, err = dbase.GetTaskDate(date)
	// 	}

	// 	if err != nil {
	// 		slog.Error(err.Error())
	// 		var resErr ResponseErr
	// 		resErr.Error = err.Error()
	// 		msg, err := json.Marshal(resErr)
	// 		if err != nil {
	// 			return
	// 		}
	// 		io.Writer.Write(w, msg)
	// 		return
	// 	}
	// 	slog.Info("Select", "tasks", tasks)
	// 	var tasksResp TasksResp
	// 	tasksResp.Tasks = tasks
	// 	msg, err := json.Marshal(tasksResp)
	// 	if err != nil {
	// 		return
	// 	}
	// 	io.Writer.Write(w, msg)
	// 	return
	// }
	// По индексу
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error(err.Error())
			var resErr ResponseErr
			resErr.Error = "Неверно указан идентификатор"
			//			resErr.Error = err.Error()
			msg, err := json.Marshal(resErr)
			if err != nil {
				return
			}
			io.Writer.Write(w, msg)
			return
		}
		task, err := dbase.GetTaskOnId(id)
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

	}

}
