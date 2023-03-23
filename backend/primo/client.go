package primo

// Primo service URL
const DefaultPrimoURL = "https://bobcat.library.nyu.edu/primo_library/libweb/webservices/rest/primo-explore/v1/pnxs"

var primoURL = DefaultPrimoURL

func Do(request *PrimoRequest) (*PrimoResponse, error) {
	return request.do()
}

func SetPrimoURL(dependencyInjectedURL string) {
	primoURL = dependencyInjectedURL
}
