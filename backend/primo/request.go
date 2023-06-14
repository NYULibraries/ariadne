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
const FRBRMemberSearchQueryParamName = "multiFacets"
const normalizedQueryParamNameISBN = "isbn"
const normalizedQueryParamNameRFTISBN = "rft.isbn"

type PrimoRequest struct {
	DumpedISBNSearchHTTPRequest string
	ISBNSearchHTTPRequest       http.Request
	QueryStringValues           url.Values
}

func (primoRequest PrimoRequest) do() (*PrimoResponse, error) {
	primoResponse := &PrimoResponse{}

	client := http.Client{}
	httpResponse, err := client.Do(&primoRequest.ISBNSearchHTTPRequest)
	if err != nil {
		return &PrimoResponse{}, fmt.Errorf("Could not do request to Primo server: %v", err)
	}
	defer httpResponse.Body.Close()

	isbnSearchResponse, err := primoResponse.addHTTPResponseData(httpResponse)
	if err != nil {
		return primoResponse, fmt.Errorf("Error adding to Primo response: %v", err)
	}

	isbn := getISBN(primoRequest.QueryStringValues)

	// Getting the links is a slightly complicated process which might require
	// additional HTTP requests to the Primo server.
	err = primoResponse.getLinks(isbn, isbnSearchResponse)
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

	primoRequest.QueryStringValues = queryStringValues

	httpRequest, err := newPrimoISBNSearchHTTPRequest(queryStringValues)
	if err != nil {
		return primoRequest, fmt.Errorf("Could not create new Primo request: %v", err)
	}
	// NOTE: This appears to drain httpRequest.Body, so when getting the dumped
	// HTTP request later, make sure to get it from primoRequest.ISBNSearchHTTPRequest
	// and not httpRequest.  If httpRequest is used later accidentally, probably
	// no harm done since currently these requests don't have a body.
	primoRequest.ISBNSearchHTTPRequest = (*httpRequest)

	dumpedISBNSearchHTTPRequest, err := httputil.DumpRequest(&primoRequest.ISBNSearchHTTPRequest, true)
	if err != nil {
		// TODO: Log this.  PrimoRequest.DumpedISBNSearchHTTPRequest field is for
		// debugging only - it should not block the user request.
	}
	primoRequest.DumpedISBNSearchHTTPRequest = string(dumpedISBNSearchHTTPRequest)

	return primoRequest, nil
}

func getISBN(queryStringValues url.Values) string {
	isbn := ""
	for queryParamName, queryParamValue := range queryStringValues {
		normalizedQueryParamName := strings.ToLower(queryParamName)
		// Prefer rft.* param name over non-rft-prefixed param name.
		if normalizedQueryParamName == normalizedQueryParamNameRFTISBN &&
			queryParamValue[0] != "" {
			isbn = queryParamValue[0]
			break
		} else if normalizedQueryParamName == normalizedQueryParamNameISBN &&
			queryParamValue[0] != "" {
			isbn = queryParamValue[0]
			break
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
func newPrimoHTTPRequest(isbn string, frbrGroupID *string) (*http.Request, error) {
	if isbn == "" {
		return nil, fmt.Errorf("query string params do not contain required ISBN param")
	}

	primoRequestParams := url.Values{
		// Same for every request
		"inst":   []string{"NYU"},
		"limit":  []string{"50"},
		"offset": []string{"0"},
		"scope":  []string{"all"},
		"vid":    []string{"NYU"},
		// ISBN query
		"q": []string{fmt.Sprintf(
			"isbn,exact,%s", isbn)},
	}

	if frbrGroupID != nil {
		primoRequestParams.Add(FRBRMemberSearchQueryParamName, fmt.Sprintf("facet_frbrgroupid,include,%s", *frbrGroupID))
	}

	queryURL := fmt.Sprintf("%s?%s", primoURL, primoRequestParams.Encode())

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return request, fmt.Errorf("Could not initialize request to Primo server: %v", err)
	}

	return request, nil
}

func newPrimoISBNSearchHTTPRequest(queryStringValues url.Values) (*http.Request, error) {
	isbn := getISBN(queryStringValues)

	return newPrimoHTTPRequest(isbn, nil)
}
