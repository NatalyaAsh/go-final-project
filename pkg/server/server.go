package server

import (
	"log/slog"
	"net/http"
	"os"

	"main.go/pkg/api"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Start() {
	port := getEnv("TODO_PORT", "7540")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	api.Init(mux)

	slog.Info("Started", "Port", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
