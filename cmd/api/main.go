package main

import (
	"context"
	"fmt"
	"log/slog"

	"url_shorter/config"
	"url_shorter/db/pool"
	"url_shorter/internal/handler"
	httpserver "url_shorter/internal/http_server"
	"url_shorter/internal/repository"
	"url_shorter/internal/service"
	"url_shorter/utils/logger"
)

func main() {
	conf, err := config.NewConfig[config.Config]()
	if err != nil {
		panic(fmt.Sprintf("error loading config %v\n", err))
	}

	logger.NewLogger(conf, "api")
	pool, err := pool.NewPool(context.TODO(), conf.Database)
	if err != nil {
		logger.Fatal("error with db", slog.Any("error", err))
	}
	repo := repository.NewLinkRepository(pool)
	serv := service.NewLinkService(repo)
	server := httpserver.NewHttpServer(conf.Server)
	handler.NewHandler(&conf.Server, serv).RegisterRoutes(server.Mux())

	if err := server.Start(); err != nil {
		logger.Fatal("error starting http server", slog.Any("error", err))
	}
}
