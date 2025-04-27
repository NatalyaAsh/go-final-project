package dbase

import (
	"database/sql"
	"log/slog"
	"os"

	"go1f/pkg/conf"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT "");
		CREATE INDEX IF NOT EXISTS idxDate ON scheduler(date);`
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Init(cfg *conf.Configuration) error {
	_, err := os.Stat(cfg.DBFile)

	var install bool
	if err != nil {
		install = true
	}

	slog.Info("Connect to db", "dbFile", cfg.DBFile)
	db, err = sql.Open("sqlite", cfg.DBFile)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
		slog.Info("Create table and index")
	}
	return nil
}

func CloseDB() {
	slog.Info("Close DB")
	db.Close()
}
