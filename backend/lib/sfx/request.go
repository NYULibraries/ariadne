package sfx

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"text/template"
	"time"
)

// SFX url
const sfxUrl string = "http://sfx.library.nyu.edu/sfxlcl41"

//go:embed templates/sfx-request.xml
var sfxRequestTemplate string

type Timestamp time.Time

// Values needed for templating an SFX request are parsed
type sfxContextObjectRequestBody struct {
	RftValues *OpenURL
	Timestamp string
	Genre     string
}

type OpenURL map[string][]string

// Object representing everything that's needed to request from SFX
type SFXContextObjectRequest struct {
	RequestXML string
}

type SFXRequest interface {
	Request()
	toRequestXML()
}

// Take a querystring from the request and convert it to a valid
// XML string for use in the POST to SFX, return SFXContextObjectRequest object
func Init(qs url.Values) (sfxContextObjectRequest *SFXContextObjectRequest, err error) {
	sfxContextObjectRequest, err = setSFXContextObjectRequest(qs)
	if err != nil {
		return sfxContextObjectRequest, fmt.Errorf("could not create context object for request: %v", err)
	}

	return
}

// Construct and run the actual POST request to the SFX server
// Expects an XML string in a SFXContextObjectRequest obj which will be appended to the PostForm params
// Body is blank because that is how SFX expects it
func (c SFXContextObjectRequest) Request() (body string, err error) {
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
// Store in SFXContextObjectRequest.RequestXML
func (c *SFXContextObjectRequest) toRequestXML(tplVals sfxContextObjectRequestBody) error {
	t := template.New("sfx-request.xml").Funcs(template.FuncMap{"ToLower": strings.ToLower})

	t, err := t.Parse(sfxRequestTemplate)
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

// Setup the SFXContextObjectTpl instance we'll need to run with
// the gotemplates to create the valid XML string param
func setSFXContextObjectRequest(qs url.Values) (sfxContext *SFXContextObjectRequest, err error) {
	rfts, err := parseOpenURL(qs)
	if err != nil {
		return sfxContext, fmt.Errorf("could not parse OpenURL: %v", err)
	}

	genre, err := validGenre((*rfts)["genre"])
	if err != nil {
		return sfxContext, fmt.Errorf("genre is not valid: %v", err)
	}

	// Set up template values, but discard after generating requestXML
	now := time.Now()
	tmpl := sfxContextObjectRequestBody{
		Timestamp: now.Format(time.RFC3339Nano),
		RftValues: rfts,
		Genre:     genre,
	}

	// Init the empty object to populate with toRequestXML
	sfxContext = &SFXContextObjectRequest{}

	if err := sfxContext.toRequestXML(tmpl); err != nil {
		return sfxContext, fmt.Errorf("could not convert request context object to XML: %v", err)
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
func validGenre(genre []string) (string, error) {
	validGenres := genresList()
	if len(genre) > 0 && validGenres[genre[0]] {
		return genre[0], nil
	}
	return "", fmt.Errorf("genre not in list of allowed genres: %v", genre)
}

// Take an openurl and return an OpenURL object of only the rft-prefixed fields
// These are the fields we are going to parse into XML as part of the
// post request params
func parseOpenURL(qs url.Values) (*OpenURL, error) {
	parsed := &OpenURL{}

	for k, v := range qs {
		// Strip the "rft." prefix from the OpenURL
		// and map into valid OpenURL fields
		if strings.HasPrefix(k, "rft.") {
			// E.g. "rft.book" becomes "book"
			newKey := strings.Split(k, ".")[1]
			(*parsed)[newKey] = v
		}
	}

	if reflect.DeepEqual(parsed, &OpenURL{}) {
		return nil, fmt.Errorf("no valid querystring values to parse")
	}

	return parsed, nil
}

// Validate XML, by marshalling and checking for a blank error
func isValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// Convert the response XML from SFX into a JSON string
func toResponseJson(from []byte) (to string, err error) {
	var p SFXContextObjectSet
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
