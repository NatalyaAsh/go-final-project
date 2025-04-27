package server

import (
	"log/slog"
	"net/http"
	"os"

	"go1f/pkg/api"
	"go1f/pkg/conf"
)

func Start(cfg *conf.Configuration) {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	api.Init(mux, cfg)

	slog.Info("Started", "Port", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, mux)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
