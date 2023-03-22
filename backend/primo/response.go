package primo

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type PrimoResponse struct {
	DumpedHTTPResponse string
	HTTPResponse       *http.Response
}

func (primoResponse *PrimoResponse) IsFound() bool {
	return false
}

func newPrimoResponse(httpResponse *http.Response) (*PrimoResponse, error) {
	// NOTE: `defer httpResponse.Body.Close()` should have already been called by the client
	// before passing to this function.

	primoResponse := &PrimoResponse{
		HTTPResponse: httpResponse,
	}

	dumpedHTTPResponse, err := httputil.DumpResponse(httpResponse, true)
	if err != nil {
		return primoResponse, fmt.Errorf("could not dump HTTP response")
	}
	primoResponse.DumpedHTTPResponse = string(dumpedHTTPResponse)

	return primoResponse, nil
}
