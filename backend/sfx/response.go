package sfx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

type SFXResponse struct {
	DumpedHTTPResponse string
	HTTPResponse       *http.Response
	JSON               string
	XMLResponseBody    XMLResponseBody
	XML                string
}

// Mapped out the entire Context Object responses possible from SFX as defined here:
// https://developers.exlibrisgroup.com/sfx/apis/web_services/openurl/
// But most of it is likely not useful for pulling out links of interest to us
type XMLResponseBody struct {
	ContextObject *[]ContextObject `xml:"ctx_obj" json:"ctx_obj"`
}

type ContextObject struct {
	// SFXContextObjectAttrs   string
	SFXContextObjectTargets *[]ContextObjectTargets `xml:"ctx_obj_targets" json:"ctx_obj_targets"`
}

type ContextObjectTargets struct {
	Targets *[]Target `xml:"target" json:"target"`
}

type Target struct {
	TargetName       string `xml:"target_name" json:"target_name"`
	TargetPublicName string `xml:"target_public_name" json:"target_public_name"`
	TargetUrl        string `xml:"target_url" json:"target_url"`
	Authentication   string `xml:"authentication" json:"authentication"`
	Proxy            string `xml:"proxy" json:"proxy"`
	// ObjectPortfolioId  string                  `xml:"object_portfolio_id"`
	// TargetId           string                  `xml:"target_id"`
	// TargetService_id   string                  `xml:"target_service_id"`
	// ServiceType        string                  `xml:"service_type"`
	// Parser             string                  `xml:"parser"`
	// ParseParam         string                  `xml:"parse_param"`
	// Crossref           string                  `xml:"crossref"`
	// Note               string                  `xml:"note"`
	// CharSet            string                  `xml:"char_set"`
	// Displayer          string                  `xml:"displayer"`
	// Isrelated          string                  `xml:"is_related"`
	// RelatedServiceInfo *[]RelatedServiceInfo `xml:"related_service_info"`
	Coverage *[]Coverage `xml:"coverage" json:"coverage,omitempty"`
}

// type RelatedServiceInfo struct {
// 	RelationType       string `xml:"relation_type"`
// 	RelatedObjectIssn  string `xml:"related_object_issn"`
// 	RelatedObjectTitle string `xml:"related_object_title"`
// 	RelatedObjectId    string `xml:"related_object_id"`
// }

type Coverage struct {
	CoverageText *[]CoverageText `xml:"coverage_text" json:"coverage_text,omitempty"`
	From         *[]FromTo       `xml:"from" json:"from,omitempty"`
	To           *[]FromTo       `xml:"to" json:"to,omitempty"`
	Embargo      *Embargo        `xml:"embargo" json:"embargo,omitempty"`
}

type Embargo struct {
	Availability string `xml:"availability" json:"availability,omitempty"`
	Month        string `xml:"month" json:"month,omitempty"`
	Days         string `xml:"days" json:"days,omitempty"`
}

type CoverageText struct {
	ThresholdText *[]ThresholdText    `xml:"threshold_text" json:"threshold_text,omitempty"`
	EmbargoText   *[]EmbargoStatement `xml:"embargo_text" json:"embargo_text,omitempty"`
}

type EmbargoStatement struct {
	EmbargoStatement string `xml:"embargo_statement" json:"embargo_statement,omitempty"`
}

type FromTo struct {
	Year   string `xml:"year" json:"year,omitempty"`
	Month  string `xml:"month" json:"month,omitempty"`
	Day    string `xml:"day" json:"day,omitempty"`
	Volume string `xml:"volume" json:"volume,omitempty"`
	Issue  string `xml:"issue" json:"issue,omitempty"`
}

type ThresholdText struct {
	CoverageStatement []string `xml:"coverage_statement" json:"coverage_statement,omitempty"`
}

const AskALibrarianLink = "http://library.nyu.edu/ask/"
const ILLLink = "ill.library.nyu.edu"

func (sfxResponse *SFXResponse) RemoveTarget(targetURL string) {
	currentTargets := (*(*sfxResponse.XMLResponseBody.ContextObject)[0].SFXContextObjectTargets)[0].Targets
	var newTargets []Target
	for _, target := range *currentTargets {
		if target.TargetUrl != targetURL {
			newTargets = append(newTargets, target)
		}
	}
	(*(*sfxResponse.XMLResponseBody.ContextObject)[0].SFXContextObjectTargets)[0].Targets = &newTargets
}

// TODO: Write unit test
func (sfxResponse *SFXResponse) IsFound() bool {
	targets := (*(*sfxResponse.XMLResponseBody.ContextObject)[0].SFXContextObjectTargets)[0].Targets

	// The only way to flip this to true is if a target is found that is neither
	// the Ask A Librarian link nor the ILL link.
	result := false
	if len(*targets) == 0 {
		// Theoretically this should never happen, as there is supposed to be a pre-built ILL
		// link included if no meaningful results were found.
		result = false
	} else {
		for _, target := range *targets {
			if !(strings.Contains(target.TargetUrl, ILLLink) ||
				target.TargetUrl == AskALibrarianLink) {
				result = true
				break
			}
		}
	}

	return result
}

func newSFXResponse(httpResponse *http.Response) (*SFXResponse, error) {
	// NOTE: `defer httpResponse.Body.Close()` should have already been called by the client
	// before passing to this function.

	sfxResponse := &SFXResponse{
		HTTPResponse: httpResponse,
	}

	dumpedHTTPResponse, err := httputil.DumpResponse(httpResponse, true)
	if err != nil {
		return sfxResponse, fmt.Errorf("could not dump HTTP response")
	}
	sfxResponse.DumpedHTTPResponse = string(dumpedHTTPResponse)

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return sfxResponse, fmt.Errorf("could not read response from SFX server: %v", err)
	}

	sfxResponse.XML = string(body)

	var xmlResponseBody XMLResponseBody
	if err = xml.Unmarshal(body, &xmlResponseBody); err != nil {
		return sfxResponse, err
	}

	if xmlResponseBody.ContextObject == nil {
		return sfxResponse, fmt.Errorf("could not identify context object in response")
	}

	sfxResponse.XMLResponseBody = xmlResponseBody

	json, err := json.MarshalIndent(xmlResponseBody, "", "    ")
	if err != nil {
		return sfxResponse, fmt.Errorf("could not marshal SFX response body to JSON: %v", err)
	}

	sfxResponse.JSON = string(json)

	return sfxResponse, nil
}
