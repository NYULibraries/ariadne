package sfx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
)

type openURL map[string][]string

// SFX service URL
const DefaultSFXURL = "http://sfx.library.nyu.edu/sfxlcl41"

var sfxURL = DefaultSFXURL

func Do(request *MultipleObjectsRequest) (*MultipleObjectsResponse, error) {
	return request.do()
}

// Take a querystring from the request and convert it to a valid
// XML string for use in the POST to SFX, return MultipleObjectsRequest object
func NewSFXMultipleObjectsRequest(qs url.Values) (multipleObjectsRequest *MultipleObjectsRequest, err error) {
	multipleObjectsRequest, err = setMultipleObjectsRequest(qs)
	if err != nil {
		return multipleObjectsRequest, fmt.Errorf("could not create a multiple objects request for query string values: %v", err)
	}

	return
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

// Convert the response XML from SFX into a JSON string
func toResponseJSON(from []byte) (to string, err error) {
	var p MultiObjXMLResponseBody
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
