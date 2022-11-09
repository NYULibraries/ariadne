package api

import (
	"ariadne/sfx"
	"ariadne/util"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//go:embed testdata/server/test-cases.json
var testCasesJSON []byte

type TestCase struct {
	// Identifier used for fixture and golden file basename.
	Key string
	// Human-readable name/description of test case
	Name string
	// OpenURL querystring
	QueryString string
}

var testCases []TestCase
var updateGoldenFiles = flag.Bool("update-golden-files", false, "update the golden files")

// --update-sfx-fake-responses flag?
// Ideally we also want to have a flag for updating SFX fake response fixture files,
// but it appears that ordering of elements in the SFX response XML and the elements
// in the escaped XML in <perldata> is not stable.
// See comment in monday.com ticket "Add sample integration test for OpenURL resolver":
// https://nyu-lib.monday.com/boards/765008773/pulses/3073776565/posts/1676502313
// Thus the same request submitted multiple times in less than a second
// might end up generating responses that differ only in element ordering.  If this
// in fact is confirmed to be the case, in order for --update-sfx-fake-responses
// to be useful, we would need to write some utility code to normalize the SFX
// responses before writing out the fixture files.

func TestMain(m *testing.M) {
	flag.Parse()

	err := json.Unmarshal(testCasesJSON, &testCases)
	if err != nil {
		panic(fmt.Sprintf("Error reading test cases file: %s", err))
	}

	os.Exit(m.Run())
}

func TestResponseJSONRoute(t *testing.T) {
	var currentTestCase TestCase

	fakeSFXServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sfxFakeResponse, err := getSFXFakeResponse(currentTestCase)
			if err != nil {
				t.Fatal(err)
			}

			_, err = fmt.Fprint(w, sfxFakeResponse)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer fakeSFXServer.Close()

	sfx.SetSFXURL(fakeSFXServer.URL)

	router := NewRouter()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			currentTestCase = testCase

			request, err := http.NewRequest(
				"GET",
				"/v0/?"+testCase.QueryString,
				nil,
			)
			if err != nil {
				t.Fatalf("Error creating new HTTP request: %s", err)
			}

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			response := responseRecorder.Result()
			body, _ := io.ReadAll(response.Body)

			if *updateGoldenFiles {
				err = updateGoldenFile(testCase, body)
				if err != nil {
					t.Fatalf("Error updating golden file: %s", err)
				}
			}

			goldenValue, err := getGoldenValue(testCase)
			if err != nil {
				t.Fatalf("Error retrieving golden value for test case \"%s\": %s",
					testCase.Name, err)
			}

			actualValue := string(body)
			if actualValue != goldenValue {
				err := writeActualToTmp(testCase, actualValue)
				if err != nil {
					t.Fatalf("Error writing actual temp file for test case \"%s\": %s",
						testCase.Name, err)
				}

				goldenFile := goldenFile(testCase)
				actualFile := tmpFile(testCase)
				diff, err := util.Diff(goldenFile, actualFile)
				if err != nil {
					t.Fatalf("Error diff'ing %s vs. %s: %s\n"+
						"Manually diff these files to determine the reasons for test failure.",
						goldenFile, actualFile, err)
				}

				t.Errorf("golden and actual values do not match\noutput of `diff %s %s`:\n%s\n",
					goldenFile, actualFile, diff)
			}
		})
	}
}

func getGoldenValue(testCase TestCase) (string, error) {
	return getTestdataFileContents(goldenFile(testCase))
}

func getSFXFakeResponse(testCase TestCase) (string, error) {
	return getTestdataFileContents(sfxFakeResponseFile(testCase))
}

func getTestdataFileContents(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return filename, err
	}

	return string(bytes), nil
}

func getTestCases() ([]TestCase, error) {
	testCases := []TestCase{}

	err := json.Unmarshal(testCasesJSON, &testCases)
	if err != nil {
		return testCases, err
	}

	return testCases, nil
}

func goldenFile(testCase TestCase) string {
	return "testdata/server/golden/" + testCase.Key + ".json"
}

func sfxFakeResponseFile(testCase TestCase) string {
	return "testdata/server/fixtures/sfx-fake-responses/" + testCase.Key + ".xml"
}

func tmpFile(testCase TestCase) string {
	return "testdata/server/tmp/actual/" + testCase.Key + ".json"
}

func updateGoldenFile(testCase TestCase, bytes []byte) error {
	return os.WriteFile(goldenFile(testCase), bytes, 0644)
}

func writeActualToTmp(testCase TestCase, actual string) error {
	return os.WriteFile(tmpFile(testCase), []byte(actual), 0644)
}
