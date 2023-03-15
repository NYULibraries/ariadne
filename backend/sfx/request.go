package sfx

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type SFXRequest struct {
	DumpedHTTPRequest string
	HTTPRequest       http.Request
}

func (c SFXRequest) do() (*SFXResponse, error) {
	client := http.Client{}
	response, err := client.Do(&c.HTTPRequest)
	if err != nil {
		return &SFXResponse{}, fmt.Errorf("could not do request to SFX server: %v", err)
	}
	defer response.Body.Close()

	sfxResponse, err := newSFXResponse(response)
	if err != nil {
		return sfxResponse, err
	}

	return sfxResponse, nil
}

func NewSFXRequest(queryString string) (*SFXRequest, error) {
	sfxRequest := &SFXRequest{}

	queryStringValues, err := url.ParseQuery(queryString)
	if err != nil {
		return sfxRequest, err
	}

	httpRequest, err := newSFXHTTPRequest(queryStringValues)
	if err != nil {
		return sfxRequest, fmt.Errorf("could not create new SFX request: %v", err)
	}
	// NOTE: This appears to drain httpRequest.Body, so when getting the dumped
	// HTTP request later, make sure to get it from sfxRequest.HTTPRequest
	// and not httpRequest.
	sfxRequest.HTTPRequest = (*httpRequest)

	dumpedHTTPRequest, err := httputil.DumpRequest(&sfxRequest.HTTPRequest, true)
	if err != nil {
		// TODO: Log this.  SFXRequest.DumpedHTTPRequest field is for
		// debugging only - it should not block the user request.
	}
	sfxRequest.DumpedHTTPRequest = string(dumpedHTTPRequest)

	return sfxRequest, nil
}

// Transforms `sid` to `rfr_id` since the former seems to trigger SFX errors when
// its value contains certain unicode encodings, while the latter doesn't.
// An example of such a request is the test case "Trouble-causing `sid`" in the
// `TestNewSFXRequest` unit test.
// Can't include the URL example in this comment because it appears to cause
// the `gofmt -l ...` test to fail in CircleCI, which currently uses Go version 1.19.1 linux/amd64.
// `gofmt -l` for this file does not yield the same result on a Mac.
func filterOpenURLParams(queryStringValues url.Values) url.Values {
	// If no `sid`, we do nothing
	sid := queryStringValues.Get("sid")
	if sid == "" {
		return queryStringValues
	}

	// Replace `sid` with `rfr_id`
	queryStringValues.Del("sid")
	queryStringValues.Add("rfr_id", sid)

	return queryStringValues
}

func newSFXHTTPRequest(queryStringValues url.Values) (*http.Request, error) {
	params := filterOpenURLParams(queryStringValues)

	// Add SFX query params
	params.Add("url_ctx_fmt", "info:ofi/fmt:xml:xsd:ctx")
	params.Add("sfx.response_type", "multi_obj_xml")
	// Do we always need these parameters? Umlaut adds them only in certain conditions: https://github.com/team-umlaut/umlaut/blob/b954895e0aa0a7cd0a9ec6bb716c1886c813601e/app/service_adaptors/sfx.rb#L145-L153
	if !(params.Has("date") || params.Has("rft.date") || params.Has("rft.year") || params.Has("year")) {
		params.Add("sfx.show_availability", "1")
		params.Add("sfx.ignore_date_threshold", "1")
	}
	params.Add("sfx.doi_url", "http://dx.doi.org")

	queryURL := sfxURL + "?" + params.Encode()
	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return request, fmt.Errorf("could not initialize request to SFX server: %v", err)
	}

	return request, nil
}
