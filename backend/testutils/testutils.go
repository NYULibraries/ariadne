package testutils

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed testdata/test-cases.json
var TestCasesJSON []byte

type TestCase struct {
	// Identifier used for fixture and golden file basename.
	Key string
	// Human-readable name/description of test case
	Name string
	// OpenURL querystring
	QueryString string
}

var TestCases []TestCase

var testutilsPath string

func init() {
	// Get parent directory for this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ERROR: `runtime.Caller(0)` failed")
	}

	// We need this absolute path to this current directory to allow callers
	// from other packages to be able to retrieve the stuff in testdata/.
	testutilsPath = filepath.Dir(filename)

	err := json.Unmarshal(TestCasesJSON, &TestCases)
	if err != nil {
		panic(fmt.Sprintf("Error reading test cases file: %s", err))
	}

}

func GetGoldenValue(testCase TestCase) (string, error) {
	return GetTestdataFileContents(GoldenFile(testCase))
}

func GetSFXFakeResponse(testCase TestCase) (string, error) {
	return GetTestdataFileContents(sfxFakeResponseFile(testCase))
}

func GetTestdataFileContents(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return filename, err
	}

	return string(bytes), nil
}

func GoldenFile(testCase TestCase) string {
	return testutilsPath + "/testdata/golden/" + testCase.Key + ".json"
}

func sfxFakeResponseFile(testCase TestCase) string {
	return testutilsPath + "/testdata/fixtures/sfx-fake-responses/" + testCase.Key + ".xml"
}
