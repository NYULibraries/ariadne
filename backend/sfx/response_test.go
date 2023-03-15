package sfx

import (
	"ariadne/testutils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

func TestNewMultipleObjectsResponse(t *testing.T) {
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
			multipleObjectsResponse, err := newSFXResponse(testCase.httpResponse)
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
			if multipleObjectsResponse.JSON != testCase.expected {
				t.Errorf("multipleObjectsResponse.JSON was '%v', expecting '%v'", multipleObjectsResponse.JSON, testCase.expected)
			}
		})
	}
}

func TestRemoveTarget(t *testing.T) {
	var testCases []testutils.TestCase
	err := json.Unmarshal(testutils.TestCasesJSON, &testCases)
	if err != nil {
		panic(fmt.Sprintf("Error reading test cases file: %s", err))
	}

	//for _, testCase := range testCases {
	//	fakeHTTPResponse := makeFakeHTTPResponse("")
	//	fakeSFXResponse := makeFakeSFXResponse(fakeHTTPResponse)
	//	fakeSFXResponse.RemoveTarget("http://library.nyu.edu/ask/")
	//}
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
