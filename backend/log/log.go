package log

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"io"
	"math"
	"os"
	"reflect"
	"sort"
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

func Debug(args ...any) {
	slogger.Debug(emptyMsg, args...)
}

func Error(args ...any) {
	slogger.Error(emptyMsg, args...)
}

func Fatal(args ...any) {
	Error(fmt.Sprint(args...))
	os.Exit(1)
}

func GetLevelOptionStringForLogLevel(levelArg Level) string {
	for optionString, levelValue := range logLevelStringOptions {
		if levelValue == levelArg {
			return optionString
		}
	}

	return ""
}

// Returns a slice of level option strings in increasing order of severity
// for use in help messages.
// Uses code based on https://stackoverflow.com/questions/18695346/how-can-i-sort-a-mapstringint-by-its-values.
func GetValidLevelOptionStrings() []string {
	var orderedLevelOptionStrings []string

	type pair struct {
		levelString string
		levelValue  Level
	}

	var temp []pair
	for levelString, levelValue := range logLevelStringOptions {
		temp = append(temp, pair{levelString, levelValue})
	}

	sort.SliceStable(temp, func(i, j int) bool {
		return temp[i].levelValue < temp[j].levelValue
	})

	for _, pair := range temp {
		orderedLevelOptionStrings = append(orderedLevelOptionStrings, pair.levelString)
	}

	return orderedLevelOptionStrings
}

func Info(args ...any) {
	slogger.Info(emptyMsg, args...)
}

func SetLevel(level Level) {
	programLevel.Set(slog.Level(level))
}

// Useful for setting the level based on a string value set by a user via a flag
// in CLI mode.
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

// Redirect output to another stream besides stdout, or to a bytes.Buffer for
// testing.
func SetOutput(logWriter io.Writer) {
	handler := slog.HandlerOptions{Level: programLevel}.NewJSONHandler(logWriter)
	slogger = slog.New(handler)
}

func Warn(args ...any) {
	slogger.Warn(emptyMsg, args...)
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
