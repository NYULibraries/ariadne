package sfx

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"strings"
	"text/template"
	"time"
)

var validGenres = map[string]bool{
	"journal":    true,
	"book":       true,
	"conference": true,
	"article":    true,
	"preprint":   true,
	"proceeding": true,
	"bookitem":   true,
}

type Timestamp time.Time

type ContextObject struct {
	ValidGenre bool
	RftValues  map[string][]string
	Timestamp  string
	Genre      string
}

type ContextObjectResp struct {
	XMLName xml.Name   `xml:"ctx_obj_set"`
	CtxObj  []CtxObjEl `xml:"ctx_obj"`
}

type CtxObjEl struct {
	XMLName       xml.Name          `xml:"ctx_obj"`
	CtxObjAttrs   string            //`xml:"ctx_obj_attributes"`
	CtxObjTargets []CtxObjTargetsEl `xml:"ctx_obj_targets"`
}

type CtxObjTargetsEl struct {
	XMLName xml.Name   `xml:"ctx_obj_targets"`
	Targets []TargetEl `xml:"target"`
}

type TargetEl struct {
	TargetName          string                 `xml:"target_name"`
	TargetPublicName    string                 `xml:"target_public_name"`
	ObjectPortfolioId   string                 `xml:"object_portfolio_id"`
	TargetId            string                 `xml:"target_id"`
	TargetService_id    string                 `xml:"target_service_id"`
	ServiceType         string                 `xml:"service_type"`
	Parser              string                 `xml:"parser"`
	ParseParam          string                 `xml:"parse_param"`
	Proxy               string                 `xml:"proxy"`
	Crossref            string                 `xml:"crossref"`
	Note                string                 `xml:"note"`
	Authentication      string                 `xml:"authentication"`
	CharSet             string                 `xml:"char_set"`
	Displayer           string                 `xml:"displayer"`
	TargetUrl           string                 `xml:"target_url"`
	IsRelated           string                 `xml:"is_related"`
	RelatedService_info []RelatedServiceInfoEl `xml:"related_service_info"`
	Coverage            []CoverageEl           `xml:"coverage"`
}

type RelatedServiceInfoEl struct {
	RelationType       string `xml:"relation_type"`
	RelatedObjectIssn  string `xml:"related_object_issn"`
	RelatedObjectTitle string `xml:"related_object_title"`
	RelatedObjectId    string `xml:"related_object_id"`
}

type CoverageEl struct {
	CoverageText []CoverageTextEl `xml:"coverage_text"`
	From         []FromToEl       `xml:"from"`
	To           []FromToEl       `xml:"to"`
	Embargo      string           `xml:"embargo"`
}

type CoverageTextEl struct {
	ThresholdText []ThresholdTextEl `xml:"threshold_text"`
	EmbargoText   string            `xml:"embargo_text"`
}

type FromToEl struct {
	Year   string `xml:"year"`
	Volume string `xml:"volume"`
	Issue  string `xml:"issue"`
}

type ThresholdTextEl struct {
	CoverageStatement string `xml:"coverage_statement"`
}

func validGenre(genre []string) (validGenre string, isValid bool, err error) {
	if len(genre) > 0 {
		return genre[0], validGenres[genre[0]], nil
	}
	return
}

func parseUrl(s string) (parsed map[string][]string, err error) {
	parsed = make(map[string][]string)

	u, err := url.Parse(s)
	if err != nil {
		return
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}

	// TODO: Dedupe
	for key, val := range m {
		if strings.HasPrefix(key, "rft.") {
			newKey := strings.Split(key, ".")
			parsed[newKey[1]] = val
		}
	}

	return
}

func IsValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

func (ctx *ContextObject) ToXML() (result string, err error) {
	t := template.New("index.goxml")

	t, err = t.ParseFiles("templates/index.goxml")
	if err != nil {
		return
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, ctx); err != nil {
		return
	}

	result = tpl.String()
	return
}

func CreateNewCtx(s string) (ctx *ContextObject, err error) {
	ctx = &ContextObject{}
	parsedQueryString, err := parseUrl(s)
	if err != nil {
		return
	}
	validGenre, isValidGenre, err := validGenre(parsedQueryString["genre"])
	if err != nil {
		return
	}
	now := time.Now()
	ctx.RftValues = parsedQueryString
	ctx.ValidGenre = isValidGenre
	ctx.Genre = validGenre
	ctx.Timestamp = now.Format(time.RFC3339Nano)
	return
}
