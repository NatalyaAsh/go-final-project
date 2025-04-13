package dbase

import (
	"database/sql"
	"fmt"
	"log/slog"
)

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

func SelectTask() {
	slog.Info("Select * from scheduler:")

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var id, date, title, comment, repeat string
	for rows.Next() {
		err := rows.Scan(&id, &date, &title, &comment, &repeat)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(id, date, title, comment, repeat)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}
}
