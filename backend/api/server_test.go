package api

import (
	"ariadne/primo"
	"ariadne/sfx"
	"ariadne/testutils"
	"ariadne/util"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

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
				err = updateGoldenFile(testCase, body)
				if err != nil {
					t.Fatalf("Error updating golden file: %s", err)
				}
			}

			goldenValue, err := testutils.GetGoldenValue(testCase)
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

				goldenFile := testutils.GoldenFile(testCase)
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

func tmpFile(testCase testutils.TestCase) string {
	return "testdata/server/tmp/actual/" + testCase.Key + ".json"
}

func updateGoldenFile(testCase testutils.TestCase, bytes []byte) error {
	return os.WriteFile(testutils.GoldenFile(testCase), bytes, 0644)
}

func writeActualToTmp(testCase testutils.TestCase, actual string) error {
	return os.WriteFile(tmpFile(testCase), []byte(actual), 0644)
}
