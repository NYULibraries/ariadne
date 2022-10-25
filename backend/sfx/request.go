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

type openURL map[string][]string

// Object representing everything that's needed to request from SFX
type multipleObjectsRequest struct {
	RequestXML string
}

// Values needed for templating an SFX request are parsed
type multipleObjectsRequestBody struct {
	RftValues *openURL
	Timestamp string
	Genre     string
}

//go:embed templates/sfx-request.xml
var sfxRequestTemplate string

// SFX service URL
const DefaultSFXURL = "http://sfx.library.nyu.edu/sfxlcl41"

var sfxURL = DefaultSFXURL

// Construct and run the actual POST request to the SFX server
// Expects an XML string in a multipleObjectsRequest obj which will be appended to the PostForm params
// Body is blank because that is how SFX expects it
func (c multipleObjectsRequest) Do() (body string, err error) {
	params := url.Values{}
	params.Add("url_ctx_fmt", "info:ofi/fmt:xml:xsd:ctx")
	params.Add("sfx.response_type", "multi_obj_xml")
	// Do we always need these parameters? Umlaut adds them only in certain conditions: https://github.com/team-umlaut/umlaut/blob/master/app/service_adaptors/sfx.rb#L145-L153
	params.Add("sfx.show_availability", "1")
	params.Add("sfx.ignore_date_threshold", "1")
	params.Add("sfx.doi_url", "http://dx.doi.org")
	params.Add("url_ctx_val", c.RequestXML)

	request, err := http.NewRequest("POST", sfxURL, strings.NewReader(params.Encode()))
	if err != nil {
		return body, fmt.Errorf("could not initialize request to SFX server: %v", err)
	}

	request.PostForm = params
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return body, fmt.Errorf("could not do post to SFX server: %v", err)
	}

	defer response.Body.Close()
	sfxResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return body, fmt.Errorf("could not read response from SFX server: %v", err)
	}

	// Convert to JSON before returning
	body, err = toResponseJSON(sfxResponse)

	if err != nil {
		return body, fmt.Errorf("could not convert SFX response XML to JSON: %v", err)
	}

	return
}

// Convert a request to an XML string
// via gotemplates, in order to set it up as a post param to SFX
// Store in multipleObjectsRequest.RequestXML
func (c *multipleObjectsRequest) toRequestXML(tplVals multipleObjectsRequestBody) error {
	t := template.New("sfx-request.xml").Funcs(template.FuncMap{"ToLower": strings.ToLower})

	t, err := t.Parse(sfxRequestTemplate)
	if err != nil {
		return fmt.Errorf("could not load template parse file: %v", err)
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, tplVals); err != nil {
		return fmt.Errorf("could not execute go template from multiple objects request: %v", err)
	}

	if !isValidXML(tpl.Bytes()) {
		return fmt.Errorf("request multiple objects XML is not valid XML: %v", err)
	}

	// Set requestXML to this converted XML string
	c.RequestXML = tpl.String()
	return nil
}

// Take a querystring from the request and convert it to a valid
// XML string for use in the POST to SFX, return multipleObjectsRequest object
func NewSFXMultipleObjectsRequest(qs url.Values) (multipleObjectsRequest *multipleObjectsRequest, err error) {
	multipleObjectsRequest, err = setMultipleObjectsRequest(qs)
	if err != nil {
		return multipleObjectsRequest, fmt.Errorf("could not create a multiple objects request for query string values: %v", err)
	}

	return
}

func SetSFXURL(dependencyInjectedURL string) {
	sfxURL = dependencyInjectedURL
}

// A list of the valid genres as defined by the OpenURL spec
// Is this correct? See genres list on NISO spec page 59: https://groups.niso.org/higherlogic/ws/public/download/14833/z39_88_2004_r2010.pdf
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

// Validate XML, by marshalling and checking for a blank error
func isValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// Take an openurl and return an OpenURL object of only the rft-prefixed fields
// These are the fields we are going to parse into XML as part of the
// post request params
func parseOpenURL(queryStringValues url.Values) (*openURL, error) {
	parsed := &openURL{}

	for k, v := range queryStringValues {
		// Strip the "rft." prefix from the OpenURL
		// and map into valid OpenURL fields
		if strings.HasPrefix(k, "rft.") {
			// E.g. "rft.book" becomes "book"
			newKey := strings.Split(k, ".")[1]
			(*parsed)[newKey] = v
		}
	}

	if reflect.DeepEqual(parsed, &openURL{}) {
		return nil, fmt.Errorf("no valid querystring values to parse")
	}

	return parsed, nil
}

// Setup the SFXContextObjectTpl instance we'll need to run with
// the gotemplates to create the valid XML string param
func setMultipleObjectsRequest(queryStringValues url.Values) (sfxContext *multipleObjectsRequest, err error) {
	rfts, err := parseOpenURL(queryStringValues)
	if err != nil {
		return sfxContext, fmt.Errorf("could not parse OpenURL: %v", err)
	}

	genre, err := validGenre((*rfts)["genre"])
	if err != nil {
		return sfxContext, fmt.Errorf("genre is not valid: %v", err)
	}

	// Set up template values, but discard after generating requestXML
	now := time.Now()
	tmpl := multipleObjectsRequestBody{
		Timestamp: now.Format(time.RFC3339Nano),
		RftValues: rfts,
		Genre:     genre,
	}

	// Init the empty object to populate with toRequestXML
	sfxContext = &multipleObjectsRequest{}

	if err := sfxContext.toRequestXML(tmpl); err != nil {
		return sfxContext, fmt.Errorf("could not convert multiple objects request to XML: %v", err)
	}

	return
}

// Convert the response XML from SFX into a JSON string
func toResponseJSON(from []byte) (to string, err error) {
	var p SFXContextObjectSet
	if err = xml.Unmarshal(from, &p); err != nil {
		return
	}

	b, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return to, fmt.Errorf("could not marshal context object struct to json: %v", err)
	}
	to = string(b)

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
