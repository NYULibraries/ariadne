package log

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"math"
	"os"
	"reflect"
	"strings"
)

type Level int

const defaultLevel = slog.LevelInfo

var DefaultLevelStringOption = getLevelOptionStringForSlogLevel(defaultLevel)

var (
	LevelDebug    = Level(reflect.ValueOf(slog.LevelDebug).Int())
	LevelInfo     = Level(reflect.ValueOf(slog.LevelInfo).Int())
	LevelWarn     = Level(reflect.ValueOf(slog.LevelWarn).Int())
	LevelError    = Level(reflect.ValueOf(slog.LevelWarn).Int())
	LevelDisabled = Level(math.MaxInt)
)

var logLevelStringOptions = map[string]Level{
	"debug":    LevelDebug,
	"info":     LevelInfo,
	"warn":     LevelWarn,
	"error":    LevelError,
	"disabled": LevelDisabled,
}

var programLevel = new(slog.LevelVar)
var slogger *slog.Logger

func init() {
	slogger = newDefaultSlogger()
}

func Error(message string, args ...interface{}) {
	slogger.Error(message, args...)
}

func Fatal(args ...interface{}) {
	Error(fmt.Sprint(args...))
	os.Exit(1)
}

func GetValidLevelOptionStrings() []string {
	// Note that we are returning the strings in a specific order: i.e. in order of
	// increasing severity.
	return []string{
		"debug",
		"info",
		"warn",
		"error",
		"disabled",
	}
}

func Info(message string, args ...interface{}) {
	slogger.Info(message, args...)
}

func Warn(message string, args ...interface{}) {
	slogger.Warn(message, args...)
}

func SetLevel(level Level) {
	programLevel.Set(slog.Level(level))
}

func SetLevelByString(levelStringArg string) error {
	level, ok := logLevelStringOptions[levelStringArg]
	if ok {
		programLevel.Set(slog.Level(level))
		return nil
	} else {
		return errors.New(fmt.Sprintf("\"%s\" is not a valid error string option.  Valid options: %s",
			levelStringArg, strings.Join(GetValidLevelOptionStrings(), ", ")))
	}
}

func SetOutput(bytesBuffer *bytes.Buffer) {
	logWriter := bufio.NewWriter(bytesBuffer)
	handler := slog.HandlerOptions{Level: programLevel}.NewJSONHandler(logWriter)
	slog.SetDefault(slog.New(handler))
	slogger = slog.New(handler)
}

func getLevelOptionStringForSlogLevel(levelArg slog.Level) string {
	for optionString, levelValue := range logLevelStringOptions {
		if levelValue == Level(reflect.ValueOf(levelArg).Int()) {
			return optionString
		}
	}

	return ""
}

func newDefaultSlogger() *slog.Logger {
	programLevel.Set(defaultLevel)
	handler := slog.HandlerOptions{Level: programLevel}.NewJSONHandler(os.Stdout)
	slog.SetDefault(slog.New(handler))

	return slog.New(handler)
}
