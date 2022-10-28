package sfx

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/context-objects.xml
var sfxRequestTemplate string

// Object representing everything that's needed to request from SFX
type MultipleObjectsRequest struct {
	DumpedHTTPRequest string
	HTTPRequest       http.Request
	RequestXML        string
}

// Values needed for templating an SFX request are parsed
type multipleObjectsRequestBodyParams struct {
	RftValues *map[string][]string
	Timestamp string
	Genre     string
}

func (c MultipleObjectsRequest) do() (*MultipleObjectsResponse, error) {
	client := http.Client{}
	response, err := client.Do(&c.HTTPRequest)
	if err != nil {
		return &MultipleObjectsResponse{}, fmt.Errorf("could not do post to SFX server: %v", err)
	}
	defer response.Body.Close()

	multipleObjectsResponse, err := newMultipleObjectsResponse(response)
	if err != nil {
		return multipleObjectsResponse, err
	}

	return multipleObjectsResponse, nil
}

func NewMultipleObjectsRequest(queryStringValues url.Values) (*MultipleObjectsRequest, error) {
	multipleObjectsRequest := &MultipleObjectsRequest{}

	multipleObjectsRequestBodyParams, err := parseMultipleObjectsRequestParams(queryStringValues)
	if err != nil {
		return multipleObjectsRequest, fmt.Errorf("could not parse required request body params from querystring: %v", err)
	}

	multipleObjectsRequest.RequestXML, err = requestXML(multipleObjectsRequestBodyParams)
	if err != nil {
		return multipleObjectsRequest, fmt.Errorf("could not convert multiple objects request to XML: %v", err)
	}

	httpRequest, err := newMultipleObjectsHTTPRequest(multipleObjectsRequest.RequestXML)
	if err != nil {
		return multipleObjectsRequest, fmt.Errorf("could not create new multiple objects request: %v", err)
	}
	// NOTE: This appears to drain httpRequest.Body, so when getting the dumped
	// HTTP request later, make sure to get it from multipleObjectsRequest.HTTPRequest
	// and not httpRequest.
	multipleObjectsRequest.HTTPRequest = (*httpRequest)

	dumpedHTTPRequest, err := httputil.DumpRequest(&multipleObjectsRequest.HTTPRequest, true)
	if err != nil {
		// TODO: Log this.  MultipleObjectsRequest.DumpedHTTPRequest field is for
		// debugging only - it should not block the user request.
	}
	multipleObjectsRequest.DumpedHTTPRequest = string(dumpedHTTPRequest)

	return multipleObjectsRequest, nil
}

func newMultipleObjectsHTTPRequest(requestXML string) (*http.Request, error) {
	params := url.Values{}
	params.Add("url_ctx_fmt", "info:ofi/fmt:xml:xsd:ctx")
	params.Add("sfx.response_type", "multi_obj_xml")
	// Do we always need these parameters? Umlaut adds them only in certain conditions: https://github.com/team-umlaut/umlaut/blob/master/app/service_adaptors/sfx.rb#L145-L153
	params.Add("sfx.show_availability", "1")
	params.Add("sfx.ignore_date_threshold", "1")
	params.Add("sfx.doi_url", "http://dx.doi.org")
	params.Add("url_ctx_val", requestXML)

	request, err := http.NewRequest("POST", sfxURL, strings.NewReader(params.Encode()))
	if err != nil {
		return request, fmt.Errorf("could not initialize request to SFX server: %v", err)
	}

	request.PostForm = params
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

// Parse SFX request body params from querystring.  For now, we use only fields
// prefixed with "rft.".
func parseMultipleObjectsRequestParams(queryStringValues url.Values) (multipleObjectsRequestBodyParams, error) {
	params := multipleObjectsRequestBodyParams{}

	rfts := &map[string][]string{}

	for k, v := range queryStringValues {
		// Strip the "rft." prefix from the param name and map to valid OpenURL fields
		if strings.HasPrefix(k, "rft.") {
			// E.g. "rft.book" becomes "book"
			newKey := strings.Split(k, ".")[1]
			(*rfts)[newKey] = v
		}
	}

	if reflect.DeepEqual(rfts, &map[string][]string{}) {
		return params, fmt.Errorf("no valid querystring values to parse")
	}

	genre, err := validGenre((*rfts)["genre"])
	if err != nil {
		return params, fmt.Errorf("genre is not valid: %v", err)
	}

	now := time.Now()
	params.Timestamp = now.Format(time.RFC3339Nano)
	params.RftValues = rfts
	params.Genre = genre

	return params, nil
}

func requestXML(templateValues multipleObjectsRequestBodyParams) (string, error) {
	t := template.New("sfx-request.xml").Funcs(template.FuncMap{"ToLower": strings.ToLower})

	t, err := t.Parse(sfxRequestTemplate)
	if err != nil {
		return "", fmt.Errorf("could not load template parse file: %v", err)
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, templateValues); err != nil {
		return "", fmt.Errorf("could not execute go template from multiple objects request: %v", err)
	}

	if !isValidXML(tpl.Bytes()) {
		return "", fmt.Errorf("request multiple objects XML is not valid XML: %v", err)
	}

	return tpl.String(), nil
}
