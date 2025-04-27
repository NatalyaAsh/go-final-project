package main

import (
	_ "database/sql"
	"go1f/pkg/conf"
	"go1f/pkg/dbase"
	"go1f/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := conf.New()

	err := dbase.Init(cfg)
	if err != nil {
		return
	}
	defer dbase.CloseDB()

	server.Start(cfg)

}
