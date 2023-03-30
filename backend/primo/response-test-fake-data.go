package primo

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

var fakePrimoISBNSearchAPIResponse = APIResponse{
	Docs: []Doc{
		{
			Delivery: Delivery{
				Link: []Link{
					{
						HyperlinkText: "1",
						LinkURL:       "https://fake.com/1/",
						LinkType:      "http://purl.org/pnx/linkType/linktoprice",
					},
					{
						HyperlinkText: "2",
						LinkURL:       "https://fake.com/2/",
						LinkType:      "http://purl.org/pnx/linkType/linktorsrc",
					},
					{
						HyperlinkText: "3",
						LinkURL:       "https://fake.com/3/",
						LinkType:      "http://purl.org/pnx/linkType/linktoprice",
					},
					{
						HyperlinkText: "4",
						LinkURL:       "https://fake.com/4/",
						LinkType:      "http://purl.org/pnx/linkType/linktorsrc",
					},
				},
			},
			PNX: PNX{
				Facets: Facets{
					FRBRType: []string{
						"5",
					},
					FRBRGroupID: []string{
						"1234567890",
					},
				},
				Search: Search{
					ISBN: []string{
						"1111111111111",
						"2222222222222",
						"3333333333333",
						"4444444444444",
					},
				},
			},
		},
	},
}

var fakePrimoISBNSearchHTTPResponseBody = `{
  "docs": [
    {
      "delivery": {
        "link": [
          {
            "hyperlinkText": "1",
            "linkURL": "https://fake.com/1/",
            "linkType": "http://purl.org/pnx/linkType/linktoprice"
          },
          {
            "hyperlinkText": "2",
            "linkURL": "https://fake.com/2/",
            "linkType": "http://purl.org/pnx/linkType/linktorsrc"
          },
          {
            "hyperlinkText": "3",
            "linkURL": "https://fake.com/3/",
            "linkType": "http://purl.org/pnx/linkType/linktoprice"
          },
          {
            "hyperlinkText": "4",
            "linkURL": "https://fake.com/4/",
            "linkType": "http://purl.org/pnx/linkType/linktorsrc"
          }
        ]
      },
      "pnx": {
        "facets": {
          "frbrtype": [
            "5"
          ],
          "frbrgroupid": [
            "1234567890"
          ]
        },
        "search": {
          "isbn": [
            "1111111111111",
            "2222222222222",
            "3333333333333",
            "4444444444444"
          ]
        }
      }
    }
  ]
}
`

var fakePrimoISBNSearchHTTPResponse = &http.Response{
	Body: ioutil.NopCloser(bytes.NewBufferString(fakePrimoISBNSearchHTTPResponseBody)),
}

var fakePrimoISBNSearchHTTPResponseInvalid = &http.Response{
	Body: ioutil.NopCloser(bytes.NewBufferString("<invalid></invalid>")),
}

var fakeDumpedPrimoISBNSearchHTTPResponse = `HTTP/0.0 000 status code 0

{
  "docs": [
    {
      "delivery": {
        "link": [
          {
            "hyperlinkText": "1",
            "linkURL": "https://fake.com/1/",
            "linkType": "http://purl.org/pnx/linkType/linktoprice"
          },
          {
            "hyperlinkText": "2",
            "linkURL": "https://fake.com/2/",
            "linkType": "http://purl.org/pnx/linkType/linktorsrc"
          },
          {
            "hyperlinkText": "3",
            "linkURL": "https://fake.com/3/",
            "linkType": "http://purl.org/pnx/linkType/linktoprice"
          },
          {
            "hyperlinkText": "4",
            "linkURL": "https://fake.com/4/",
            "linkType": "http://purl.org/pnx/linkType/linktorsrc"
          }
        ]
      },
      "pnx": {
        "facets": {
          "frbrtype": [
            "5"
          ],
          "frbrgroupid": [
            "1234567890"
          ]
        },
        "search": {
          "isbn": [
            "1111111111111",
            "2222222222222",
            "3333333333333",
            "4444444444444"
          ]
        }
      }
    }
  ]
}`

var fakeDumpedPrimoISBNSearchHTTPResponseInvalid = `HTTP/0.0 000 status code 0

<invalid></invalid>`
