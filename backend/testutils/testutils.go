package testutils

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
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

// We need to get the absolute path to this package in order to enable the function
// for golden file and fixture file retrieval to be called from other packages
// which would not be able to resolve the hardcoded relative paths used here.
func init() {
	// The `filename` string is the absolute path to this source file, which should
	// be located at the root of the package directory.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ERROR: `runtime.Caller(0)` failed")
	}

	// Get the path to the parent directory of this file.  Again, this is assuming
	// that this `init()` function is defined in a package top level file -- or
	// more precisely, that this file is in the same directory at the `testdata/`
	// directory that is referenced in the relative paths used in the functions
	// defined in this file.
	testutilsPath = filepath.Dir(filename)

	err := json.Unmarshal(TestCasesJSON, &TestCases)
	if err != nil {
		panic(fmt.Sprintf("Error reading test cases file: %s", err))
	}

}

func GetGoldenValue(testCase TestCase) (string, error) {
	return GetTestdataFileContents(GoldenFile(testCase))
}

func GetPrimoFakeResponseISBNSearch(testCase TestCase) (string, error) {
	return GetTestdataFileContents(primoFakeResponseFileISBNSearch(testCase))
}

func GetPrimoFakeResponseFRBRMemberSearch(testCase TestCase) (string, error) {
	return GetTestdataFileContents(primoFakeResponseFileFRBRMemberSearch(testCase))
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

// Returns a url.Values consisting of everything in urlValues1 with everything in
// urlValues2 added, where collision of keys is handled by urlValues2 values
// overriding those in urlValues1.
func MergeURLValues(urlValues1 url.Values, urlValues2 url.Values) url.Values {
	mergedURLValues := url.Values{}
	for queryParamName, queryParamValue := range urlValues1 {
		mergedURLValues[queryParamName] = queryParamValue
	}

	for queryParamName, queryParamValue := range urlValues2 {
		mergedURLValues[queryParamName] = queryParamValue
	}

	return mergedURLValues
}

func NormalizeDumpedHTTPRequest(dumpedHTTPRequest string) string {
	multipleWhitespaceRegexp := regexp.MustCompile(`\s+`)

	return multipleWhitespaceRegexp.ReplaceAllString(strings.TrimSpace(dumpedHTTPRequest), " ")
}

func NormalizeDumpedHTTPResponse(dumpedHTTPRequest string) string {
	multipleWhitespaceRegexp := regexp.MustCompile(`\s+`)

	return multipleWhitespaceRegexp.ReplaceAllString(strings.TrimSpace(dumpedHTTPRequest), " ")
}

func StringifyURLValues(urlValues url.Values) string {
	return fmt.Sprintf("%v", urlValues)
}

func primoFakeResponseFileFRBRMemberSearch(testCase TestCase) string {
	return testutilsPath + "/testdata/fixtures/primo-fake-responses/frbr-member-search-data/" + testCase.Key + ".json"
}

func primoFakeResponseFileISBNSearch(testCase TestCase) string {
	return testutilsPath + "/testdata/fixtures/primo-fake-responses/" + testCase.Key + ".json"
}

func sfxFakeResponseFile(testCase TestCase) string {
	return testutilsPath + "/testdata/fixtures/sfx-fake-responses/" + testCase.Key + ".xml"
}
