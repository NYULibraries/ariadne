package sfx

import (
	"bytes"
	"encoding/xml"
	"errors"
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

// Values needed for an SFX request are parsed
type ContextObjectReq struct {
	RftValues map[string][]string
	Timestamp string
	Genre     string
}

// Only return a valid
func validGenre(genre []string) (string, error) {
	if len(genre) > 0 && validGenres[genre[0]] {
		return genre[0], nil
	}
	return "", errors.New("Not a valid genre")
}

// Take an openurl and return a map of only the rft-prefixed fields
func parseOpenURL(s string) (parsed map[string][]string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}

	parsed = make(map[string][]string)
	for key, val := range m {
		if strings.HasPrefix(key, "rft.") {
			// TODO: Dedupe
			newKey := strings.Split(key, ".")
			parsed[newKey[1]] = val
		}
	}

	return
}

func IsValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// Convert a context object request to an XML string
func (c *ContextObjectReq) ToXML() (result string, err error) {
	t := template.New("index.goxml")

	t, err = t.ParseFiles("templates/index.goxml")
	if err != nil {
		return
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, c); err != nil {
		return
	}

	result = tpl.String()
	return
}

// Setup the ContextObject making a request to SFX
func SetContextObjectReq(s string) (ctx *ContextObjectReq, err error) {
	qs, err := parseOpenURL(s)
	if err != nil {
		return
	}

	validGenre, err := validGenre(qs["genre"])
	if err != nil {
		return
	}

	type Timestamp time.Time
	now := time.Now()
	ctx = &ContextObjectReq{
		Timestamp: now.Format(time.RFC3339Nano),
		RftValues: qs,
		Genre:     validGenre,
	}

	return
}
