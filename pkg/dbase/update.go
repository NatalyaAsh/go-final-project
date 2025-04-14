package dbase

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
)

func GetTask(id string) (*Task, error) {
	_, err := strconv.Atoi(id)
	if err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf("неверно указан идентификатор")
	}

	rows, err := db.Query(`SELECT id, date, title, comment, repeat 
	FROM scheduler WHERE id=:id`, sql.Named("id", id))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var task Task
	err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &task, nil
}
