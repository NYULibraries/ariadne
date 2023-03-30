package primo

import (
	"ariadne/testutils"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestAddHTTPResponseData(t *testing.T) {
	testCases := []struct {
		httpResponse                *http.Response
		expectedAPIResponses        []APIResponse
		expectedDumpedHTTPResponses []string
		expectedError               error
	}{
		{
			httpResponse: fakePrimoISBNSearchHTTPResponse,
			expectedAPIResponses: []APIResponse{
				fakePrimoISBNSearchAPIResponse,
			},
			expectedDumpedHTTPResponses: []string{
				testutils.NormalizeDumpedHTTPResponse(fakeDumpedPrimoISBNSearchHTTPResponse),
			},
			expectedError: nil,
		},
		{
			httpResponse:         fakePrimoISBNSearchHTTPResponseInvalid,
			expectedAPIResponses: []APIResponse{},
			expectedDumpedHTTPResponses: []string{
				testutils.NormalizeDumpedHTTPResponse(fakeDumpedPrimoISBNSearchHTTPResponseInvalid),
			},
			expectedError: errors.New("invalid character '<' looking for beginning of value"),
		},
	}

	for _, testCase := range testCases {
		primoResponse := PrimoResponse{}
		returnedAPIResponse, err := primoResponse.addHTTPResponseData(testCase.httpResponse)

		// Test returnedAPIResponse, which in this case should be the same as primoResponse.APIResponses[0]
		var stringifiedExpectedAPIResponse string
		if len(testCase.expectedAPIResponses) > 0 {
			stringifiedExpectedAPIResponse = stringifyAPIResponse(testCase.expectedAPIResponses[0])
		} else {
			stringifiedExpectedAPIResponse = stringifyAPIResponse(APIResponse{})
		}
		stringifiedGotAPIResponse := stringifyAPIResponse(returnedAPIResponse)
		if stringifiedGotAPIResponse != stringifiedExpectedAPIResponse {
			t.Errorf("addHTTPResponseData returned incorrect APIResponse: "+
				"expected \"%s\"; got \"%s\"", stringifiedExpectedAPIResponse, stringifiedGotAPIResponse)
		}

		// Test primoResponse.APIResponses
		numAPIResponses := len(primoResponse.APIResponses)
		if numAPIResponses != len(testCase.expectedAPIResponses) {
			t.Errorf("addHTTPResponseData added incorrect number of APIResponse objects: "+
				"expected \"%d\"; got \"%d\"", numAPIResponses, len(testCase.expectedDumpedHTTPResponses))
		}
		for i := 0; i < numAPIResponses; i++ {
			stringifiedExpectedAPIResponse := stringifyAPIResponse(testCase.expectedAPIResponses[i])
			stringifiedGotAPIResponse := stringifyAPIResponse(primoResponse.APIResponses[i])
			if stringifiedGotAPIResponse != stringifiedExpectedAPIResponse {
				t.Errorf("addHTTPResponseData added incorrect APIResponse: "+
					"expected \"%s\"; got \"%s\"", stringifiedExpectedAPIResponse, stringifiedGotAPIResponse)
			}
		}

		// Test primoResponse.DumpedHTTPResponses
		numDumpedHTTPResponses := len(primoResponse.DumpedHTTPResponses)
		if numDumpedHTTPResponses != len(testCase.expectedDumpedHTTPResponses) {
			t.Errorf("addHTTPResponseData added incorrect number of DumpedHTTPResponse objects: "+
				"expected \"%d\"; got \"%d\"", numDumpedHTTPResponses, len(testCase.expectedDumpedHTTPResponses))
		}
		for i := 0; i < numDumpedHTTPResponses; i++ {
			expectedDumpedHTTPResponse := testCase.expectedDumpedHTTPResponses[i]
			gotDumpedHTTResponse :=
				testutils.NormalizeDumpedHTTPResponse(primoResponse.DumpedHTTPResponses[i])
			if gotDumpedHTTResponse != expectedDumpedHTTPResponse {
				t.Errorf("addHTTPResponseData added incorrect DumpedHTTPResponse: "+
					"expected \"%s\"; got \"%s\"", expectedDumpedHTTPResponse, gotDumpedHTTResponse)
			}
		}

		if testCase.expectedError != nil {
			if err == nil {
				t.Errorf("addHTTPResponseData returned no error, expecting '%v'", testCase.expectedError)
			} else if err.Error() != testCase.expectedError.Error() {
				t.Errorf("addHTTPResponseData returned error '%v', expecting '%v'", err, testCase.expectedError)
			}
		}
		if err != nil && testCase.expectedError == nil {
			t.Errorf("addHTTPResponseData returned error '%v', expecting no errors", err)
		}
	}
}

func TestAddLinks(t *testing.T) {
	testCases := []struct {
		doc           Doc
		expectedLinks []Link
	}{
		{
			doc: fakePrimoISBNSearchAPIResponse.Docs[0],
			expectedLinks: []Link{
				{
					HyperlinkText: "2",
					LinkURL:       "https://fake.com/2/",
					LinkType:      "http://purl.org/pnx/linkType/linktorsrc",
				},
				{
					HyperlinkText: "4",
					LinkURL:       "https://fake.com/4/",
					LinkType:      "http://purl.org/pnx/linkType/linktorsrc",
				},
			},
		},
	}

	for _, testCase := range testCases {
		primoResponse := PrimoResponse{}
		primoResponse.addLinks(testCase.doc)
		expectedStringifiedLinks := stringifyLinks(testCase.expectedLinks)
		gotStringifiedLinks := stringifyLinks(primoResponse.Links)
		if gotStringifiedLinks != expectedStringifiedLinks {
			t.Errorf("addLinks did not correctly add links: "+
				"expected \"%s\"; got \"%s\"", expectedStringifiedLinks, gotStringifiedLinks)
		}
	}
}

func TestIsFound(t *testing.T) {
	testCases := []struct {
		name           string
		links          []Link
		expectedResult bool
	}{
		{
			name: "PrimoResponse has at least one link",
			links: []Link{
				Link{},
			},
			expectedResult: true,
		},
		{
			name:           "PrimoResponse has no links",
			links:          []Link{},
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		primoResponse := PrimoResponse{
			Links: testCase.links,
		}
		got := primoResponse.IsFound()
		if primoResponse.IsFound() != testCase.expectedResult {
			t.Errorf(
				"IsFound returned an incorrect result for test case '%s': "+
					"expected %t, got %t",
				testCase.name,
				testCase.expectedResult,
				got,
			)
		}
	}
}

func TestIsMatch(t *testing.T) {
	testCases := []struct {
		name           string
		frbrGroupDoc   Doc
		isbn           string
		expectedResult bool
	}{
		{
			name:           "ISBN match found",
			frbrGroupDoc:   fakePrimoISBNSearchAPIResponse.Docs[0],
			isbn:           "3333333333333",
			expectedResult: true,
		},
		{
			name:           "ISBN match not found",
			frbrGroupDoc:   fakePrimoISBNSearchAPIResponse.Docs[0],
			isbn:           "5555555555555",
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		got := isMatch(testCase.frbrGroupDoc, testCase.isbn)
		if got != testCase.expectedResult {
			t.Errorf(
				"IsMatch returned an incorrect result for test case \"%s\": "+
					"expected %t, got %t",
				testCase.name,
				testCase.expectedResult,
				got,
			)
		}
	}
}

func stringifyAnything(thing any) string {
	return fmt.Sprintf("%v", thing)
}

func stringifyAPIResponse(apiResponse APIResponse) string {
	return stringifyAnything(apiResponse)
}

func stringifyLinks(links []Link) string {
	return stringifyAnything(links)
}
