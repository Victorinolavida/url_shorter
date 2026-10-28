package logger

import (
	"log/slog"
	"os"
	"strings"

	"url_shorter/config"
)

func parseLogLevel(level string) slog.Leveler {
	levelClean := strings.ToLower(strings.TrimSpace(level))
	switch levelClean {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func NewLogger(conf *config.Config, appName string) *slog.Logger {
	isProd := conf.Env == "production"

	var handler slog.Handler
	level := parseLogLevel(conf.LogLevel)
	if isProd {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	}

	log := slog.New(handler)
	if appName == "" {
		log.Warn("service name is empty")
	} else {
		log = log.With(slog.String("service", appName))
	}
	slog.SetDefault(log)
	return log
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
