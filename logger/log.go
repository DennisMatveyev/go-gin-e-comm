package logger

import (
	"go-gin-e-comm/configs"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

func Setup(cfg *configs.Config) *slog.Logger {
	var log *slog.Logger
	var logWriter io.Writer

	if cfg.Log.Path != "" {
		logWriter = &lumberjack.Logger{
			Filename:   cfg.Log.Path, // Path to the log file
			MaxSize:    10,           // Max size in MB before rotation
			MaxBackups: 5,            // Max number of old log files to keep
			MaxAge:     30,           // Max number of days to retain old log files
			Compress:   true,         // Compress old log files
		}
	} else {
		logWriter = os.Stdout
	}

	switch cfg.Env {
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	case "staging":
		log = slog.New(
			slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(logWriter, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	}
	return log
}
