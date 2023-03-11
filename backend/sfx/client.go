package sfx

import (
	"encoding/xml"
)

// SFX service URL
const DefaultSFXURL = "http://sfx.library.nyu.edu/sfxlcl41"

var sfxURL = DefaultSFXURL

func Do(request *SFXRequest) (*SFXResponse, error) {
	return request.do()
}

func SetSFXURL(dependencyInjectedURL string) {
	sfxURL = dependencyInjectedURL
}

// Validate XML, by marshalling and checking for a blank error
func isValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}
