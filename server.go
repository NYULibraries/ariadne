package main

import (
	"fmt"
	"getit/lib/sfx"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const sfxUrl string = "http://sfx.library.nyu.edu/sfxlcl41"

type xmlStr string

func postToSfx(xmlBody string) (bodyString xmlStr, err error) {
	client := http.Client{}
	params := url.Values{}

	params.Add("url_ctx_fmt", "info:ofi/fmt:xml:xsd:ctx")
	params.Add("sfx.response_type", "multi_obj_xml")
	params.Add("sfx.show_availability", "1")
	params.Add("sfx.ignore_date_threshold", "1")
	params.Add("sfx.doi_url", "http://dx.doi.org")
	params.Add("url_ctx_val", xmlBody)

	req, err := http.NewRequest("POST", sfxUrl, strings.NewReader(params.Encode()))
	req.PostForm = params
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyString = xmlStr(body)
	return
}

func (str xmlStr) toJSON() (json string, err error) {
	return
}

func main() {
	s := "http://yourserver:3000/resolve?sid=FirstSearch%3AWorldCat&genre=book&title=Fairy+tales&date=1898&aulast=Andersen&aufirst=H&auinitm=C&rfr_id=info%3Asid%2Ffirstsearch.oclc.org%3AWorldCat&rft.genre=book&rft_id=info%3Aoclcnum%2F7675437&rft.aulast=Andersen&rft.aufirst=H&rft.auinitm=C&rft.btitle=Fairy+tales&rft.date=1898&rft.place=Philadelphia&rft.pub=H.+Altemus+Co.&rft.genre=book"

	ctx, err := sfx.CreateNewCtx(s)
	if err != nil {
		panic(err)
	}
	result, err := ctx.ToXML()
	if err != nil {
		panic(err)
	}
	if !sfx.IsValidXML([]byte(result)) {
		panic("PANIC!!!")
	}

	bodyString, err := postToSfx(result)
	if err != nil {
		panic(err)
	}
	json, err := bodyString.toJSON()
	fmt.Println(json)

}
