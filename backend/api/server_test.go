package api

import (
	"ariadne/log"
	"ariadne/primo"
	"ariadne/sfx"
	"ariadne/testutils"
	"ariadne/util"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"
)

const elidedHost = "Host: [ELIDED]"
const elidedDatestamp = "Date: [ELIDED]"
const elidedTimestamp = "\"time\":\"[ELIDED]\""

const loggingTestCaseKey = "contrived-frbr-group-test-case"

var logOutputStringDatestampRegexp = regexp.MustCompile("Date:.*GMT")
var logOutputStringHostRegexp = regexp.MustCompile("Host: 127.0.0.1:\\d*")
var logOutputStringTimestampRegexp = regexp.MustCompile("\"time\":\"[^\"]*\"")

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

	os.Exit(m.Run())
}

func TestResponseJSONRoute(t *testing.T) {
	var currentTestCase testutils.TestCase

	// Set up Primo service fake
	fakePrimoServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				t.Fatal(err)
			}

			// There potentially two kinds of requests:
			//     - ISBN search request: this is the initial request that is
			//       always made if Primo is being used at all
			//     - FRBR member search request: if the response to the initial
			//       ISBN search request returns docs that indicate an active FRBR
			//       group, more requests are made with an extra query param added
			//       to the query string of the ISBN search request.
			var primoFakeResponse string
			if params.Get(primo.FRBRMemberSearchQueryParamName) == "" {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseISBNSearch(currentTestCase)
			} else {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseFRBRMemberSearch(currentTestCase)
			}

			if err != nil {
				t.Fatal(err)
			}

			_, err = fmt.Fprint(w, primoFakeResponse)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer fakePrimoServer.Close()

	primo.SetPrimoURL(fakePrimoServer.URL)

	// Set up SFX service fake
	fakeSFXServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sfxFakeResponse, err := testutils.GetSFXFakeResponse(currentTestCase)
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

	// Disable logging or else we'll have a ton on noise in the test results
	// output.
	log.SetLevel(log.LevelDisabled)

	for _, testCase := range testutils.TestCases {
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
				err = updateAPIResponseGoldenFile(testCase, body)
				if err != nil {
					t.Fatalf("Error updating golden file: %s", err)
				}
			}

			goldenValue, err := testutils.GetAPIResponseGoldenValue(testCase)
			if err != nil {
				t.Fatalf("Error retrieving golden value for test case \"%s\": %s",
					testCase.Name, err)
			}

			actualValue := string(body)
			if actualValue != goldenValue {
				err := writeActualAPIResponseToTmp(testCase, actualValue)
				if err != nil {
					t.Fatalf("Error writing actual temp file for test case \"%s\": %s",
						testCase.Name, err)
				}

				goldenFile := testutils.APIResponseGoldenFile(testCase)
				actualFile := tmpAPIResponseFile(testCase)
				diff, err := util.DiffFiles(goldenFile, actualFile)
				if err != nil {
					t.Fatalf("Error diff'ing %s vs. %s: %s\n"+
						"Manually diff these files to determine the reasons for test failure.",
						goldenFile, actualFile, err)
				}

				t.Errorf("golden and actual values do not match:\n%s\n", diff)
			}
		})
	}
}

