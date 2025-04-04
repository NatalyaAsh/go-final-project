package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// func getRoot(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r)
// 	io.WriteString(w, "root request\n")
// }

func getHello(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("hello, %s\n", r.PathValue("name"))
	io.WriteString(w, msg)
}

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
	mux.HandleFunc("GET /hello/{name}", getHello)

	slog.Info("Started", "Port", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
