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
		Pnx  struct {
			Facets struct {
				Frbrtype    []string `json:"frbrtype"`
				Frbrgroupid []string `json:"frbrgroupid"`
			} `json:"facets"`
			Search struct {
				Isbn []string `json:"isbn"`
			} `json:"search"`
		} `json:"pnx"`
	}
}

type Link struct {
	HyperlinkText string `json:"hyperlinkText"`
	LinkURL       string `json:"linkURL"`
	LinkType      string `json:"linkType"`
}

type PrimoResponse struct {
	DumpedHTTPResponses []string
	HTTPResponses       []http.Response
	APIResponses        []APIResponse
	Links               []Link
}

type APIResponse struct {
	Docs []Doc `json:"docs"`
}

const linkToSrcType = "http://purl.org/pnx/linkType/linktorsrc"

func (primoResponse *PrimoResponse) IsFound() bool {
	return false
}

func (primoResponse *PrimoResponse) addToPrimoResponse(httpResponse *http.Response) error {
	// NOTE: `defer httpResponse.Body.Close()` should have already been called by the client
	// before passing to this function.

	primoResponse.HTTPResponses = append(primoResponse.HTTPResponses, *httpResponse)

	dumpedHTTPResponse, err := httputil.DumpResponse(httpResponse, true)
	if err != nil {
		return fmt.Errorf("could not dump HTTP response")
	}

	primoResponse.DumpedHTTPResponses = append(primoResponse.DumpedHTTPResponses, string(dumpedHTTPResponse))

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return fmt.Errorf("could not read response from Primo server: %v", err)
	}

	var primoSearchAPIResponse APIResponse
	if err = json.Unmarshal(body, &primoSearchAPIResponse); err != nil {
		return err
	}

	primoResponse.APIResponses =
		append(primoResponse.APIResponses, primoSearchAPIResponse)

	return nil
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
