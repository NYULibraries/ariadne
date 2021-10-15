package sfx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

// A list of the valid genres as defined by the OpenURL spec
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

// Values needed for an SFX request are parsed
type ContextObjectReq struct {
	RftValues map[string][]string
	Timestamp string
	Genre     string
}

// This is the NYU SFX API url
const sfxUrl string = "http://sfx.library.nyu.edu/sfxlcl41"

// Only return a valid genre that has been allowed by the OpenURL spec
func validGenre(genre []string) (string, error) {
	if len(genre) > 0 && validGenres[genre[0]] {
		return genre[0], nil
	}
	return "", fmt.Errorf("genre not in list of allowed genres")
}

// Take an openurl and return a map of only the rft-prefixed fields
// These are the fields we are going to parse into XML as part of the
// post request params
func parseOpenURL(qs url.Values) (parsed map[string][]string, err error) {
	parsed = make(map[string][]string)
	for key, val := range qs {
		if strings.HasPrefix(key, "rft.") {
			// TODO: Dedupe
			newKey := strings.Split(key, ".")
			parsed[newKey[1]] = val
		}
	}

	return
}

func isValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// Convert a context object request to an XML string
// via gotemplates, in order to set it up as a post param to SFX
func (c *ContextObjectReq) toXML() (result string, err error) {
	t := template.New("index.goxml")

	t, err = t.ParseFiles("templates/index.goxml")
	if err != nil {
		return result, fmt.Errorf("could not load template parse file: %v", err)
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, c); err != nil {
		return result, fmt.Errorf("could not execute go template from context object request: %v", err)
	}

	result = tpl.String()
	return
}

// Setup the ContextObjectReq instance we'll need to run with
// the gotemplates to create the valid XML string param
func setContextObjectReq(qs url.Values) (ctx *ContextObjectReq, err error) {
	rfts, err := parseOpenURL(qs)
	if err != nil {
		return ctx, fmt.Errorf("could not parse OpenURL: %v", err)
	}

	validGenre, err := validGenre(qs["genre"])
	if err != nil {
		return ctx, fmt.Errorf("genre is not valid: %v", err)
	}

	now := time.Now()
	ctx = &ContextObjectReq{
		Timestamp: now.Format(time.RFC3339Nano),
		RftValues: rfts,
		Genre:     validGenre,
	}

	return
}

// Convert the response XML from SFX into a JSON string
func toJson(from []byte) (to string, err error) {
	var p CtxObjSet
	if err = xml.Unmarshal(from, &p); err != nil {
		return
	}

	b, err := json.Marshal(p)
	if err != nil {
		return to, fmt.Errorf("could not marshal contect object struct to json: %v", err)
	}
	to = string(b)

	return
}

// Take a querystring from the request and convert it to a valid
// XML string for use in the POST to SFX
func ToCtxObjReq(qs url.Values) (ctxObjReqXml string, err error) {
	ctxObjReq, err := setContextObjectReq(qs)
	if err != nil {
		return ctxObjReqXml, fmt.Errorf("could not create context object for request: %v", err)
	}

	ctxObjReqXml, err = ctxObjReq.toXML()
	if err != nil {
		return ctxObjReqXml, fmt.Errorf("could not convert request context object to XML: %v", err)
	}

	if !isValidXML([]byte(ctxObjReqXml)) {
		return ctxObjReqXml, fmt.Errorf("request context object XML is not valid XML: %v", err)
	}

	return
}

// Construct and run the actual POST request to the SFX server
// Expects an XML string which will be appended to the PostForm params
// Body is blank because that is how SFX expects it
func Post(xmlBody string) (body string, err error) {
	params := url.Values{}
	params.Add("url_ctx_fmt", "info:ofi/fmt:xml:xsd:ctx")
	params.Add("sfx.response_type", "multi_obj_xml")
	params.Add("sfx.show_availability", "1")
	params.Add("sfx.ignore_date_threshold", "1")
	params.Add("sfx.doi_url", "http://dx.doi.org")
	params.Add("url_ctx_val", xmlBody)

	req, err := http.NewRequest("POST", sfxUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return body, fmt.Errorf("could not initialize request to SFX server: %v", err)
	}

	req.PostForm = params
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return body, fmt.Errorf("could not do post to SFX server: %v", err)
	}

	defer resp.Body.Close()
	sfxResp, err := ioutil.ReadAll(resp.Body)

	// Convert to JSON before returning
	body, err = toJson(sfxResp)
	if err != nil {
		return body, fmt.Errorf("could not convert SFX response XML to JSON: %v", err)
	}

	return
}
