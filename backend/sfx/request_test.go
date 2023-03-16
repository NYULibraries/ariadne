package sfx

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestNewSFXRequest(t *testing.T) {
	var tests = []struct {
		name              string
		dumpedHTTPRequest string
		expectedError     error
		queryString       string
	}{
		{
			// This request as-is was causing SFX to return a "XSS violation occured [sic]."
			// error.  Ariadne currently remediates by replacing `sid` with `rfr_id` (set to `sid` value).
			// We do not know exactly how/why this appears to eliminate the error,
			// but somehow it does.
			// NOTE: This is the can-community-task-groups-learn-from-the-principles-of-group-therapy
			// test case from backend/testutils/testdata/test-cases.json.
			name: "Trouble-causing `sid`",
			dumpedHTTPRequest: `GET /sfxlcl41?atitle=Can+community+task+groups+learn+from+the+principles+of+group+therapy%3F&aulast=Zanbar%2C+L.&date=20181020&genre=article&isbn=&issn=19447485&issue=5&pid=Zanbar%2C+L.edselc.2-52.0-8505573399120181020Scopus%5C%C2%AE&rfr_id=EBSCO%3AScopus%5C%C2%AE&sfx.doi_url=http%3A%2F%2Fdx.doi.org&sfx.response_type=multi_obj_xml&spage=574&title=Community+Development&url_ctx_fmt=info%3Aofi%2Ffmt%3Axml%3Axsd%3Actx&volume=49 HTTP/1.1
Host: sfx.library.nyu.edu`,
			expectedError: nil,
			queryString:   "genre=article&isbn=&issn=19447485&title=Community%20Development&volume=49&issue=5&date=20181020&atitle=Can%20community%20task%20groups%20learn%20from%20the%20principles%20of%20group%20therapy?&aulast=Zanbar,%20L.&spage=574&sid=EBSCO:Scopus\\®&pid=Zanbar,%20L.edselc.2-52.0-8505573399120181020Scopus\\®",
		},
		// These unit tests were originally written when Ariadne was making POST
		// requests to the SFX API, with complicated query string params validation
		// and massaging and a somewhat brittle XML request body.  There were plenty
		// of error conditions to test for.
		// It is extremely difficult and perhaps impossible to cause the current
		// `NewSFXRequest` code to return an error.  Should we need to write unit
		// tests to check for errors, we would use test cases with this structure:
		//
		// {
		//   name:              "",
		//	 dumpedHTTPRequest: "",
		//	 expectedError:     errors.New(""),
		//	 queryString:       string,
		// },
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.queryString)
		t.Run(testName, func(t *testing.T) {
			sfxRequest, err := NewSFXRequest(testCase.queryString)
			if testCase.dumpedHTTPRequest != "" {
				expected := normalizeDumpedHTTPRequest(testCase.dumpedHTTPRequest)
				got := normalizeDumpedHTTPRequest(sfxRequest.DumpedHTTPRequest)
				if got != expected {
					t.Errorf(
						"NewSFXRequest returned an SFXRequest with incorrect DumpedHTTPRequest string for %s: "+
							"expected '%s', got '%s'",
						testCase.name,
						expected,
						got,
					)
				}
			}
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("NewSFXRequest returned no error, expecting '%v'", testCase.expectedError)
				} else if err.Error() != testCase.expectedError.Error() {
					t.Errorf("NewSFXRequest returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("NewSFXRequest returned error '%v', expecting no errors", err)
			}
		})
	}
}

func TestFilterOpenURLParams(t *testing.T) {
	var testCases = []struct {
		testName    string
		queryString url.Values
		expected    url.Values
	}{
		{"query string contains sid", map[string][]string{"sid": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}, map[string][]string{"rfr_id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}},
		{"query string doesn't contain sid", map[string][]string{"id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}, map[string][]string{"id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			actual := filterOpenURLParams(testCase.queryString)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("filterOpenURLParams returned '%v', expecting '%v'", actual, testCase.expected)
			}
		})
	}
}

func normalizeDumpedHTTPRequest(dumpedHTTPRequest string) string {
	multipleWhitespaceRegexp := regexp.MustCompile(`\s+`)

	return multipleWhitespaceRegexp.ReplaceAllString(strings.TrimSpace(dumpedHTTPRequest), " ")
}
