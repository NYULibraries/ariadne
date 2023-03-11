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
		queryStringValues url.Values
	}{
		{
			name: "Trouble-causing `sid`",
			dumpedHTTPRequest: `GET /sfxlcl41?rfr_id=genre%3Darticle%26isbn%3D%26issn%3D19447485%26title%3DCommunity%2520Development%26volume%3D49%26issue%3D5%26date%3D20181020%26atitle%3DCan%2520community%2520task%2520groups%2520learn%2520from%2520the%2520principles%2520of%2520group%2520therapy%3F%26aulast%3DZanbar%2C%2520L.%26spage%3D574%26sid%3DEBSCO%3AScopus%5C%5Cu00ae%26pid%3DZanbar%2C%2520L.edselc.2-52.0-8505573399120181020Scopus%5C%5Cu00ae&sfx.doi_url=http%3A%2F%2Fdx.doi.org&sfx.ignore_date_threshold=1&sfx.response_type=multi_obj_xml&sfx.show_availability=1&url_ctx_fmt=info%3Aofi%2Ffmt%3Axml%3Axsd%3Actx HTTP/1.1
Host: sfx.library.nyu.edu`,
			expectedError:     nil,
			queryStringValues: map[string][]string{"sid": {"genre=article&isbn=&issn=19447485&title=Community%20Development&volume=49&issue=5&date=20181020&atitle=Can%20community%20task%20groups%20learn%20from%20the%20principles%20of%20group%20therapy?&aulast=Zanbar,%20L.&spage=574&sid=EBSCO:Scopus\\\\u00ae&pid=Zanbar,%20L.edselc.2-52.0-8505573399120181020Scopus\\\\u00ae"}},
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
		//	 queryStringValues: map[string][]string{},
		// },
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.queryStringValues)
		t.Run(testName, func(t *testing.T) {
			sfxRequest, err := NewSFXRequest(testCase.queryStringValues)
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
