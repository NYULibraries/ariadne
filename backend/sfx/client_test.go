package sfx

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestIsValidXML(t *testing.T) {
	var tests = []struct {
		testXMLFile string
		expected    bool
	}{
		{"./testdata/sfx-context-object-valid.xml", true},
		{"./testdata/sfx-context-object-invalid-truncated.xml", false},
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
