package sfx

// SFX service URL
const DefaultSFXURL = "http://sfx.library.nyu.edu/sfxlcl41"

var sfxURL = DefaultSFXURL

func Do(request *SFXRequest) (*SFXResponse, error) {
	return request.do()
}

func SetSFXURL(dependencyInjectedURL string) {
	sfxURL = dependencyInjectedURL
}
