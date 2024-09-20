package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	logPath := os.Getenv("LOG_FILE_PATH")
	if logPath == "" {
		logPath = "log.txt"
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("FATAL ERROR:" + err.Error())

	}

	return slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{}))
}
