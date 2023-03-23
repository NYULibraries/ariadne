package primo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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
	DisplayLabel string `json:"displayLabel"`
	LinkURL      string `json:"linkURL"`
	LinkType     string `json:"linkType"`
}

type PrimoResponse struct {
	DumpedHTTPResponses     []string
	HTTPResponses           []http.Response
	PrimoSearchAPIResponses []PrimoSearchAPIResponse
	Links                   []Link
}

type PrimoSearchAPIResponse struct {
	Docs []Doc `json:"docs"`
}

func (primoResponse *PrimoResponse) IsFound() bool {
	return false
}

func addToPrimoResponse(primoResponse *PrimoResponse, httpResponse *http.Response) error {
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

	var primoSearchAPIResponse PrimoSearchAPIResponse
	if err = json.Unmarshal(body, &primoSearchAPIResponse); err != nil {
		return err
	}

	primoResponse.PrimoSearchAPIResponses =
		append(primoResponse.PrimoSearchAPIResponses, primoSearchAPIResponse)

	return nil
}
