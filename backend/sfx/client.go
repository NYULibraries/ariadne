package sfx

import (
	"encoding/xml"
	"fmt"
)

// SFX service URL
const DefaultSFXURL = "http://sfx.library.nyu.edu/sfxlcl41"

// "Add support for OpenURLs with no genre parameter"
// https://nyu-lib.monday.com/boards/765008773/pulses/4036893558
// If genre unspecified, default genre to 'journal' https://github.com/team-umlaut/umlaut/blob/b954895e0aa0a7cd0a9ec6bb716c1886c813601e/app/service_adaptors/sfx.rb#L642
const defaultGenre = "journal"

var sfxURL = DefaultSFXURL

func Do(request *MultipleObjectsRequest) (*MultipleObjectsResponse, error) {
	return request.do()
}

func SetSFXURL(dependencyInjectedURL string) {
	sfxURL = dependencyInjectedURL
}

// Source: https://github.com/NYULibraries/umlaut/blob/105593ef5237f5b4b3de76b091c73dfc9e3e722f/config/locales/en.yml#L28-L43
func genresList() (genresList map[string]bool) {
	genresList = map[string]bool{
		"article":      true,
		"book":         true,
		"bookitem":     true,
		"conference":   true,
		"dissertation": true,
		"document":     true,
		"issue":        true,
		"journal":      true,
		"preprint":     true,
		"proceeding":   true,
		"report":       true,
		"unknown":      true,
	}

	return
}

// Validate XML, by marshalling and checking for a blank error
func isValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// Only return a valid genre that has been allowed by the OpenURL spec
func validGenre(genre []string) (string, error) {
	validGenres := genresList()
	if len(genre) > 0 && validGenres[genre[0]] {
		return genre[0], nil
	}
	return "", fmt.Errorf("genre not in list of allowed genres: %v", genre)
}
