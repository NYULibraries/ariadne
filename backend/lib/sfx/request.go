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

// SFX url
const sfxUrl string = "http://sfx.library.nyu.edu/sfxlcl41"

type Timestamp time.Time

// Values needed for templating an SFX request are parsed
type ctxObjTpl struct {
	RftValues map[string][]string
	Timestamp string
	Genre     string
}

// Object representing everything that's needed to request from SFX
type CtxObjReq struct {
	RequestXML string
}

// Mainly for documenting
type SFXRequest interface {
	Request()
	toRequestXML()
}

// Take a querystring from the request and convert it to a valid
// XML string for use in the POST to SFX, return CtxObjReq object
func Init(qs url.Values) (ctxObjReq *CtxObjReq, err error) {
	ctxObjReq, err = setCtxObjReq(qs)
	if err != nil {
		return ctxObjReq, fmt.Errorf("could not create context object for request: %v", err)
	}

	return
}

// Construct and run the actual POST request to the SFX server
// Expects an XML string in a CtxObjReq obj which will be appended to the PostForm params
// Body is blank because that is how SFX expects it
func (c CtxObjReq) Request() (body string, err error) {
	params := url.Values{}
	params.Add("url_ctx_fmt", "info:ofi/fmt:xml:xsd:ctx")
	params.Add("sfx.response_type", "multi_obj_xml")
	params.Add("sfx.show_availability", "1")
	params.Add("sfx.ignore_date_threshold", "1")
	params.Add("sfx.doi_url", "http://dx.doi.org")
	params.Add("url_ctx_val", c.RequestXML)

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
	if err != nil {
		return body, fmt.Errorf("could not read response from SFX server: %v", err)
	}

	// Convert to JSON before returning
	body, err = toResponseJson(sfxResp)

	if err != nil {
		return body, fmt.Errorf("could not convert SFX response XML to JSON: %v", err)
	}

	return
}

// Convert a context object request to an XML string
// via gotemplates, in order to set it up as a post param to SFX
// Store in CtxObjReq.RequestXML
func (c *CtxObjReq) toRequestXML(tplVals ctxObjTpl) error {
	t := template.New("index.goxml")

	t, err := t.ParseFiles("templates/index.goxml")
	if err != nil {
		return fmt.Errorf("could not load template parse file: %v", err)
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, tplVals); err != nil {
		return fmt.Errorf("could not execute go template from context object request: %v", err)
	}

	if !isValidXML(tpl.Bytes()) {
		return fmt.Errorf("request context object XML is not valid XML: %v", err)
	}

	// Set requestXML to this converted XML string
	c.RequestXML = tpl.String()
	return nil
}

// Setup the CtxObjTpl instance we'll need to run with
// the gotemplates to create the valid XML string param
func setCtxObjReq(qs url.Values) (ctx *CtxObjReq, err error) {
	rfts, err := parseOpenURL(qs)
	if err != nil {
		return ctx, fmt.Errorf("could not parse OpenURL: %v", err)
	}

	isValidGenre, err := isValidGenre(rfts["genre"])
	if err != nil {
		return ctx, fmt.Errorf("genre is not valid: %v", err)
	}

	// Set up template values, but discard after generating requestXML
	now := time.Now()
	tmpl := ctxObjTpl{
		Timestamp: now.Format(time.RFC3339Nano),
		RftValues: rfts,
		Genre:     isValidGenre,
	}

	// Init the empty object to populate with toRequestXML
	ctx = &CtxObjReq{}

	if err := ctx.toRequestXML(tmpl); err != nil {
		return ctx, fmt.Errorf("could not convert request context object to XML: %v", err)
	}

	return
}

// A list of the valid genres as defined by the OpenURL spec
func genresList() (genresList map[string]bool) {
	genresList = map[string]bool{
		"journal":    true,
		"book":       true,
		"conference": true,
		"article":    true,
		"preprint":   true,
		"proceeding": true,
		"bookitem":   true,
	}

	return
}

// Only return a valid genre that has been allowed by the OpenURL spec
func isValidGenre(genre []string) (string, error) {
	validGenres := genresList()
	if len(genre) > 0 && validGenres[genre[0]] {
		return genre[0], nil
	}
	return "", fmt.Errorf("genre not in list of allowed genres: %v", genre)
}

// Take an openurl and return a map of only the rft-prefixed fields
// These are the fields we are going to parse into XML as part of the
// post request params
func parseOpenURL(qs url.Values) (map[string][]string, error) {
	parsed := make(map[string][]string)

	for key, val := range qs {
		if strings.HasPrefix(key, "rft.") {
			// TODO: Dedupe on values
			newKey := strings.Split(key, ".")
			parsed[newKey[1]] = val
		}
	}
	if len(parsed) == 0 {
		return nil, fmt.Errorf("no valid OpenURL querystring params")
	}

	return parsed, nil
}

// Validate XML, by marshalling and checking for a blank error
func isValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// Convert the response XML from SFX into a JSON string
func toResponseJson(from []byte) (to string, err error) {
	var p CtxObjSet
	if err = xml.Unmarshal(from, &p); err != nil {
		return
	}

	b, err := json.Marshal(p)
	if err != nil {
		return to, fmt.Errorf("could not marshal context object struct to json: %v", err)
	}
	to = string(b)

	return
}
