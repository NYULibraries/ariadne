package log

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

// This is for set level closures.  Examples:
// SetLevel(LevelDebug)
// SetLevelByString("warn")
// These closures are passed into the main testing function which contains code
// that we want to keep DRY -- i.e. we don't want to duplicate it just to change
// `SetLevel(...)` to `SetLevelByString(...)` and vice-versa.
type setLevelFunction func()

// For a JSON string like this:
// {"time":"2023-04-17T11:11:16.545242-04:00","level":"DEBUG","msg":"","message":"debug"}
// ...match "DEBUG" and "debug" in the "level" and "message" properties.
var simpleLogEntryStringRegexp = regexp.MustCompile(
	`.*"level":"([A-Z]+)","msg":"","message":"([a-z]+)"}`,
)

const messageKey = "message"

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

// Tests whether GetLevelOptionStringForLogLevel() returns the correct human
// read-able strings for the Level* values (which are basically `int` values).
func TestGetLevelOptionStringForLogLevel(t *testing.T) {
	testCases := []struct {
		name                      string
		expectedLevelOptionString string
		level                     Level
	}{
		{name: "LevelDebug", expectedLevelOptionString: "debug", level: LevelDebug},
		{name: "LevelInfo", expectedLevelOptionString: "info", level: LevelInfo},
		{name: "LevelWarn", expectedLevelOptionString: "warn", level: LevelWarn},
		{name: "LevelError", expectedLevelOptionString: "error", level: LevelError},
		{name: "LevelDisabled", expectedLevelOptionString: "disabled", level: LevelDisabled},
	}

	for _, testCase := range testCases {
		actualLevelOptionString := GetLevelOptionStringForLogLevel(testCase.level)
		if actualLevelOptionString != testCase.expectedLevelOptionString {
			t.Errorf("Incorrect log level string for %s:\nExpected: \"%s\"\nActual: \"%s\"",
				testCase.name, testCase.expectedLevelOptionString, actualLevelOptionString)
		}
	}
}

// Tests whether GetValidLevelOptionStrings() returns the correct list of level
// options strings, ordered by increasing level of severity.
func TestGetValidLevelOptionStrings(t *testing.T) {
	expectedLevelOptionStrings := strings.Join([]string{
		"debug",
		"info",
		"warn",
		"error",
		"disabled",
	}, "\n")

	actualLevelOptionStrings := strings.Join(GetValidLevelOptionStrings(), "\n")

	if actualLevelOptionStrings != expectedLevelOptionStrings {
		t.Errorf("GetValidLevelOptionStrings() did not return correct list of strings:\nExpected:\n\n%s\n\nActual:\n\n%s",
			expectedLevelOptionStrings, actualLevelOptionStrings)
	}
}

// Tests that both `SetLevel()` and `SetLevelByString()` work correctly.
func TestSetLevel(t *testing.T) {
	// This reflects the log output that we expect the test fixture function `logSeries()`
	// to generate at level debug.
	expectedLogSeriesFull := []string{
		"DEBUG|debug",
		"INFO|info",
		"WARN|warn",
		"ERROR|error",
	}

	// We explicitly set `name` for test cases even though this package has available
	// a `GetLevelOptionStringForLogLevel` function which can provide the corresponding
	// string for a level.  The test case name is used in test failure messages,
	// so a bug in `GetLevelOptionStringForLogLevel()` could potentially produce
	// incorrect test fail messages. For example: at one point `LevelError` was
	// accidentally mapped to `slog.LevelWarn`, leading to the test failure message
	// indicating a problem with `SetLevel(LevelWarn)` rather than `SetLevel(LevelError)`.
	testCases := []struct {
		name  string
		level Level
	}{
		{name: "Debug", level: LevelDebug},
		{name: "Info", level: LevelInfo},
		{name: "Warn", level: LevelWarn},
		{name: "Error", level: LevelError},
		{name: "Disabled", level: LevelDisabled},
	}

	for i, testCase := range testCases {
		expectedLogSeriesForLevelString := strings.Join(expectedLogSeriesFull[i:], "\n")
		testSetLevelByLevel(t, testCase.name, testCase.level, expectedLogSeriesForLevelString)
		testSetLevelByLevelString(t, testCase.name, strings.ToLower(testCase.name), expectedLogSeriesForLevelString)
	}
}

