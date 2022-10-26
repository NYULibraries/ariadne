package sfx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type MultipleObjectsResponse struct {
	DumpedHTTPResponse      string
	HTTPResponse            *http.Response
	JSON                    string
	MultiObjXMLResponseBody MultiObjXMLResponseBody
	XML                     string
}

// Mapped out the entire Context Object responses possible from SFX as defined here:
// https://developers.exlibrisgroup.com/sfx/apis/web_services/openurl/
// But most of it is likely not useful for pulling out links of interest to us
type MultiObjXMLResponseBody struct {
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

func newMultipleObjectsResponse(httpResponse *http.Response) (*MultipleObjectsResponse, error) {
	// NOTE: `defer httpResponse.Body.Close()` should have already been called by the client
	// before passing to this function.

	multipleObjectsResponse := &MultipleObjectsResponse{
		HTTPResponse: httpResponse,
	}

	dumpedHTTPResponse, err := httputil.DumpResponse(httpResponse, true)
	if err != nil {
		return multipleObjectsResponse, fmt.Errorf("could not dump HTTP response")
	}
	multipleObjectsResponse.DumpedHTTPResponse = string(dumpedHTTPResponse)

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return multipleObjectsResponse, fmt.Errorf("could not read response from SFX server: %v", err)
	}

	multipleObjectsResponse.XML = string(body)

	var multiObjXMLResponseBody MultiObjXMLResponseBody
	if err = xml.Unmarshal(body, &multiObjXMLResponseBody); err != nil {
		return multipleObjectsResponse, err
	}

	json, err := json.MarshalIndent(multiObjXMLResponseBody, "", "    ")
	if err != nil {
		return multipleObjectsResponse, fmt.Errorf("could not marshal SFX response body to JSON: %v", err)
	}

	multipleObjectsResponse.JSON = string(json)

	return multipleObjectsResponse, nil
}
