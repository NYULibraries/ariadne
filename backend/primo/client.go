package primo

// Primo service URL
const DefaultPrimoURL = "http://primo.library.nyu.edu/primolcl41"

var primoURL = DefaultPrimoURL

func Do(request *PrimoRequest) (*PrimoResponse, error) {
	return request.do()
}

func SetPrimoURL(dependencyInjectedURL string) {
	primoURL = dependencyInjectedURL
}
