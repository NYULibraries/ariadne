package sfx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const mockTimestamp = "2017-10-27T10:49:40-04:00"

func TestValidGenre(t *testing.T) {
	var tests = []struct {
		genre       []string
		expected    string
		expectedErr error
	}{
		{[]string{"book", "book"}, "book", nil},
		{[]string{"journal", "book"}, "journal", nil},
		{[]string{"book"}, "book", nil},
		{[]string{"wrong"}, "", errors.New("error")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.genre)
		t.Run(testname, func(t *testing.T) {
			ans, err := validGenre(tt.genre)
			if ans != tt.expected {
				t.Errorf("validGenre returned '%v', expecting '%v'", ans, tt.expected)
			}
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("validGenre err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestParseOpenURL(t *testing.T) {
	var tests = []struct {
		queryString map[string][]string
		expected    map[string][]string
		expectedErr error
	}{
		{map[string][]string{"genre": {"book"}, "rft.genre": {"book"}}, OpenURL{"genre": {"book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal", "book"}}, OpenURL{"genre": {"journal", "book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal"}}, OpenURL{"genre": {"journal"}}, nil},
		{map[string][]string{"genre": {"book"}}, OpenURL{}, errors.New("error")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.queryString)

		t.Run(testname, func(t *testing.T) {
			ans, err := parseOpenURL(tt.queryString)
			if reflect.DeepEqual(ans, tt.expected) {
				t.Errorf("parseOpenURL returned '%v', expecting '%v'", ans, tt.expected)
			}
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("parseOpenURL err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestIsValidXML(t *testing.T) {
	var tests = []struct {
		testXMLFile string
		expected    bool
	}{
		{"../../testdata/sfx-context-object-valid.xml", true},
		{"../../testdata/sfx-context-object-invalid-truncated.xml", false},
		{"", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.testXMLFile)
		t.Run(testname, func(t *testing.T) {
			data, _ := ioutil.ReadFile(tt.testXMLFile)
			ans := isValidXML(data)
			if ans != tt.expected {
				t.Errorf("isValidXML returned '%v', expecting '%v'", ans, tt.expected)
			}
		})
	}
}

func TestToRequestXML(t *testing.T) {
	var tests = []struct {
		sfxContext  *SFXContextObjectRequest
		tpl         sfxContextObjectRequestBody
		expectedErr error
	}{
		{&SFXContextObjectRequest{}, sfxContextObjectRequestBody{RftValues: &OpenURL{"genre": {"book"}, "btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
		{&SFXContextObjectRequest{}, sfxContextObjectRequestBody{}, errors.New("error")},
		{&SFXContextObjectRequest{}, sfxContextObjectRequestBody{RftValues: &OpenURL{"genre": {"<rft:"}}, Timestamp: mockTimestamp, Genre: "book"}, errors.New("error")},
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

func TestSetSFXContextObjectReq(t *testing.T) {
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
			ans, err := setSFXContextObjectReq(tt.querystring)
			// if err != nil {
			// 	t.Errorf("error %v", err)
			// }
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("setSFXContextObjectReq err was '%v', expecting '%v'", err, tt.expectedErr)
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
	dummyJSONResponse := `{"ctx_obj":[{"ctx_obj_targets":[{"target":[{"target_name":"","target_public_name":"","target_url":"http://answers.library.newschool.edu/","authentication":"","proxy":""}]}]}]}`
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
			ans, err := toResponseJson(tt.from)
			// if err != nil {
			// 	t.Errorf("error %v", err)
			// }
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("toResponseJson err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
			if ans != tt.expected {
				t.Errorf("toResponseJson was '%v', expecting '%v'", ans, tt.expected)
			}
		})
	}
}

// func (c SFXContextObjectRequest) Request() (body string, err error) {
// func Init(qs url.Values) (sfxContextObjectRequest *SFXContextObjectRequest, err error) {
