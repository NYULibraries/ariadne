package primo

import (
	"ariadne/testutils"
	"errors"
	"fmt"
	"net/http/httputil"
	"net/url"
	"testing"
)

const testISBN = "1111111111111"

// Can't be const because need to generate a pointer to it.
var testFRBRGroupID = "2222222222"

var sharedTestCases = []struct {
	name                                string
	expectedDumpedISBNSearchHTTPRequest string
	expectedError                       error
	expectedQueryStringValues           url.Values
	queryString                         string
}{
	{
		name: "Query string with `isbn` only",
		expectedDumpedISBNSearchHTTPRequest: `GET /primo_library/libweb/webservices/rest/primo-explore/v1/pnxs?inst=NYU&limit=50&offset=0&q=isbn%2Cexact%2C` + testISBN + `&scope=all&vid=NYU HTTP/1.1
Host: bobcat.library.nyu.edu`,
		expectedError: nil,
		expectedQueryStringValues: url.Values{
			"isbn": {testISBN},
		},
		queryString: "isbn=" + testISBN,
	},
	{
		name: "3 generic query string params and `isbn`",
		expectedDumpedISBNSearchHTTPRequest: `GET /primo_library/libweb/webservices/rest/primo-explore/v1/pnxs?inst=NYU&limit=50&offset=0&q=isbn%2Cexact%2C` + testISBN + `&scope=all&vid=NYU HTTP/1.1
Host: bobcat.library.nyu.edu`,
		expectedError: nil,
		expectedQueryStringValues: url.Values{
			"param1": {"1"},
			"param2": {"2"},
			"param3": {"3"},
			"isbn":   {testISBN},
		},
		queryString: "param1=1&param2=2&param3=3&isbn=" + testISBN,
	},
	{
		name:                                "Query string without `isbn` param",
		expectedDumpedISBNSearchHTTPRequest: "",
		expectedError:                       errors.New("Could not create new Primo request: query string params do not contain required ISBN param"),
		expectedQueryStringValues: url.Values{
			"param1": {"1"},
			"param2": {"2"},
			"param3": {"3"},
		},
		queryString: "param1=1&param2=2&param3=3",
	},
	{
		name:                                "Empty query string",
		expectedDumpedISBNSearchHTTPRequest: "",
		expectedError:                       errors.New("Could not create new Primo request: query string params do not contain required ISBN param"),
		expectedQueryStringValues:           url.Values{},
		queryString:                         "",
	},
}

func TestGetISBN(t *testing.T) {
	genericParams := url.Values{
		"param1": {"1"},
		"param2": {"2"},
		"param3": {"3"},
	}
	testCases := []struct {
		queryStringValues url.Values
		expectedISBN      string
	}{
		{
			queryStringValues: testutils.MergeURLValues(
				genericParams,
				url.Values{"isbn": {testISBN}},
			),
			expectedISBN: testISBN,
		},
		{
			queryStringValues: testutils.MergeURLValues(
				genericParams,
				url.Values{"rft.isbn": {testISBN}},
			),
			expectedISBN: testISBN,
		},
		{
			queryStringValues: testutils.MergeURLValues(
				genericParams,
				url.Values{"ISBN": {testISBN}},
			),
			expectedISBN: testISBN,
		},
		{
			queryStringValues: testutils.MergeURLValues(
				genericParams,
				url.Values{"RFT.ISBN": {testISBN}},
			),
			expectedISBN: testISBN,
		},
		{
			queryStringValues: testutils.MergeURLValues(
				genericParams,
				url.Values{"iSbN": {testISBN}},
			),
			expectedISBN: testISBN,
		},
		{
			queryStringValues: testutils.MergeURLValues(
				genericParams,
				url.Values{"rFt.iSbN": {testISBN}},
			),
			expectedISBN: testISBN,
		},
		{
			queryStringValues: genericParams,
			expectedISBN:      "",
		},
		{
			queryStringValues: url.Values{},
			expectedISBN:      "",
		},
	}

	for _, testCase := range testCases {
		isbn := getISBN(testCase.queryStringValues)
		if isbn != testCase.expectedISBN {
			t.Errorf(
				"getISBN returned incorrect ISBN value for '%v': "+
					"expected '%s', got '%s'",
				testCase.queryStringValues,
				testCase.expectedISBN,
				isbn,
			)
		}
	}
}

