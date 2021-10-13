package ctx_obj

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
