package api

import (
	"ariadne/log"
	"ariadne/primo"
	"ariadne/sfx"
	"ariadne/testutils"
	"ariadne/util"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"
)

const elidedHost = "\"Host\":\"[ELIDED]\""
const elidedDatestamp = "\"Date\":\"[ELIDED]\""
const elidedTimestamp = "\"time\":\"[ELIDED]\""

const loggingTestCaseKey = "contrived-frbr-group-test-case"

var logOutputStringDatestampRegexp = regexp.MustCompile("Date:.*GMT")
var logOutputStringHostRegexp = regexp.MustCompile("Host: 127.0.0.1:\\d*")
var logOutputStringTimestampRegexp = regexp.MustCompile("\"time\":\"[^\"]*\"")

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
	var currentTestCase testutils.TestCase

	// Set up Primo service fake
	fakePrimoServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				t.Fatal(err)
			}

			// There potentially two kinds of requests:
			//     - ISBN search request: this is the initial request that is
			//       always made if Primo is being used at all
			//     - FRBR member search request: if the response to the initial
			//       ISBN search request returns docs that indicate an active FRBR
			//       group, more requests are made with an extra query param added
			//       to the query string of the ISBN search request.
			var primoFakeResponse string
			if params.Get(primo.FRBRMemberSearchQueryParamName) == "" {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseISBNSearch(currentTestCase)
			} else {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseFRBRMemberSearch(currentTestCase)
			}

			if err != nil {
				t.Fatal(err)
			}

			_, err = fmt.Fprint(w, primoFakeResponse)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer fakePrimoServer.Close()

	primo.SetPrimoURL(fakePrimoServer.URL)

	// Set up SFX service fake
	fakeSFXServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sfxFakeResponse, err := testutils.GetSFXFakeResponse(currentTestCase)
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

	// Disable logging or else we'll have a ton on noise in the test results
	// output.
	log.SetLevel(log.LevelDisabled)

	for _, testCase := range testutils.TestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			currentTestCase = testCase

			request, err := http.NewRequest(
				"GET",
				"/v0/?"+testCase.QueryString,
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

			goldenValue, err := testutils.GetGoldenValue(testCase)
			if err != nil {
				t.Fatalf("Error retrieving golden value for test case \"%s\": %s",
					testCase.Name, err)
			}

			actualValue := string(body)
			if actualValue != goldenValue {
				err := writeActualToTmp(testCase, actualValue)
				if err != nil {
					t.Fatalf("Error writing actual temp file for test case \"%s\": %s",
						testCase.Name, err)
				}

				goldenFile := testutils.GoldenFile(testCase)
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

func TestLogging(t *testing.T) {
	var loggingTestCase testutils.TestCase

	// Set up SFX service fake
	fakeSFXServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sfxFakeResponse, err := testutils.GetSFXFakeResponse(loggingTestCase)
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

	// Set up Primo service fake
	fakePrimoServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				t.Fatal(err)
			}

			// There potentially two kinds of requests:
			//     - ISBN search request: this is the initial request that is
			//       always made if Primo is being used at all
			//     - FRBR member search request: if the response to the initial
			//       ISBN search request returns docs that indicate an active FRBR
			//       group, more requests are made with an extra query param added
			//       to the query string of the ISBN search request.
			var primoFakeResponse string
			if params.Get(primo.FRBRMemberSearchQueryParamName) == "" {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseISBNSearch(loggingTestCase)
			} else {
				primoFakeResponse, err = testutils.GetPrimoFakeResponseFRBRMemberSearch(loggingTestCase)
			}

			if err != nil {
				t.Fatal(err)
			}

			_, err = fmt.Fprint(w, primoFakeResponse)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer fakePrimoServer.Close()

	primo.SetPrimoURL(fakePrimoServer.URL)

	router := NewRouter()

	for _, testCase := range testutils.TestCases {
		if testCase.Key == loggingTestCaseKey {
			loggingTestCase = testCase
			break
		}
	}

	// Set logging level and redirect output to a buffer.
	log.SetLevel(log.LevelInfo)
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	t.Run(loggingTestCase.Name, func(t *testing.T) {
		request, err := http.NewRequest(
			"GET",
			"/v0/?"+loggingTestCase.QueryString,
			nil,
		)
		if err != nil {
			t.Fatalf("Error creating new HTTP request: %s", err)
		}

		// We're not recording any responses at the moment, but we need to pass
		// in a http.ResponseWriter anyway to router.ServeHTTP, so why not make
		// it a recorder.
		responseRecorder := httptest.NewRecorder()
		router.ServeHTTP(responseRecorder, request)
		actualLogOutputString := normalizeLogOutputString(logOutput.String())
		expectedLogOutputString := normalizeLogOutputString("{\"time\":\"2023-04-04T13:54:42.80309-04:00\",\"level\":\"INFO\",\"msg\":\"SFX request\",\"sfxRequest.DumpedHTTPRequest\":\"GET /?%3Ftitle=Contrived+FRBR+Group+Test+Case&date=1999&isbn=1111111111111&sfx.doi_url=http%3A%2F%2Fdx.doi.org&sfx.response_type=multi_obj_xml&url_ctx_fmt=info%3Aofi%2Ffmt%3Axml%3Axsd%3Actx HTTP/1.1\\r\\nHost: 127.0.0.1:63502\\r\\n\\r\\n\"}\n{\"time\":\"2023-04-04T13:54:43.644324-04:00\",\"level\":\"INFO\",\"msg\":\"SFX response\",\"sfxResponse.DumpedHTTPResponse\":\"HTTP/1.1 200 OK\\r\\nTransfer-Encoding: chunked\\r\\nContent-Type: text/plain; charset=utf-8\\r\\nDate: Tue, 04 Apr 2023 17:54:43 GMT\\r\\n\\r\\n1289\\r\\nHTTP/1.1 200 OK\\r\\nTransfer-Encoding: chunked\\r\\nContent-Type: application/xml; charset=ISO-8859-1\\r\\nDate: Thu, 30 Mar 2023 22:18:01 GMT\\r\\nServer: Apache\\r\\n\\r\\n11e4\\r\\n<?xml version=\\\"1.0\\\" encoding=\\\"utf-8\\\"?>\\n\\n<ctx_obj_set>\\n <ctx_obj identifier=\\\"\\\">\\n  <ctx_obj_attributes>&lt;perldata&gt;\\n &lt;hash&gt;\\n  &lt;item key=\\\"fetchid\\\"&gt;1111111111111&lt;/item&gt;\\n  &lt;item key=\\\"_stash\\\"&gt;\\n   &lt;hash&gt;\\n   &lt;/hash&gt;\\n  &lt;/item&gt;\\n  &lt;item key=\\\"req.session_id\\\"&gt;sBBC5CFFC-CF48-11ED-AF63-75004131B499&lt;/item&gt;\\n  &lt;item key=\\\"rft.btitle\\\"&gt;5-Minute Clinical Suite: Version 9.0&lt;/item&gt;\\n  &lt;item key=\\\"sfx.doi_url\\\"&gt;http://dx.doi.org&lt;/item&gt;\\n  &lt;item key=\\\"url_ctx_fmt\\\"&gt;info:ofi/fmt:xml:xsd:ctx&lt;/item&gt;\\n  &lt;item key=\\\"rft.isbn_10\\\"&gt;&lt;/item&gt;\\n  &lt;item key=\\\"sfx.response_type\\\"&gt;multi_obj_xml&lt;/item&gt;\\n  &lt;item key=\\\"rft.year\\\"&gt;1999&lt;/item&gt;\\n  &lt;item key=\\\"rft.date\\\"&gt;1999&lt;/item&gt;\\n  &lt;item key=\\\"rft.isbn\\\"&gt;1111111111111&lt;/item&gt;\\n  &lt;item key=\\\"rft.object_type\\\"&gt;BOOK&lt;/item&gt;\\n  &lt;item key=\\\"sfx.sourcename\\\"&gt;DEFAULT&lt;/item&gt;\\n  &lt;item key=\\\"rft.language\\\"&gt;eng&lt;/item&gt;\\n  &lt;item key=\\\"sfx.request_id\\\"&gt;25793894&lt;/item&gt;\\n  &lt;item key=\\\"sfx.ignore_char_set\\\"&gt;1&lt;/item&gt;\\n  &lt;item key=\\\"rft.genre\\\"&gt;book&lt;/item&gt;\\n  &lt;item key=\\\"sfx.sid\\\"&gt;DEFAULT&lt;/item&gt;\\n  &lt;item key=\\\"rft.pub\\\"&gt;Lippincott Williams &amp;amp; Wilkins&lt;/item&gt;\\n  &lt;item key=\\\"rft.object_id\\\"&gt;4100000012052805&lt;/item&gt;\\n  &lt;item key=\\\"rft.title\\\"&gt;5-Minute Clinical Suite: Version 9.0&lt;/item&gt;\\n  &lt;item key=\\\"@rfe_id\\\"&gt;\\n   &lt;array&gt;\\n   &lt;/array&gt;\\n  &lt;/item&gt;\\n  &lt;item key=\\\"@sfx.searched_by_identifier\\\"&gt;\\n   &lt;array&gt;\\n    &lt;item key=\\\"0\\\"&gt;\\n     &lt;hash&gt;\\n      &lt;item key=\\\"VALUE\\\"&gt;\\n       &lt;array&gt;\\n        &lt;item key=\\\"0\\\"&gt;1111111111111&lt;/item&gt;\\n       &lt;/array&gt;\\n      &lt;/item&gt;\\n      &lt;item key=\\\"SUBTYPE\\\"&gt;&lt;/item&gt;\\n      &lt;item key=\\\"TYPE\\\"&gt;ISBN&lt;/item&gt;\\n     &lt;/hash&gt;\\n    &lt;/item&gt;\\n   &lt;/array&gt;\\n  &lt;/item&gt;\\n  &lt;item key=\\\"existing_ts_ids\\\"&gt;\\n   &lt;array&gt;\\n    &lt;item key=\\\"0\\\"&gt;20430000000000002&lt;/item&gt;\\n    &lt;item key=\\\"1\\\"&gt;111027614344001&lt;/item&gt;\\n   &lt;/array&gt;\\n  &lt;/item&gt;\\n  &lt;item key=\\\"rft.isbn_13\\\"&gt;1111111111111&lt;/item&gt;\\n  &lt;item key=\\\"@rft_id\\\"&gt;\\n   &lt;array&gt;\\n   &lt;/array&gt;\\n  &lt;/item&gt;\\n  \\n &lt;/hash&gt;\\n&lt;/perldata&gt;\\n</ctx_obj_attributes>\\n  <ctx_obj_targets>\\n   <target>\\n    <target_name>DOCDEL_ILLIAD</target_name>\\n    <target_public_name>Request via Interlibrary Loan</target_public_name>\\n    <object_portfolio_id></object_portfolio_id>\\n    <target_id>111027614344000</target_id>\\n    <interface_id>111027614344000</interface_id>\\n    <interface_name>DOCDEL_ILLIAD</interface_name>\\n    <target_service_id>111027614344001</target_service_id>\\n    <service_type>getDocumentDelivery</service_type>\\n    <parser>ILLiad::DDL</parser>\\n    <parse_param>url=https://ill.library.nyu.edu/illiad/illiad.dll/OpenURL &amp; id_type=</parse_param>\\n    <proxy>yes</proxy>\\n    <crossref>yes</crossref>\\n    <note></note>\\n    <authentication></authentication>\\n    <char_set>utf8</char_set>\\n    <displayer></displayer>\\n    <target_url>http://proxy.library.nyu.edu/login?url=https://ill.library.nyu.edu/illiad/illiad.dll/OpenURL?title=5-Minute%20Clinical%20Suite%3A%20Version%209.0&amp;isbn=1111111111111&amp;genre=book&amp;sid=DEFAULT%20(Via%20SFX)&amp;date=1999&amp;year=1999</target_url>\\n    <is_related>no</is_related>\\n    <coverage>\\n     <coverage_text>\\n      <threshold_text></threshold_text>\\n      <embargo_text></embargo_text>\\n     </coverage_text>\\n     <embargo></embargo>\\n    </coverage>\\n   </target>\\n   <target>\\n    <target_name>ASK_A_LIBRARIAN_LCL</target_name>\\n    <target_public_name>Ask a Librarian</target_public_name>\\n    <object_portfolio_id></object_portfolio_id>\\n    <target_id>20430000000000002</target_id>\\n    <interface_id>20430000000000002</interface_id>\\n    <interface_name>ASK_A_LIBRARIAN</interface_name>\\n    <target_service_id>20430000000000002</target_service_id>\\n    <service_type>getWebService</service_type>\\n    <parser>Generic</parser>\\n    <parse_param>IF () \\\"http://library.nyu.edu/ask/\\\"</parse_param>\\n    <proxy>no</proxy>\\n    <crossref>no</crossref>\\n    <note></note>\\n    <authentication></authentication>\\n    <char_set>iso-8859-1</char_set>\\n    <displayer></displayer>\\n    <target_url>http://library.nyu.edu/ask/</target_url>\\n    <is_related>no</is_related>\\n    <coverage>\\n     <coverage_text>\\n      <threshold_text></threshold_text>\\n      <embargo_text></embargo_text>\\n     </coverage_text>\\n     <embargo></embargo>\\n    </coverage>\\n   </target>\\n  </ctx_obj_targets>\\n </ctx_obj>\\n</ctx_obj_set>\\r\\n0\\r\\n\\r\\n\\n\\r\\n0\\r\\n\\r\\n\"}\n{\"time\":\"2023-04-04T13:54:44.404936-04:00\",\"level\":\"INFO\",\"msg\":\"Primo request\",\"primoRequest.DumpedISBNSearchHTTPRequest\":\"GET /?inst=NYU&limit=50&offset=0&q=isbn%2Cexact%2C1111111111111&scope=all&vid=NYU HTTP/1.1\\r\\nHost: 127.0.0.1:63503\\r\\n\\r\\n\"}\n{\"time\":\"2023-04-04T13:54:47.742169-04:00\",\"level\":\"INFO\",\"msg\":\"Primo HTTP FRBR member requests\",\"primoResponse.DumpedFRBRMemberHTTPRequests\":[\"GET /?inst=NYU&limit=50&multiFacets=facet_frbrgroupid%2Cinclude%2C1234567890&offset=0&q=isbn%2Cexact%2C1111111111111&scope=all&vid=NYU HTTP/1.1\\r\\nHost: 127.0.0.1:63503\\r\\n\\r\\n\"]}\n{\"time\":\"2023-04-04T13:54:47.742232-04:00\",\"level\":\"INFO\",\"msg\":\"Primo HTTP responses (initial ISBN search and FRBR member searches)\",\"primoResponse.DumpedHTTPResponses\":[\"HTTP/1.1 200 OK\\r\\nTransfer-Encoding: chunked\\r\\nContent-Type: text/plain; charset=utf-8\\r\\nDate: Tue, 04 Apr 2023 17:54:45 GMT\\r\\n\\r\\ndcc\\r\\n{\\n    \\\"docs\\\": [\\n        {\\n            \\\"delivery\\\": {\\n                \\\"link\\\": [\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C THIS IS AN ACTIVE FRBR GROUP] ISBN search results doc 1, link 1\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/1/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C THIS IS AN ACTIVE FRBR GROUP] ISBN search results doc 1, link 2\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/2/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C THIS IS AN ACTIVE FRBR GROUP] ISBN search results doc 1, link 3\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/3/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C THIS IS AN ACTIVE FRBR GROUP] ISBN search results doc 1, link 4\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/4/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    }\\n                ]\\n            },\\n            \\\"pnx\\\": {\\n                \\\"facets\\\": {\\n                    \\\"frbrtype\\\": [\\n                        \\\"5\\\"\\n                    ],\\n                    \\\"frbrgroupid\\\": [\\n                        \\\"1234567890\\\"\\n                    ]\\n                }\\n            }\\n        },\\n        {\\n            \\\"delivery\\\": {\\n                \\\"link\\\": [\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"ISBN search results doc 2, link 4\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/4/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT THE RIGHT LINK TYPE] ISBN search results doc 2, link 3\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/3/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"ISBN search results doc 2, link 2\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/2/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"ISBN search results doc 2, link 2\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/2/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT THE RIGHT LINK TYPE] ISBN search results doc 2, link 1\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/1/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    }\\n                ]\\n            },\\n            \\\"pnx\\\": {\\n                \\\"facets\\\": {\\n                    \\\"frbrtype\\\": [\\n                        \\\"6\\\"\\n                    ],\\n                    \\\"frbrgroupid\\\": [\\n                        \\\"1234567890\\\"\\n                    ]\\n                }\\n            }\\n        }\\n    ]\\n}\\n\\r\\n0\\r\\n\\r\\n\",\"HTTP/1.1 200 OK\\r\\nTransfer-Encoding: chunked\\r\\nContent-Type: text/plain; charset=utf-8\\r\\nDate: Tue, 04 Apr 2023 17:54:47 GMT\\r\\n\\r\\nf48\\r\\n{\\n    \\\"docs\\\": [\\n        {\\n            \\\"delivery\\\": {\\n                \\\"link\\\": [\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT THE RIGHT LINK TYPE] FRBR member search results doc 1, link 4\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/4/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"ISBN search results doc 2, link 4\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-isbn-search.com/4/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"FRBR member search results doc 1, link 3\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/3/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT THE RIGHT LINK TYPE] FRBR member search results doc 1, link 2\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/2/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"FRBR member search results doc 1, link 1\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/1/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"FRBR member search results doc 1, link 1\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/1/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    }\\n                ]\\n            },\\n            \\\"pnx\\\": {\\n                \\\"search\\\": {\\n                    \\\"isbn\\\": [\\n                        \\\"1111111111111\\\",\\n                        \\\"2222222222222\\\",\\n                        \\\"3333333333333\\\",\\n                        \\\"4444444444444\\\"\\n                    ]\\n                }\\n            }\\n        },\\n        {\\n            \\\"delivery\\\": {\\n                \\\"link\\\": [\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT AN ISBN MATCH] FRBR member search results doc 2, link 1\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/1/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT AN ISBN MATCH] FRBR member search results doc 2, link 2\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/2/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT AN ISBN MATCH] FRBR member search results doc 2, link 3\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/3/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktorsrc\\\"\\n                    },\\n                    {\\n                        \\\"hyperlinkText\\\": \\\"[SHOULD NEVER SEE THIS B/C NOT AN ISBN MATCH] FRBR member search results doc 2, link 4\\\",\\n                        \\\"linkURL\\\": \\\"https://fake-frbr-member-search.com/4/\\\",\\n                        \\\"linkType\\\": \\\"http://purl.org/pnx/linkType/linktoprice\\\"\\n                    }\\n                ]\\n            },\\n            \\\"pnx\\\": {\\n                \\\"search\\\": {\\n                    \\\"isbn\\\": [\\n                        \\\"2222222222222\\\",\\n                        \\\"3333333333333\\\",\\n                        \\\"4444444444444\\\"\\n                    ]\\n                }\\n            }\\n        }\\n    ]\\n}\\n\\r\\n0\\r\\n\\r\\n\"]}\n")
		if actualLogOutputString != expectedLogOutputString {
			t.Errorf("Log output is not correct:\n\nexpected:\n\n%s\n\nactual:\n\n%s",
				expectedLogOutputString, actualLogOutputString)
		}
	})
}

func normalizeLogOutputString(logOutputString string) string {
	result := logOutputStringDatestampRegexp.ReplaceAllString(logOutputString, elidedDatestamp)
	result = logOutputStringHostRegexp.ReplaceAllString(result, elidedHost)
	result = logOutputStringTimestampRegexp.ReplaceAllString(result, elidedTimestamp)

	return result
}

func tmpFile(testCase testutils.TestCase) string {
	return "testdata/server/tmp/actual/" + testCase.Key + ".json"
}

func updateGoldenFile(testCase testutils.TestCase, bytes []byte) error {
	return os.WriteFile(testutils.GoldenFile(testCase), bytes, 0644)
}

func writeActualToTmp(testCase testutils.TestCase, actual string) error {
	return os.WriteFile(tmpFile(testCase), []byte(actual), 0644)
}
