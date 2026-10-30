package main

import (
	"context"
	"fmt"
	"log/slog"

	"url_shorter/cache/redis"
	"url_shorter/config"
	"url_shorter/db/pool"
	"url_shorter/internal/entity"
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

	cache, err := redis.NewRedisClient(conf.Redis)
	if err != nil {
		logger.Fatal("error with redis", slog.Any("error", err))
	}

	serv := service.NewLinkService(repo, cache)
	server := httpserver.NewHttpServer(conf.Server)
	ipHasher := entity.NewIPHasher([]byte(conf.IPHashKey))
	handler.NewHandler(&conf.Server, serv, ipHasher).RegisterRoutes(server.Mux())

	if err := server.Start(); err != nil {
		logger.Fatal("error starting http server", slog.Any("error", err))
	}
}
