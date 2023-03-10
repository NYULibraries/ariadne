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
			ans, err := NewSFXRequest(testCase.querystring)
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
		tpl         sfxRequestBodyParams
		expectedErr error
	}{
		{"genre=\"book\"; btitle=\"a book\"", sfxRequestBodyParams{RftValues: &map[string][]string{"genre": {"book"}, "btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
		{"[empty request body]", sfxRequestBodyParams{}, errors.New("error")},
		{"genre=\"<rft:\"", sfxRequestBodyParams{RftValues: &map[string][]string{"genre": {"<rft:"}}, Timestamp: mockTimestamp, Genre: "book"}, errors.New("error")},
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

func TestFilterOpenURLParams(t *testing.T) {
	var testCases = []struct {
		testName    string
		queryString url.Values //map[string][]string
		expected    url.Values //map[string][]string
	}{
		//{"querystring contains sid", map[string][]string{"sid": {"unicode+garbage"}}, map[string][]string{"rfr_id": {"unicode+garbage"}}},
		//{"querystring doesn't contain sid", map[string][]string{"id": {"unicode+garbage"}}, map[string][]string{"id": {"unicode+garbage"}}},
		{"querystring contains sid", map[string][]string{"sid": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}, map[string][]string{"rfr_id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}},
		{"querystring doesn't contain sid", map[string][]string{"id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}, map[string][]string{"id": {"unicode+garbage+EBSCO:Scopus\\u00ae"}}},
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

func TestParseMultipleObjectsRequestParams(t *testing.T) {
	var testCases = []struct {
		queryString   map[string][]string
		expected      *map[string][]string
		expectedError error
	}{
		// Even though "rft." prefix is supposed to be required, many past requests
		// have had non-prefixed equivalent params, so we've decided to allow them.
		{map[string][]string{"genre": {"book"}}, &map[string][]string{"genre": {"book"}}, nil},

		// "rft."-prefixed query params should always have priority over their non-prefixed equivalents.
		// When we first made the change from requiring the prefix to not requiring it,
		// we had a bug where priority was determined by ordering, so we test opposite
		// orderings to prevent regression.
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal", "book"}}, &map[string][]string{"genre": {"journal", "book"}}, nil},
		{map[string][]string{"rft.genre": {"journal", "book"}, "genre": {"book"}}, &map[string][]string{"genre": {"journal", "book"}}, nil},

		{map[string][]string{"genre": {"podcast"}}, nil, errors.New("error")},
	}

	for _, testCase := range testCases {
		testname := fmt.Sprintf("%s", testCase.queryString)

		t.Run(testname, func(t *testing.T) {
			params, err := parseRequestParams(testCase.queryString)
			actual := params.RftValues
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("parseRequestParams returned '%v', expecting '%v'", actual, testCase.expected)
			}
			if testCase.expectedError != nil {
				if err == nil {
					t.Errorf("parseRequestParams err was '%v', expecting '%v'", err, testCase.expectedError)
				}
			}
		})
	}
}

// func (c SFXRequest) Do() (body string, err error) {
// func Init(qs url.Values) (SFXRequest *SFXRequest, err error) {
