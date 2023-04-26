package sfx

import (
	"ariadne/testutils"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

const dummyGoodXMLResponse = `
<ctx_obj_set>
	<ctx_obj>
		<ctx_obj_targets>
			<target>
				<target_url>http://answers.library.newschool.edu/</target_url>
			</target>
		</ctx_obj_targets>
	</ctx_obj>
</ctx_obj_set>`

const dummyBadXMLResponse = `
<ctx_obj_set`

const dummyErrorXMLResponse = `
37
<html><body><p>XSS violation occured.</p></body></html>
0
`

const dummyJSONResponse = `{
    "ctx_obj": [
        {
            "ctx_obj_targets": [
                {
                    "target": [
                        {
                            "target_name": "",
                            "target_public_name": "",
                            "target_url": "http://answers.library.newschool.edu/",
                            "authentication": "",
                            "proxy": ""
                        }
                    ]
                }
            ]
        }
    ]
}`

func TestNewSFXResponse(t *testing.T) {
	testCases := []struct {
		name          string
		httpResponse  *http.Response
		expected      string
		expectedError error
	}{
		{"Good SFX response", makeFakeHTTPResponse(dummyGoodXMLResponse), dummyJSONResponse, nil},
		{"Bad SFX response", makeFakeHTTPResponse(dummyBadXMLResponse), "", errors.New("XML syntax error on line 2: unexpected EOF")},
		{"Error SFX response", makeFakeHTTPResponse(dummyErrorXMLResponse), "", errors.New("could not identify context object in response")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sfxResponse, err := newSFXResponse(testCase.httpResponse)
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("newSFXResponse returned no error, expecting '%v'", testCase.expectedError)
				}
				if err.Error() != testCase.expectedError.Error() {
					t.Errorf("newSFXResponse returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("newSFXResponse returned error '%v', expecting no errors", err)
			}
			if sfxResponse.JSON != testCase.expected {
				t.Errorf("sfxResponse.JSON was '%v', expecting '%v'", sfxResponse.JSON, testCase.expected)
			}
		})
	}
}

func TestRemoveTarget(t *testing.T) {
	targetURLsToRemove := []string{
		"http://library.nyu.edu/ask/",
		"http://proxy.library.nyu.edu/login\\\\?url=http://www.newyorker.com/archive",
	}

	for _, testCase := range testutils.TestCases {
		for _, targetURLToRemove := range targetURLsToRemove {
			sfxFakeResponseFixture, err := testutils.GetSFXFakeResponse(testCase)
			if err != nil {
				t.Fatal(err)
			}

			// This is an error test case, where no SFX response is ever received
			// due to failure to construct or execute a request to the SFX API.
			// This test only makes sense with a valid SFX response.
			if sfxFakeResponseFixture == "" {
				continue
			}

			fakeSFXResponse := makeFakeSFXResponse(makeFakeHTTPResponse(sfxFakeResponseFixture))

			// Generate stringified expected targets slice
			originalTargetsStringified := fmt.Sprintf("%v", (*(*fakeSFXResponse.XMLResponseBody.ContextObject)[0].SFXContextObjectTargets)[0].Targets)
			expectedTargetsStringified := getExpectedTargetsStringified(originalTargetsStringified, targetURLToRemove)

			fakeSFXResponse.RemoveTarget(targetURLToRemove)

			// Generate stringified targets slice that should now have `targetURLToRemove` removed
			newTargetsStringified := fmt.Sprintf("%v", (*(*fakeSFXResponse.XMLResponseBody.ContextObject)[0].SFXContextObjectTargets)[0].Targets)

			if newTargetsStringified != expectedTargetsStringified {
				t.Errorf(
					fmt.Sprintf(
						"%s: expected %s; got %s",
						testCase.Name,
						expectedTargetsStringified,
						newTargetsStringified,
					),
				)
			}
		}
	}
}

func getExpectedTargetsStringified(targetsStringified string, targetURLToRemove string) string {
	// Example string matched by this regexp: {ASK_A_LIBRARIAN_LCL Ask a Librarian http://library.nyu.edu/ask/  no 0xc00000ed50}
	targetURLRegexp := regexp.MustCompile("{[^}]*" + targetURLToRemove + "[^}]*}")
	expectedTargetsStringified := targetURLRegexp.ReplaceAllString(targetsStringified, "")

	// Clean up whitespace
	// ...when removed target came after another target
	expectedTargetsStringified = strings.ReplaceAll(expectedTargetsStringified, "}  ", "} ")
	// ...when removed target came before another target
	expectedTargetsStringified = strings.ReplaceAll(expectedTargetsStringified, "  {", " {")
	// ...when removed target was at the end of the slice
	expectedTargetsStringified = strings.ReplaceAll(expectedTargetsStringified, " ]", "]")
	// ...when removed target was at the beginning of the slice
	expectedTargetsStringified = strings.ReplaceAll(expectedTargetsStringified, "[ ", "[")

	return expectedTargetsStringified
}

func makeFakeHTTPResponse(body string) *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}

func makeFakeSFXResponse(httpResponse *http.Response) *SFXResponse {
	fakeSFXResponse, _ := newSFXResponse(httpResponse)

	return fakeSFXResponse
}
