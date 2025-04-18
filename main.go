package main

import (
	_ "database/sql"
	"go1f/pkg/dbase"

	"main.go/pkg/server"
	_ "modernc.org/sqlite"
)

func main() {
	err := dbase.Init("scheduler.db")
	if err != nil {
		return
	}
	defer dbase.CloseDB()

	server.Start()

}
