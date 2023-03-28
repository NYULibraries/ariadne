package primo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"sort"
)

type Doc struct {
	Delivery struct {
		Link []Link `json:"link"`
	} `json:"delivery"`
	PNX struct {
		Facets struct {
			FRBRType    []string `json:"frbrtype"`
			FRBRGroupID []string `json:"frbrgroupid"`
		} `json:"facets"`
		Search struct {
			ISBN []string `json:"isbn"`
		} `json:"search"`
	} `json:"pnx"`
}

type Link struct {
	HyperlinkText string `json:"hyperlinkText"`
	LinkURL       string `json:"linkURL"`
	LinkType      string `json:"linkType"`
}

type PrimoResponse struct {
	DumpedFRBRMemberHTTPRequests []string
	DumpedHTTPResponses          []string
	FRBRMemberHTTPRequests       []http.Request
	HTTPResponses                []http.Response
	APIResponses                 []APIResponse
	Links                        []Link
}

type APIResponse struct {
	Docs []Doc `json:"docs"`
}

const linkToSrcType = "http://purl.org/pnx/linkType/linktorsrc"

func (primoResponse *PrimoResponse) IsFound() bool {
	return len(primoResponse.Links) > 0
}

func (primoResponse *PrimoResponse) addHTTPResponseData(httpResponse *http.Response) (APIResponse, error) {
	// NOTE: `defer httpResponse.Body.Close()` should have already been called by the client
	// before passing to this function.

	primoResponse.HTTPResponses = append(primoResponse.HTTPResponses, *httpResponse)

	dumpedHTTPResponse, err := httputil.DumpResponse(httpResponse, true)
	if err != nil {
		return APIResponse{}, fmt.Errorf("could not dump HTTP response")
	}

	primoResponse.DumpedHTTPResponses = append(primoResponse.DumpedHTTPResponses, string(dumpedHTTPResponse))

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return APIResponse{}, fmt.Errorf("could not read response from Primo server: %v", err)
	}

	var apiResponse APIResponse
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return apiResponse, err
	}

	primoResponse.APIResponses =
		append(primoResponse.APIResponses, apiResponse)

	// Returning apiResponse because httpResponse.Body has been drained and so
	// caller will not be able to get an APIResponse from it.
	return apiResponse, nil
}

func (primoResponse *PrimoResponse) addLinks(doc Doc) {
	for _, link := range doc.Delivery.Link {
		if link.LinkType == linkToSrcType {
			primoResponse.Links = append(primoResponse.Links, link)
		}
	}
}

func (primoResponse *PrimoResponse) dedupeAndSortLinks() []Link {
	processed := make(map[string]struct{})

	links := []Link{}
	for _, link := range primoResponse.Links {
		if _, ok := processed[link.LinkURL]; ok {
			continue
		}

		links = append(links, link)

		processed[link.LinkURL] = struct{}{}
	}

	sort.SliceStable(links, func(i, j int) bool { return links[i].HyperlinkText < links[j].HyperlinkText })

	return links
}

func (primoResponse *PrimoResponse) getDocsForFRBRGroup(isbn, frbrGroupID string) ([]Doc, error) {
	docs := []Doc{}

	httpRequest, err := newPrimoHTTPRequest(isbn, &frbrGroupID)
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

	apiResponse, err := primoResponse.addHTTPResponseData(httpResponse)
	if err != nil {
		return docs, fmt.Errorf("error adding to Primo response: %v", err)
	}

	return apiResponse.Docs, nil
}

func (primoResponse *PrimoResponse) getLinks(isbn string, isbnSearchResponse APIResponse) error {
	for _, doc := range isbnSearchResponse.Docs {
		if isActiveFRBRGroupType(doc) {
			// This makes another HTTP request to Primo and fetches docs for the
			// active FRBR group.
			docsForFRBRGroup, err :=
				primoResponse.getDocsForFRBRGroup(isbn, doc.PNX.Facets.FRBRGroupID[0])
			if err != nil {
				return fmt.Errorf("error fetching FRBR group links: %v", err)
			}
			// Only collect links from docs that match the user-specified ISBN.
			for _, frbrGroupDoc := range docsForFRBRGroup {
				if isMatch(frbrGroupDoc, isbn) {
					primoResponse.addLinks(frbrGroupDoc)
				}
			}
		} else {
			// No FRBR groups involved, just collect the links straight from this doc.
			primoResponse.addLinks(doc)
		}
	}

	primoResponse.dedupeAndSortLinks()

	return nil
}

func isMatch(frbrGroupDoc Doc, isbn string) bool {
	isMatch := false
	for _, isbnToTest := range frbrGroupDoc.PNX.Search.ISBN {
		if isbnToTest == isbn {
			isMatch = true
			break
		}
	}

	return isMatch
}
