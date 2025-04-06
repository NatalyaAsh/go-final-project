package main

import (
	"go1f/pkg/server"

	"go1f/pkg/db"
)

func main() {
	db.Init("scheduler.db")
	server.Start()

}