func TestLogging(t *testing.T) {
	var loggingTestCase testutils.TestCase

	// Set up SFX service fake
	fakeSFXServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sfxFakeResponse, err := testutils.GetSFXFakeResponse(loggingTestCase)
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

	// Set up Primo service fake
	fakePrimoServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				t.Fatal(err)
			}

			// There potentially two kinds of requests:
			//     - ISBN search request: this is the initial request that is
			//       always made if Primo is being used at all
			//     - FRBR member search request: if the response to the initial
			//       ISBN search request returns docs that indicate an active FRBR
			//       group, more requests are made with an extra query param added
			//       to the query string of the ISBN search request.
			var primoFakeResponse string
			if params.Get(primo.FRBRMemberSearchQueryParamName) == "" {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseISBNSearch(loggingTestCase)
			} else {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseFRBRMemberSearch(loggingTestCase)
			}

			if err != nil {
				t.Fatal(err)
			}

			_, err = fmt.Fprint(w, primoFakeResponse)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer fakePrimoServer.Close()

	primo.SetPrimoURL(fakePrimoServer.URL)

	router := NewRouter()

	for _, testCase := range testutils.TestCases {
		if testCase.Key == loggingTestCaseKey {
			loggingTestCase = testCase
			break
		}
	}

	for _, level := range []log.Level{log.LevelDebug, log.LevelInfo} {
		// Needed for golden file stuff
		levelString := log.GetLevelOptionStringForLogLevel(level)

		// Set logging level and redirect output to a buffer.
		log.SetLevel(level)
		var logOutput bytes.Buffer
		logOutputWriter := bufio.NewWriter(&logOutput)
		log.SetOutput(logOutputWriter)

		t.Run(loggingTestCase.Name, func(t *testing.T) {
			request, err := http.NewRequest(
				"GET",
				"/v0/?"+loggingTestCase.QueryString,
				nil,
			)
			if err != nil {
				t.Fatalf("Error creating new HTTP request: %s", err)
			}

			// We're not recording any responses at the moment, but we need to pass
			// in a http.ResponseWriter anyway to router.ServeHTTP, so why not make
			// it a recorder.
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			err = logOutputWriter.Flush()
			if err != nil {
				t.Fatalf("logOutputWriter.Flush() error: %s", err)
			}

			actualLogOutputString := normalizeLogOutputString(logOutput.String())

			if *updateGoldenFiles {
				err = updateLogOutputGoldenFile(loggingTestCase, levelString, actualLogOutputString)
				if err != nil {
					t.Fatalf("Error updating golden file: %s", err)
				}
			}

			goldenValue, err := testutils.GetLogOutputGoldenValue(loggingTestCase, levelString)
			if err != nil {
				t.Fatalf("Error retrieving golden value for test case \"%s\": %s",
					loggingTestCase.Name, err)
			}

			if actualLogOutputString != goldenValue {
				// We don't programmatically diff an actual file vs. the golden file
				// because the golden file contents are not normalized, but we do
				// nevertheless want to write out an actual file in case the user
				// wants to do a manual diff themselves.
				err := writeActualLogOutputToTmp(loggingTestCase, actualLogOutputString)
				if err != nil {
					t.Fatalf("Error writing actual temp file for test case \"%s\": %s",
						loggingTestCase.Name, err)
				}
				actualFile := tmpLogOutputFile(loggingTestCase)
				goldenFile := testutils.LogOutputGoldenFile(loggingTestCase, levelString)
				// We pass in the file paths because util.DiffStrings prints labels
				// in the output header, and we want to use the paths for these labels.
				diff := util.DiffStrings(goldenFile, goldenValue, actualFile, actualLogOutputString)

				t.Errorf("golden and actual values do not match:\n%s\n", diff)
			}
		})
	}
}

func normalizeLogOutputString(logOutputString string) string {
	result := logOutputStringDatestampRegexp.ReplaceAllString(logOutputString, elidedDatestamp)
	result = logOutputStringHostRegexp.ReplaceAllString(result, elidedHost)
	result = logOutputStringTimestampRegexp.ReplaceAllString(result, elidedTimestamp)

	return result
}

func tmpAPIResponseFile(testCase testutils.TestCase) string {
	return "testdata/server/tmp/actual/api-responses/" + testCase.Key + ".json"
}

func tmpLogOutputFile(testCase testutils.TestCase) string {
	return "testdata/server/tmp/actual/log-output/" + testCase.Key + ".txt"
}

func updateAPIResponseGoldenFile(testCase testutils.TestCase, bytes []byte) error {
	return os.WriteFile(testutils.APIResponseGoldenFile(testCase), bytes, 0644)
}

func updateLogOutputGoldenFile(testCase testutils.TestCase, level string, newGoldenValue string) error {
	return os.WriteFile(testutils.LogOutputGoldenFile(testCase, level), []byte(newGoldenValue), 0644)
}

func writeActualAPIResponseToTmp(testCase testutils.TestCase, actual string) error {
	return os.WriteFile(tmpAPIResponseFile(testCase), []byte(actual), 0644)
}

func writeActualLogOutputToTmp(testCase testutils.TestCase, actual string) error {
	return os.WriteFile(tmpLogOutputFile(testCase), []byte(actual), 0644)
}
