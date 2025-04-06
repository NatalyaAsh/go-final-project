package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

const (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128),
		comment TEXT,
		repeat VARCHAR(128)
		);
		CREATE INDEX IF NOT EXISTS idxDate ON scheduler(date);`
)

var (
	db *sql.DB
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Init(dbFileDefault string) error {
	dbFile := getEnv("TODO_DBFILE", dbFileDefault)
	//dbFile := "scheduler.db"
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// настройте подключение к БД

	slog.Info("Connect to db", "dbFile", dbFile)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer db.Close()

	if install {
		res, err := db.Exec(schema)
		if err != nil {
			fmt.Println(res)
			slog.Error(err.Error())
			return err
		}
		slog.Info("Create table and index")
	}
	return nil
}

// store := NewParcelStore(db) // создайте объект ParcelStore функцией NewParcelStore
// service := NewParcelService(store)

// // регистрация посылки
// client := 1
// address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
// p, err := service.Register(client, address)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
