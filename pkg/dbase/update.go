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

	row := db.QueryRow(`SELECT id, date, title, comment, repeat 
	FROM scheduler WHERE id=:id`, sql.Named("id", id))
	if row == nil {
		slog.Error("Задача не найдена")
		return nil, fmt.Errorf("задача не найдена")
	}

	var task Task
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=:date, title=:title, 
	comment=:comment, repeat=:repeat WHERE id=:id`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	slog.Info("Update", "task", task, "count", count)
	return nil
}
