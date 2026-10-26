package main

import (
	"log/slog"
	"os"

	"url_shorter/config"
	"url_shorter/internal/adapters/in/http"
	httpPkg "url_shorter/internal/pkg/http"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(log)

	conf, err := config.NewConfig[config.Config]()
	if err != nil {
		log.Error("error loading config", slog.Any("error", err))
		os.Exit(1)
	}
	server := httpPkg.NewServer(conf.Server)
	handler := http.NewHandler()

	server.Mux().Mount("/", handler.Routes())

	if err := server.Start(); err != nil {
		slog.Error("error starting http server", slog.Any("error", err))
		os.Exit(1)
	}
}
