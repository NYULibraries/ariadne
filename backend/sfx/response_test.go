package sfx

import (
	"bytes"
	"errors"
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			multipleObjectsResponse, err := newMultipleObjectsResponse(testCase.httpResponse)
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("newMultipleObjectsResponse returned no error, expecting '%v'", testCase.expectedError)
				}
				if err.Error() != testCase.expectedError.Error() {
					t.Errorf("newMultipleObjectsResponse returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("newMultipleObjectsResponse returned error '%v', expecting no errors", err)
			}
			if multipleObjectsResponse.JSON != testCase.expected {
				t.Errorf("multipleObjectsResponse.JSON was '%v', expecting '%v'", multipleObjectsResponse.JSON, testCase.expected)
			}
		})
	}
}

func makeFakeHTTPResponse(body string) *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}
