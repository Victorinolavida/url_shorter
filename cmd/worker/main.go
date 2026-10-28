package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"url_shorter/config"
	"url_shorter/db/pool"
	"url_shorter/internal/repository"
	"url_shorter/internal/service"
	"url_shorter/utils/logger"
)

var defaultIntervalDuration = time.Second

func main() {
	conf, err := config.NewConfig[config.Config]()
	if err != nil {
		panic(fmt.Sprintf("error loading config %v\n", err))
	}

	log := logger.NewLogger(conf, "worker")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pool.NewPool(ctx, conf.Database)
	if err != nil {
		logger.Fatal("error with db", slog.Any("error", err))
	}
	repo := repository.NewLinkRepository(pool)
	serv := service.NewLinkService(repo)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Go(func() {
		defer wg.Done()
		runWorker(ctx, conf.WorkerInterval, serv)
	})

	<-ctx.Done()
	log.Info("shutdown signal received")
	wg.Wait() // let in-flight work finish
	log.Info("worker stopped")
}

func runWorker(ctx context.Context, duration time.Duration, service service.LinkService) {
	interval, err := time.ParseDuration(duration.String())
	if err != nil {
		slog.Error("error parsing duration", slog.Any("error", err))
	}
	if interval <= 0 {
		slog.Warn("duration is less than 0, using the default instead")
		duration = defaultIntervalDuration
	}

	tick := time.NewTicker(duration)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			doWork(ctx, service)
		}
	}
}

func doWork(ctx context.Context, service service.LinkService) {
	slog.Info("trying to expire links")
	err := service.ExpireActiveLinks(ctx)
	if err != nil {
		slog.Error("error expiring links", slog.Any("error", err))
		return
	}
}
