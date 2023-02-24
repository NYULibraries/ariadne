package sfx

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const mockTimestamp = "2017-10-27T10:49:40-04:00"

func TestNewMultipleObjectsRequest(t *testing.T) {
	var tests = []struct {
		querystring   url.Values
		expectedError error
	}{
		{map[string][]string{}, errors.New("could not parse required request body params from querystring: no valid querystring values to parse")},
		{map[string][]string{"rft.genre": {"podcast"}}, errors.New("could not parse required request body params from querystring: genre is not valid: genre not in list of allowed genres: [podcast]")},
		// "<" should be XML escaped properly
		{map[string][]string{"rft.genre": {"book"}, "rft.aulast": {"<rft:"}}, nil},
		// "&" should be XML escaped properly
		{map[string][]string{"rft.genre": {"journal"}, "title": {"Journal of the Gilded Age & Progressive Era"}}, nil},
		{map[string][]string{"rft.genre": {"book"}, "rft.btitle": {"dune"}}, nil},
		// If genre is missing, automatically add genre param with value of "journal"
		{map[string][]string{"title": {"Journal of the Gilded Age & Progressive Era"}}, nil},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("%s", testCase.querystring)
		t.Run(testName, func(t *testing.T) {
			ans, err := NewMultipleObjectsRequest(testCase.querystring)
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("NewMultipleObjectsRequest returned no error, expecting '%v'", testCase.expectedError)
				} else if err.Error() != testCase.expectedError.Error() {
					t.Errorf("NewMultipleObjectsRequest returned error '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
			if err != nil && testCase.expectedError == nil {
				t.Errorf("NewMultipleObjectsRequest returned error '%v', expecting no errors", err)
			}
			if err == nil {
				if !strings.HasPrefix(ans.RequestXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
					t.Errorf("requestXML isn't an XML document")
				}
			}
		})
	}
}

func TestRequestXML(t *testing.T) {
	var tests = []struct {
		name        string
		tpl         multipleObjectsRequestBodyParams
		expectedErr error
	}{
		{"genre=\"book\"; btitle=\"a book\"", multipleObjectsRequestBodyParams{RftValues: &map[string][]string{"genre": {"book"}, "btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
		{"[empty request body]", multipleObjectsRequestBodyParams{}, errors.New("error")},
		{"genre=\"<rft:\"", multipleObjectsRequestBodyParams{RftValues: &map[string][]string{"genre": {"<rft:"}}, Timestamp: mockTimestamp, Genre: "book"}, errors.New("error")},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actualXML, err := requestXML(testCase.tpl)
			if testCase.expectedErr == nil && !strings.HasPrefix(actualXML, `<?xml version="1.0" encoding="UTF-8"?>`) {
				t.Errorf("toRequestXML didn't return an XML document")
			}
			if testCase.expectedErr != nil {
				if err == nil {
					t.Errorf("toRequestXML err was '%v', expecting '%v'", err, testCase.expectedErr)
				}
			}
		})
	}
}

func TestParseMultipleObjectsRequestParams(t *testing.T) {
	var testCases = []struct {
		queryString   map[string][]string
		expected      *map[string][]string
		expectedError error
	}{
		{map[string][]string{"genre": {"book"}, "rft.genre": {"book"}}, &map[string][]string{"genre": {"book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal", "book"}}, &map[string][]string{"genre": {"journal", "book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal"}}, &map[string][]string{"genre": {"journal"}}, nil},
		{map[string][]string{"genre": {"book"}}, &map[string][]string{"genre": {"book"}}, nil},
		{map[string][]string{"genre": {"podcast"}}, nil, errors.New("error")},
	}

	for _, testCase := range testCases {
		testname := fmt.Sprintf("%s", testCase.queryString)

		t.Run(testname, func(t *testing.T) {
			params, err := parseMultipleObjectsRequestParams(testCase.queryString)
			actual := params.RftValues
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("parseMultipleObjectsRequestParams returned '%v', expecting '%v'", actual, testCase.expected)
			}
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("parseMultipleObjectsRequestParams err was '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
		})
	}
}

// func (c MultipleObjectsRequest) Do() (body string, err error) {
// func Init(qs url.Values) (MultipleObjectsRequest *MultipleObjectsRequest, err error) {
