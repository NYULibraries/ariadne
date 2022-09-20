package api

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"resolve/sfx"
	"resolve/util"
	"testing"
)

type TestCase struct {
	// Identifier used for fixture and golden file basename.
	key string
	// Human-readable name/description of test case
	name string
	// OpenURL querystring
	queryString string
}

var updateGoldenFiles = flag.Bool("update-golden-files", false, "update the golden files")

// --update-sfx-fake-responses flag?
// Ideally we also want to have a flag for updating SFX fake response fixture files,
// but it appears that ordering of elements in the SFX response XML and the elements
// in the escaped XML in <perldata> is not stable.
// See comment in monday.com ticket "Add sample integration test for OpenURL resolver":
// https://nyu-lib.monday.com/boards/765008773/pulses/3073776565/posts/1676502313
// Thus the same request submitted multiple times in less than a second
// might end up generating responses that differ only in element ordering.  If this
// in fact is confirmed to be the case, in order for --update-sfx-fake-responses
// to be useful, we would need to write some utility code to normalize the SFX
// responses before writing out the fixture files.

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestResponseJSONRoute(t *testing.T) {
	var currentTestCase TestCase

	testCases := []TestCase{
		{
			key:         "corriere-fiorentino",
			name:        "Corriere Fiorentino",
			queryString: "ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=<NYUMARCIT>3400000000000901</NYUMARCIT><grp_id>582323038</grp_id><oa></oa><url></url>&rft_id=info:oai/&req.language=eng",
		},
		{
			key:         "the-new-yorker",
			name:        "The New Yorker",
			queryString: "url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404<fssessid>0<%2Ffssessid>&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat",
		},
	}

	fakeSFXServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sfxFakeResponse, err := getSFXFakeResponse(currentTestCase)
			if err != nil {
				t.Fatal(err)
			}

			_, err = fmt.Fprint(w, sfxFakeResponse)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer fakeSFXServer.Close()

	sfx.SetSFXURL(fakeSFXServer.URL)

	router := NewRouter()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			currentTestCase = testCase

			request, err := http.NewRequest(
				"GET",
				"/v0/?"+testCase.queryString,
				nil,
			)
			if err != nil {
				t.Fatalf("Error creating new HTTP request: %s", err)
			}

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)

			response := responseRecorder.Result()
			body, _ := io.ReadAll(response.Body)

			if *updateGoldenFiles {
				err = updateGoldenFile(testCase, body)
				if err != nil {
					t.Fatalf("Error updating golden file: %s", err)
				}
			}

			goldenValue, err := getGoldenValue(testCase)
			if err != nil {
				t.Fatalf("Error retrieving golden value for test case \"%s\": %s",
					testCase.name, err)
			}

			actualValue := string(body)
			if actualValue != goldenValue {
				err := writeActualToTmp(testCase, actualValue)
				if err != nil {
					t.Fatalf("Error writing actual temp file for test case \"%s\": %s",
						testCase.name, err)
				}

				goldenFile := goldenFile(testCase)
				actualFile := tmpFile(testCase)
				diff, err := util.Diff(goldenFile, actualFile)
				if err != nil {
					t.Fatalf("Error diff'ing %s vs. %s: %s\n"+
						"Manually diff these files to determine the reasons for test failure.",
						goldenFile, actualFile, err)
				}

				t.Errorf("golden and actual values do not match\noutput of `diff %s %s`:\n%s\n",
					goldenFile, actualFile, diff)
			}
		})
	}
}

func getGoldenValue(testCase TestCase) (string, error) {
	return getTestdataFileContents(goldenFile(testCase))
}

func getSFXFakeResponse(testCase TestCase) (string, error) {
	return getTestdataFileContents(sfxFakeResponseFile(testCase))
}

func getTestdataFileContents(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return filename, err
	}

	return string(bytes), nil
}

func goldenFile(testCase TestCase) string {
	return "testdata/server/golden/" + testCase.key + ".json"
}

func sfxFakeResponseFile(testCase TestCase) string {
	return "testdata/server/fixtures/sfx-fake-responses/" + testCase.key + ".xml"
}

func tmpFile(testCase TestCase) string {
	return "testdata/server/tmp/actual/" + testCase.key + ".json"
}

func updateGoldenFile(testCase TestCase, bytes []byte) error {
	return os.WriteFile(goldenFile(testCase), bytes, 0644)
}

func writeActualToTmp(testCase TestCase, actual string) error {
	return os.WriteFile(tmpFile(testCase), []byte(actual), 0644)
}
