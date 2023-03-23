package primo

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const activeFRBRGroupType = "5"
const normalizedQueryParamNameISBN = "isbn"

var primoDefaultRequestParams = url.Values{
	"inst":   []string{"NYU"},
	"limit":  []string{"50"},
	"offset": []string{"0"},
	"scope":  []string{"all"},
	"vid":    []string{"NYU"},
}

type PrimoRequest struct {
	DumpedHTTPRequest string
	HTTPRequest       http.Request
	QueryStringValues url.Values
}

func (primoRequest PrimoRequest) do() (*PrimoResponse, error) {
	primoResponse := &PrimoResponse{}

	client := http.Client{}
	httpResponse, err := client.Do(&primoRequest.HTTPRequest)
	if err != nil {
		return &PrimoResponse{}, fmt.Errorf("could not do request to Primo server: %v", err)
	}
	defer httpResponse.Body.Close()

	err = primoResponse.addToPrimoResponse(httpResponse)
	if err != nil {
		return primoResponse, fmt.Errorf("error added to Primo response: %v", err)
	}

	isbnSearchResponse := primoResponse.APIResponses[0]
	for _, doc := range isbnSearchResponse.Docs {
		if isActiveFRBRGroupType(doc) {
			// TODO: recursively collect links
		} else {
			primoResponse.addLinks(doc)
		}
	}

	primoResponse.dedupeAndSortLinks()

	return primoResponse, nil
}

func NewPrimoRequest(queryString string) (*PrimoRequest, error) {
	primoRequest := &PrimoRequest{}

	queryStringValues, err := url.ParseQuery(queryString)
	if err != nil {
		return primoRequest, err
	}

	primoRequest.QueryStringValues = queryStringValues

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

// Keep this for now even though it doesn't do anything, because it matches SFX
// processing and because we might do filtering later.
func filterOpenURLParams(queryStringValues url.Values) url.Values {
	return queryStringValues
}

func getISBN(queryStringValues url.Values) string {
	isbn := ""
	for queryParamName, queryParamValue := range queryStringValues {
		normalizedQueryParamName := strings.ToLower(queryParamName)
		if normalizedQueryParamName == normalizedQueryParamNameISBN {
			isbn = queryParamValue[0]
		}
	}

	return isbn
}

func isActiveFRBRGroupType(doc Doc) bool {
	result := false

	for _, frbrType := range doc.PNX.Facets.FRBRType {
		if frbrType == activeFRBRGroupType {
			result = true
			break
		}
	}

	return result
}

func newPrimoHTTPRequest(queryStringValues url.Values) (*http.Request, error) {
	return newPrimoHTTPRequestFRBR(queryStringValues, nil)
}

func newPrimoHTTPRequestFRBR(queryStringValues url.Values, frbrGroupId *string) (*http.Request, error) {
	params := filterOpenURLParams(queryStringValues)

	isbn := getISBN(queryStringValues)
	if isbn == "" {
		return nil, fmt.Errorf("query string params do not contain required ISBN param: %v", queryStringValues)
	}

	primoRequestParams := primoDefaultRequestParams
	primoRequestParams.Add("q", fmt.Sprintf(
		"isbn,exact,%s", isbn))

	if frbrGroupId != nil {
		params.Add("multiFacets", fmt.Sprintf("facet_frbrgroupid,include%s", frbrGroupId))
	}

	queryURL := fmt.Sprintf("%s?%s", primoURL, primoRequestParams.Encode())

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return request, fmt.Errorf("could not initialize request to Primo server: %v", err)
	}

	return request, nil
}
