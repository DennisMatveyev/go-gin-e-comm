package logger

import (
	"go-gin-e-comm/configs"
	"log/slog"

	"github.com/natefinch/lumberjack"
)

func InitSlogLogger(cfg *configs.Config) {
	logFile := &lumberjack.Logger{
		Filename:   cfg.Log.Path,
		MaxSize:    cfg.Log.MaxSize,    // Max size in MB before rotation
		MaxBackups: cfg.Log.MaxBackups, // Keep last 5 logs
		MaxAge:     cfg.Log.MaxAge,     // Keep logs for 30 days
		Compress:   cfg.Log.Compress,   // Compress old logs
	}
	// multiWriter := io.MultiWriter(logFile, os.Stderr)
	log := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(log)
}
