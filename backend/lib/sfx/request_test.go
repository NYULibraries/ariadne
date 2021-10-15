package sfx

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
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
		{map[string][]string{"genre": {"book"}, "rft.genre": {"book"}}, map[string][]string{"rft:genre": {"book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal", "book"}}, map[string][]string{"rft:genre": {"journal", "book"}}, nil},
		{map[string][]string{"genre": {"book"}, "rft.genre": {"journal"}}, map[string][]string{"rft:genre": {"journal"}}, nil},
		{map[string][]string{"genre": {"book"}}, map[string][]string{}, errors.New("error")},
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
		{"../../testdata/ctxObj_good.xml", true},
		{"../../testdata/ctxObj_bad.xml", false},
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

func TestToXML(t *testing.T) {
	var tests = []struct {
		ctx         *ContextObjectReq
		expectedErr error
	}{
		{&ContextObjectReq{RftValues: map[string][]string{"rft:genre": {"book"}, "rft:btitle": {"a book"}}, Timestamp: mockTimestamp, Genre: "book"}, nil},
	}

	// Create the templates/index.goxml in the current test context temporarily
	// and delete after the test completes
	err := os.Mkdir("templates", 0755)
	if err != nil {
		t.Errorf("could not create temp templates dir")
	}
	_, err = copy("../../templates/index.goxml", "./templates/index.goxml")
	if err != nil {
		t.Errorf("could not copy template file")
	}
	defer os.RemoveAll("templates")

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.ctx)
		t.Run(testname, func(t *testing.T) {
			c := tt.ctx
			ans, err := c.toXML()
			if !strings.HasPrefix(ans, `<?xml version="1.0" encoding="UTF-8"?>`) {
				t.Errorf("toXML didn't return an XML document")
			}
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("toXML err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestSetContextObjectReq(t *testing.T) {
	var tests = []struct {
		querystring   url.Values
		expectedGenre string
		expectedRfts  map[string][]string
		expectedErr   error
	}{
		{map[string][]string{"genre": {"book"}, "rft.genre": {"book"}, "rft.btitle": {"a book"}}, "book", map[string][]string{"rft:genre": {"book"}, "rft:btitle": {"a book"}}, nil},
		{map[string][]string{}, "book", map[string][]string{"rft:genre": {"book"}, "rft:btitle": {"a book"}}, nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.querystring)
		t.Run(testname, func(t *testing.T) {
			ans, err := setContextObjectReq(tt.querystring)
			fmt.Printf("this is the obj %v", err)
			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("setContextObjectReq err was '%v', expecting '%v'", err, tt.expectedErr)
				}
			}
			if err == nil {
				if ans.Genre != tt.expectedGenre {
					t.Errorf("setContextObjectReq.Genre returned '%v', expecting '%v'", ans, tt.expectedGenre)
				}
				if reflect.DeepEqual(ans.RftValues, tt.expectedRfts) {
					t.Errorf("setContextObjectReq.RftValues returned '%v', expecting '%v'", ans, tt.expectedRfts)
				}
			}
		})
	}
}

// type ContextObjectReq struct {
// 	RftValues map[string][]string
// 	Timestamp string
// 	Genre     string
// }

// func setContextObjectReq(qs url.Values) (ctx *ContextObjectReq, err error) {
// func toJson(from []byte) (to string, err error) {
// func ToCtxObjReq(qs url.Values) (ctxObjReqXml string, err error) {
// func Post(xmlBody string) (body string, err error) {

// Util function for copying a file from a source to a new dest
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
