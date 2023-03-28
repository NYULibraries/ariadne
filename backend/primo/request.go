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
		return &PrimoResponse{}, fmt.Errorf("could not do request to Primo server: %v", err)
	}
	defer httpResponse.Body.Close()

	// Note that this drains httpResponse.Body and saves it as the first element
	// of primoResponse.APIResponses as a JSON string, so when we need to parse
	// the response later we'll need to get it from primoResponse.APIResponses[0]/
	err = primoResponse.addToPrimoResponse(httpResponse)
	if err != nil {
		return primoResponse, fmt.Errorf("error adding to Primo response: %v", err)
	}

	// As mentioned above, we need to retrieve the JSON from primoResponse.APIResponses[0]
	// because httpResponse.Body has been drained.
	isbnSearchResponse := primoResponse.APIResponses[0]
	for _, doc := range isbnSearchResponse.Docs {
		if isActiveFRBRGroupType(doc) {
			// This makes another HTTP request to Primo and fetches docs for the
			// active FRBR group.
			docsForFRBRGroup, err :=
				getDocsForFRBRGroup(primoRequest.QueryStringValues, doc.PNX.Facets.FRBRGroupID[0], primoResponse)
			if err != nil {
				return primoResponse, fmt.Errorf("error fetching FRBR group links: %v", err)
			}
			// Only collect links from docs that match the user-specified ISBN.
			for _, frbrGroupDoc := range docsForFRBRGroup {
				isMatch := false
				for _, isbnToTest := range frbrGroupDoc.PNX.Search.ISBN {
					if isbnToTest == primoRequest.QueryStringValues.Get(normalizedQueryParamNameISBN) {
						isMatch = true
						break
					}
				}
				if isMatch {
					primoResponse.addLinks(frbrGroupDoc)
				}
			}
		} else {
			// No FRBR groups involved, just collect the links straight from this doc.
			primoResponse.addLinks(doc)
		}
	}

	primoResponse.dedupeAndSortLinks()

	return primoResponse, nil
}

func getDocsForFRBRGroup(queryStringValues url.Values, frbrGroupID string, primoResponse *PrimoResponse) ([]Doc, error) {
	docs := []Doc{}

	httpRequest, err := newPrimoFRBRMemberHTTPRequest(queryStringValues, &frbrGroupID)
	if err != nil {
		return docs, fmt.Errorf("could not create new FRBR group Primo request: %v", err)
	}

	// NOTE: This appears to drain httpRequest.Body, but currently these requests
	// don't have a body, so we should be okay.
	primoResponse.FRBRMemberHTTPRequests = append(primoResponse.FRBRMemberHTTPRequests, (*httpRequest))

	dumpedHTTPRequest, err := httputil.DumpRequest(httpRequest, true)
	if err != nil {
		// TODO: Log this.  PrimoRequest.DumpedISBNSearchHTTPRequest field is for
		// debugging only - it should not block the user request.
	}
	primoResponse.DumpedFRBRMemberHTTPRequests =
		append(primoResponse.DumpedFRBRMemberHTTPRequests, string(dumpedHTTPRequest))

	client := http.Client{}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return docs, fmt.Errorf("could not do FRBR group request to Primo server: %v", err)
	}
	defer httpResponse.Body.Close()

	err = primoResponse.addToPrimoResponse(httpResponse)
	if err != nil {
		return docs, fmt.Errorf("error adding to Primo response: %v", err)
	}

	// Get the unmarshalled response that primoResponse.addToPrimoResponse added.
	apiResponse := primoResponse.APIResponses[len(primoResponse.APIResponses)-1]

	return apiResponse.Docs, nil
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
		return primoRequest, fmt.Errorf("could not create new Primo request: %v", err)
	}
	// NOTE: This appears to drain httpRequest.Body, so when getting the dumped
	// HTTP request later, make sure to get it from primoRequest.ISBNSearchHTTPRequest
	// and not httpRequest.  If httpRequest is used later accidentally, probably
	// no harm done since currently these requests don't have a body.
	primoRequest.ISBNSearchHTTPRequest = (*httpRequest)

	dumpedHTTPRequest, err := httputil.DumpRequest(&primoRequest.ISBNSearchHTTPRequest, true)
	if err != nil {
		// TODO: Log this.  PrimoRequest.DumpedISBNSearchHTTPRequest field is for
		// debugging only - it should not block the user request.
	}
	primoRequest.DumpedISBNSearchHTTPRequest = string(dumpedHTTPRequest)

	return primoRequest, nil
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

func newPrimoISBNSearchHTTPRequest(queryStringValues url.Values) (*http.Request, error) {
	return newPrimoFRBRMemberHTTPRequest(queryStringValues, nil)
}

func newPrimoFRBRMemberHTTPRequest(queryStringValues url.Values, frbrGroupID *string) (*http.Request, error) {
	isbn := getISBN(queryStringValues)
	if isbn == "" {
		return nil, fmt.Errorf("query string params do not contain required ISBN param: %v", queryStringValues)
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
		return request, fmt.Errorf("could not initialize request to Primo server: %v", err)
	}

	return request, nil
}
