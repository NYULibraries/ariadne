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

func stringifyAPIResponse(apiResponse APIResponse) string {
	return fmt.Sprintf("%v", apiResponse)
}