// Helper function to generate a simple multiline string view of the log output JSON,
// showing only the data of interest.
//
// Example:
//
// For:
//
// {"time":"2023-04-17T11:11:16.545242-04:00","level":"DEBUG","msg":"","message":"debug"}
// {"time":"2023-04-17T11:11:16.545399-04:00","level":"INFO","msg":"","message":"info"}
// {"time":"2023-04-17T11:11:16.545403-04:00","level":"WARN","msg":"","message":"warn"}
// {"time":"2023-04-17T11:11:16.545406-04:00","level":"ERROR","msg":"","message":"error"}
//
// ...return:
//
// DEBUG|debug
// INFO|info
// WARN|warn
// ERROR|error
func getLogSeriesString(logOutput string) string {
	var logSeries []string

	results := simpleLogEntryStringRegexp.FindAllStringSubmatch(logOutput, -1)

	for _, result := range results {
		logSeries = append(logSeries, fmt.Sprintf("%s|%s", result[1], result[2]))
	}

	return strings.Join(logSeries, "\n")
}

// Helper fixture function that issues a series of simple log commands in ascending
// order of severity.
func logSeries() {
	Debug(messageKey, "debug")
	Info(messageKey, "info")
	Warn(messageKey, "warn")
	Error(messageKey, "error")
}

// Called in the main loop for `TestSetLevel()`.  It calls the workhorse test function
// `testSetLevel()`, passing in a closure allowing `testSetLevel()` to set the
// level without needing to know which setter to use and what level is being
// specified.  In this case, the closure uses `SetLevel(Level)`.
func testSetLevelByLevel(t *testing.T, testCaseName string, level Level, expectedLogSeriesForLevelString string) {
	var setLevelClosure setLevelFunction
	setLevelClosure = func() {
		SetLevel(level)
	}
	testSetLevel(t, "SetLevel", setLevelClosure, testCaseName, expectedLogSeriesForLevelString)
}

// Called in the main loop for `TestSetLevel()`.  It calls the workhorse test function
// `testSetLevel()`, passing in a closure allowing `testSetLevel()` to set the
// level without needing to know which setter to use and what level is being
// specified.  In this case, the closure uses `SetLevelByLevelString(string)`.
func testSetLevelByLevelString(t *testing.T, testCaseName string, levelString string, expectedLogSeriesForLevelString string) {
	var setLevelClosure setLevelFunction
	setLevelClosure = func() {
		err := SetLevelByString(levelString)
		if err != nil {
			t.Fatalf("SetLevelByString(\"%s\") failed with error: %s",
				levelString, err)
		}
	}
	testSetLevel(t, "SetLevelByString", setLevelClosure, testCaseName, expectedLogSeriesForLevelString)
}

// The actual workhorse test function.  Do a test using the provided `setLevelClosure`,
// examples of which are:
// SetLevel(LevelDebug)
// SetLevelByString("warn")
// We implement it this way because we don't want to have two workhorse functions
// that differ only by a single statement: e.g. `SetLevel(...)`vs. `SetLevelByString(...)`.
func testSetLevel(t *testing.T, setLevelFunctionName string, setLevelClosure setLevelFunction,
	testCaseName string, expectedLogSeriesString string) {
	// Capture all log output in a `bytes.Buffer`.
	var logOutput bytes.Buffer
	logOutputWriter := bufio.NewWriter(&logOutput)
	SetOutput(logOutputWriter)

	// We've been provided a set level function which has been closed off with the
	// desired level.  We just need to call it to set the package log level.
	setLevelClosure()

	// Logs simple messages for each testCase in ascending order of severity.
	logSeries()

	// Must flush the writer to prevent risk truncation of `logOutput`
	err := logOutputWriter.Flush()
	if err != nil {
		t.Fatalf("Error calling `logOutputWriter.Flush`: %s", err)
	}

	// Example of actual log series string to be compared against the expected log series string:
	//
	// DEBUG|debug
	// INFO|info
	// WARN|warn
	// ERROR|error
	actualLogSeriesString := getLogSeriesString(logOutput.String())

	// Example output for a bug where LevelError is accidentally mapped to slog.LevelWarn
	// (this is a real bug that was caught by this test):
	// --------
	// Incorrect log series for level Error:
	// Expected:
	//
	// ERROR|error
	//
	// Actual:
	//
	// WARN|warn
	// ERROR|error
	// --------
	if actualLogSeriesString != expectedLogSeriesString {
		t.Errorf("Incorrect log series for test case \"%s\", using function `%s`:\nExpected:\n\n%s\n\nActual:\n\n%s",
			testCaseName, setLevelFunctionName, expectedLogSeriesString, actualLogSeriesString)
	}
}
