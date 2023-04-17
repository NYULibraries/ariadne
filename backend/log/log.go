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

// Slog levels are internal, so we don't export this const.
// We do export the string for the Level corresponding to this slog level.
const defaultSlogLevel = slog.LevelInfo

// We are wrapping slog calls which take an initial argument that is always output
// to `msg`, which doesn't really fit into our desired scheme for log output
// as we'd like it to appear in Kibana.  For now, we'll just pass `msg` an empty
// string.
const emptyMsg = ""

var DefaultLevelStringOption = getLevelOptionStringForSlogLevel(defaultSlogLevel)

var (
	LevelDebug    = Level(reflect.ValueOf(slog.LevelDebug).Int())
	LevelInfo     = Level(reflect.ValueOf(slog.LevelInfo).Int())
	LevelWarn     = Level(reflect.ValueOf(slog.LevelWarn).Int())
	LevelError    = Level(reflect.ValueOf(slog.LevelError).Int())
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

func Debug(message string, args ...interface{}) {
	slogger.Debug(message, args...)
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
	programLevel.Set(defaultSlogLevel)
	handler := slog.HandlerOptions{Level: programLevel}.NewJSONHandler(os.Stdout)

	return slog.New(handler)
}