func TestIsActiveFRBRGroupType(t *testing.T) {
	testCases := []struct {
		doc            Doc
		expectedResult bool
	}{
		{
			doc: Doc{
				Delivery: Delivery{},
				PNX: PNX{
					Facets: Facets{
						FRBRType: []string{activeFRBRGroupType},
					},
				},
			},
			expectedResult: true,
		},
		{
			doc: Doc{
				Delivery: Delivery{},
				PNX: PNX{
					Facets: Facets{
						FRBRType: []string{activeFRBRGroupType + "[FALSE]"},
					},
				},
			},
			expectedResult: false,
		},
		{
			doc:            Doc{},
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		got := isActiveFRBRGroupType(testCase.doc)
		if got != testCase.expectedResult {
			t.Errorf(
				"isActiveFRBRGroupType returned an incorrect result for '%v': "+
					"expected %t, got %t",
				testCase.doc,
				testCase.expectedResult,
				got,
			)
		}
	}
}

func TestNewPrimoRequest(t *testing.T) {
	for _, testCase := range sharedTestCases {
		testName := fmt.Sprintf("%s", testCase.queryString)
		t.Run(testName, func(t *testing.T) {
			primoRequest, err := NewPrimoRequest(testCase.queryString)
			if testCase.expectedDumpedISBNSearchHTTPRequest != "" {
				expected := testutils.NormalizeDumpedHTTPRequest(testCase.expectedDumpedISBNSearchHTTPRequest)
				got := testutils.NormalizeDumpedHTTPRequest(primoRequest.DumpedISBNSearchHTTPRequest)
				if got != expected {
					t.Errorf(
						"NewPrimoRequest returned an PrimoRequest with incorrect DumpedISBNSearchHTTPRequest string for '%s': "+
							"expected '%s', got '%s'",
						testCase.name,
						expected,
						got,
					)
				}
			}
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("NewPrimoRequest returned no error, expecting '%v'", testCase.expectedError)
				} else if err.Error() != testCase.expectedError.Error() {
					t.Errorf("NewPrimoRequest returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("NewPrimoRequest returned error '%v', expecting no errors", err)
			}

			expectedStringifiedQueryStringValues := testutils.StringifyURLValues(testCase.expectedQueryStringValues)
			if expectedStringifiedQueryStringValues != "" {
				got := testutils.StringifyURLValues(primoRequest.QueryStringValues)
				if got != expectedStringifiedQueryStringValues {
					t.Errorf(
						"NewPrimoRequest returned an PrimoRequest with incorrect stringified QueryStringValues for '%s': "+
							"expected '%s', got '%s'",
						testCase.name,
						expectedStringifiedQueryStringValues,
						got,
					)
				}
			}
		})
	}
}

func TestNewPrimoHTTPRequest(t *testing.T) {
	testCases := []struct {
		isbn                                string
		frbrGroupID                         *string
		expectedDumpedFRBRMemberHTTPRequest string
		expectedError                       error
	}{
		{
			isbn:        testISBN,
			frbrGroupID: &testFRBRGroupID,
			expectedDumpedFRBRMemberHTTPRequest: `GET /primo_library/libweb/webservices/rest/primo-explore/v1/pnxs?inst=NYU&limit=50&multiFacets=facet_frbrgroupid%2Cinclude%2C` +
				testFRBRGroupID +
				`&offset=0&q=isbn%2Cexact%2C` +
				testISBN +
				`&scope=all&vid=NYU HTTP/1.1
Host: bobcat.library.nyu.edu`,
			expectedError: nil,
		},
		{
			isbn:                                "",
			frbrGroupID:                         &testFRBRGroupID,
			expectedDumpedFRBRMemberHTTPRequest: "",
			expectedError:                       errors.New("query string params do not contain required ISBN param"),
		},
		{
			isbn:        testISBN,
			frbrGroupID: nil,
			expectedDumpedFRBRMemberHTTPRequest: `GET /primo_library/libweb/webservices/rest/primo-explore/v1/pnxs?inst=NYU&limit=50&offset=0&q=isbn%2Cexact%2C` + testISBN + `&scope=all&vid=NYU HTTP/1.1
		Host: bobcat.library.nyu.edu`,
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		testCaseName := fmt.Sprintf("ISBN: %s; FRBR Group ID: %v", testCase.isbn, testCase.frbrGroupID)
		t.Run(testCaseName, func(t *testing.T) {
			frbrMemberRequest, err := newPrimoHTTPRequest(testCase.isbn, testCase.frbrGroupID)
			if testCase.expectedDumpedFRBRMemberHTTPRequest != "" {
				gotDumpedFRBRMemberRequest, _ := httputil.DumpRequest(frbrMemberRequest, true)
				expected := testutils.NormalizeDumpedHTTPRequest(testCase.expectedDumpedFRBRMemberHTTPRequest)
				got := testutils.NormalizeDumpedHTTPRequest(string(gotDumpedFRBRMemberRequest))
				if got != expected {
					t.Errorf(
						"newPrimoFRBRMemberHTTPRequest returned a request with incorrect dumped HTTP request string: "+
							"expected '%s', got '%s'",
						expected,
						got,
					)
				}
			}
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("newPrimoFRBRMemberHTTPRequest returned no error, expecting '%v'", testCase.expectedError)
				} else if err.Error() != testCase.expectedError.Error() {
					t.Errorf("newPrimoFRBRMemberHTTPRequest returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("newPrimoFRBRMemberHTTPRequest returned error '%v', expecting no errors", err)
			}
		})
	}
}

func TestNewPrimoISBNSearchHTTPRequest(t *testing.T) {
	for _, testCase := range sharedTestCases {
		testName := fmt.Sprintf("%s", testCase.queryString)
		t.Run(testName, func(t *testing.T) {
			primoRequest, err := NewPrimoRequest(testCase.queryString)
			if testCase.expectedDumpedISBNSearchHTTPRequest != "" {
				expected := testutils.NormalizeDumpedHTTPRequest(testCase.expectedDumpedISBNSearchHTTPRequest)
				got := testutils.NormalizeDumpedHTTPRequest(primoRequest.DumpedISBNSearchHTTPRequest)
				if got != expected {
					t.Errorf(
						"NewPrimoRequest returned an PrimoRequest with incorrect DumpedISBNSearchHTTPRequest string for '%s': "+
							"expected '%s', got '%s'",
						testCase.name,
						expected,
						got,
					)
				}
			}
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("NewPrimoRequest returned no error, expecting '%v'", testCase.expectedError)
				} else if err.Error() != testCase.expectedError.Error() {
					t.Errorf("NewPrimoRequest returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("NewPrimoRequest returned error '%v', expecting no errors", err)
			}
		})
	}
}
