package sfx

import (
	"errors"
	"fmt"
	"testing"
)

func TestValidGenre(t *testing.T) {
	expectedErr := errors.New("genre not in list of allowed genres")
	var tests = []struct {
		genre       []string
		expected    string
		expectedErr error
	}{
		{[]string{"book", "book"}, "book", nil},
		{[]string{"journal", "book"}, "journal", nil},
		{[]string{"book"}, "book", nil},
		{[]string{"wrong"}, "", expectedErr},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.genre)
		t.Run(testname, func(t *testing.T) {
			ans, _ := validGenre(tt.genre)
			if ans != tt.expected {
				t.Errorf("validGenre returned '%v', expecting '%v'", ans, tt.expected)
			}
			// if err.Error() != tt.expectedErr.Error() {
			// 	t.Errorf("validGenre err was '%v', expecting '%v'", err, tt.expectedErr)
			// }
		})
	}
}

func TestParseOpenURL(t *testing.T) {
	var tests = []struct {
		queryString map[string][]string
		expected    string
		expectedErr error
	}{
		{map[string][]string{"genre": ["book"], "rft.genre": "book"}, "book", nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.genre)
		t.Run(testname, func(t *testing.T) {
			ans, _ := parseOpenURL(tt.queryString)
			if ans != tt.expected {
				t.Errorf("parseOpenURL returned '%v', expecting '%v'", ans, tt.expected)
			}
			// if err.Error() != tt.expectedErr.Error() {
			// 	t.Errorf("validGenre err was '%v', expecting '%v'", err, tt.expectedErr)
			// }
		})
	}
}
