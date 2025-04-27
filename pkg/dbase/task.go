package dbase

import (
	"database/sql"
	"log/slog"
	"time"
)

const limit = 50

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		slog.Error("Exec", "err", err)
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		slog.Error("Result LastInsertId", "err", err)
		return 0, err
	}
	return id, nil
}

func SelectTask() ([]Task, error) {
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
	ORDER BY date LIMIT :limit`, sql.Named("limit", limit))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	arrTasks := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		arrTasks = append(arrTasks, task)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return []Task{}, err
	}
	return arrTasks, nil
}

func GetTaskTitle(search string) ([]Task, error) {
	searchStr := "%" + search + "%"
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
	WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`,
		sql.Named("search", searchStr), sql.Named("limit", limit))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	arrTasks := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		arrTasks = append(arrTasks, task)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return arrTasks, nil
}

func GetTaskDate(date time.Time) ([]Task, error) {
	dateStr := date.Format("20060102")
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
	WHERE date=:date ORDER BY date LIMIT :limit`, sql.Named("date", dateStr), sql.Named("limit", limit))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	arrTasks := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		arrTasks = append(arrTasks, task)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return arrTasks, nil
}
