package main

import (
	"context"
	"log/slog"
	"os"

	"url_shorter/config"
	"url_shorter/db/pool"
	"url_shorter/internal/handler"
	httpserver "url_shorter/internal/http_server"
	"url_shorter/internal/repository"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(log)

	conf, err := config.NewConfig[config.Config]()
	if err != nil {
		log.Error("error loading config", slog.Any("error", err))
		os.Exit(1)
	}
	pool, err := pool.NewPool(context.TODO(), conf.Database)
	if err != nil {
		log.Error("error with db", slog.Any("error", err))
		os.Exit(1)
	}
	repo := repository.NewUrlCodeRepository(pool)
	server := httpserver.NewHttpServer(conf.Server)
	handler := handler.NewHandler(repo)
	server.Mux().Mount("/", handler.Routes())

	if err := server.Start(); err != nil {
		slog.Error("error starting http server", slog.Any("error", err))
		os.Exit(1)
	}
}
