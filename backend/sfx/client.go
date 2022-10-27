package sfx

import (
	"encoding/xml"
	"fmt"
)

// SFX service URL
const DefaultSFXURL = "http://sfx.library.nyu.edu/sfxlcl41"

var sfxURL = DefaultSFXURL

func Do(request *MultipleObjectsRequest) (*MultipleObjectsResponse, error) {
	return request.do()
}

func SetSFXURL(dependencyInjectedURL string) {
	sfxURL = dependencyInjectedURL
}

// A list of the valid genres as defined by the OpenURL spec
// Is this correct? See genres list on NISO spec page 59: https://groups.niso.org/higherlogic/ws/public/download/14833/z39_88_2004_r2010.pdf
func genresList() (genresList map[string]bool) {
	genresList = map[string]bool{
		"journal":    true,
		"book":       true,
		"conference": true,
		"article":    true,
		"preprint":   true,
		"proceeding": true,
		"bookitem":   true,
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
