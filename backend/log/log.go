package log

import (
	"fmt"
	"golang.org/x/exp/slog"
	"os"
)

var slogger *slog.Logger

func init() {
	handler := slog.HandlerOptions{Level: slog.LevelInfo}.NewJSONHandler(os.Stdout)
	slog.SetDefault(slog.New(handler))

	slogger = slog.New(handler)
}

func Fatal(args ...interface{}) {
	Error(fmt.Sprint(args...))
	os.Exit(1)
}

func Error(message string, args ...interface{}) {
	slogger.Error(message, args...)
}

func Info(message string, args ...interface{}) {
	slogger.Info(message, args...)
}
