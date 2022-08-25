package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"resolve/lib/sfx"
	"testing"
)

//go:embed all:testdata
var testdataFS embed.FS

// When set, golden files are update with the outputs of the test run.
var updateGolden = flag.Bool("update-golden", false, "update the golden files")
// When set, requests are made to the real SFX service and the fake response fixture files are updated.
var updateSFX = flag.Bool("update-sfx", false, "update the SFX fixtures")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestResponseJSONRoute(t *testing.T) {
	var err error

	sfxFakeResponse, err := getTestdataFileContents("testdata/server/fixtures/sfx-fake-responses/new-yorker.xml")
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON, err := getTestdataFileContents("testdata/server/golden/new-yorker.golden.json")
	if err != nil {
		t.Fatal(err)
	}

	fakeSFXServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, sfxFakeResponse)
		}),
	)
	defer fakeSFXServer.Close()

	sfx.SetSFXURL(fakeSFXServer.URL)

	responseRecorder := httptest.NewRecorder()

	request, err := http.NewRequest(
		"GET",
		"/v0/?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat",
		nil,
	)
	if err != nil {
		t.Fail()
	}

	router := NewRouter()
	router.ServeHTTP(responseRecorder, request)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	if string(body) != expectedJSON {
		t.Errorf(expectedJSON)
	}
}

func getTestdataFileContents(filename string) (string, error) {
	var bytes, err = testdataFS.ReadFile(filename)
	if err != nil {
		return filename, err
	}

	return string(bytes), nil
}
