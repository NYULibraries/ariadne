package primo

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type PrimoRequest struct {
	DumpedHTTPRequest string
	HTTPRequest       http.Request
}

func (c PrimoRequest) do() (*PrimoResponse, error) {
	client := http.Client{}
	response, err := client.Do(&c.HTTPRequest)
	if err != nil {
		return &PrimoResponse{}, fmt.Errorf("could not do request to Primo server: %v", err)
	}
	defer response.Body.Close()

	primoResponse, err := newPrimoResponse(response)
	if err != nil {
		return primoResponse, err
	}

	return primoResponse, nil
}

func NewPrimoRequest(queryString string) (*PrimoRequest, error) {
	primoRequest := &PrimoRequest{}

	queryStringValues, err := url.ParseQuery(queryString)
	if err != nil {
		return primoRequest, err
	}

	httpRequest, err := newPrimoHTTPRequest(queryStringValues)
	if err != nil {
		return primoRequest, fmt.Errorf("could not create new Primo request: %v", err)
	}
	// NOTE: This appears to drain httpRequest.Body, so when getting the dumped
	// HTTP request later, make sure to get it from primoRequest.HTTPRequest
	// and not httpRequest.
	primoRequest.HTTPRequest = (*httpRequest)

	dumpedHTTPRequest, err := httputil.DumpRequest(&primoRequest.HTTPRequest, true)
	if err != nil {
		// TODO: Log this.  PrimoRequest.DumpedHTTPRequest field is for
		// debugging only - it should not block the user request.
	}
	primoRequest.DumpedHTTPRequest = string(dumpedHTTPRequest)

	return primoRequest, nil
}

// Keeping this here even though it currently doesn't do anything in order to be
// consistent with the SFX request code.
func filterOpenURLParams(queryStringValues url.Values) url.Values {
	return queryStringValues
}

func newPrimoHTTPRequest(queryStringValues url.Values) (*http.Request, error) {
	params := filterOpenURLParams(queryStringValues)
	queryURL := primoURL + "?" + params.Encode()

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return request, fmt.Errorf("could not initialize request to Primo server: %v", err)
	}

	return request, nil
}
