package sfx

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

const mockTimestamp = "2017-10-27T10:49:40-04:00"

func TestNewMultipleObjectsRequest(t *testing.T) {
	var tests = []struct {
		querystring url.Values
		expectedErr error
	}{
		{map[string][]string{"genre": {"book"}}, errors.New("error")},
		{map[string][]string{"rft.genre": {"podcast"}}, errors.New("error")},
		{map[string][]string{"rft.genre": {"book"}, "rft.aulast": {"<rft:"}}, nil},
		{map[string][]string{"rft.genre": {"book"}, "rft.btitle": {"dune"}}, nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.querystring)
		t.Run(testname, func(t *testing.T) {
			ans, err := NewMultipleObjectsRequest(tt.querystring)
			// if err != nil {
			// 	t.Errorf("error %v", err)
			// }
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("setMultipleObjectsRequest err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
			if err == nil {
				if !strings.HasPrefix(ans.RequestXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
					t.Errorf("requestXML isn't an XML document")
				}
			}
		})
	}
}

func TestToRequestXML(t *testing.T) {
	var tests = []struct {
		sfxContext  *MultipleObjectsRequest
		tpl         multipleObjectsRequestBody
		expectedErr error
	}{
		{&MultipleObjectsRequest{}, multipleObjectsRequestBody{RftValues: &openURL{"genre": {"book"}, "btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
		{&MultipleObjectsRequest{}, multipleObjectsRequestBody{}, errors.New("error")},
		{&MultipleObjectsRequest{}, multipleObjectsRequestBody{RftValues: &openURL{"genre": {"<rft:"}}, Timestamp: mockTimestamp, Genre: "book"}, errors.New("error")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.sfxContext)
		t.Run(testname, func(t *testing.T) {
			c := tt.sfxContext
			err := c.toRequestXML(tt.tpl)
			if tt.expectedErr == nil && !strings.HasPrefix(c.RequestXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
				t.Errorf("toRequestXML didn't return an XML document")
			}
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("toRequestXML err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestToResponseJson(t *testing.T) {
	dummyGoodXMLResponse := `
<ctx_obj_set>
	<ctx_obj>
		<ctx_obj_targets>
			<target>
				<target_url>http://answers.library.newschool.edu/</target_url>
			</target>
		</ctx_obj_targets>
	</ctx_obj>
</ctx_obj_set>`
	dummyBadXMLResponse := `
<ctx_obj_set`
	dummyJSONResponse := `{
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
	var tests = []struct {
		from        []byte
		expected    string
		expectedErr error
	}{
		{[]byte(dummyGoodXMLResponse), dummyJSONResponse, nil},
		{[]byte(dummyBadXMLResponse), "", errors.New("error")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.expected)
		t.Run(testname, func(t *testing.T) {
			ans, err := toResponseJSON(tt.from)
			// if err != nil {
			// 	t.Errorf("error %v", err)
			// }
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("toResponseJSON err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
			if ans != tt.expected {
				t.Errorf("toResponseJSON was '%v', expecting '%v'", ans, tt.expected)
			}
		})
	}
}

// func (c MultipleObjectsRequest) Do() (body string, err error) {
// func Init(qs url.Values) (MultipleObjectsRequest *MultipleObjectsRequest, err error) {
