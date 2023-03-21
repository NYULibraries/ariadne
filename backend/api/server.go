package api

import (
	"ariadne/sfx"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", healthCheck)
	router.HandleFunc("/v0/", ResolverHandler)

	return router
}

// Handler for the endpoint used by the frontend
func ResolverHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	sfxRequest, err := sfx.NewSFXRequest(r.URL.RawQuery)
	if err != nil {
		handleError(err, w, "Invalid OpenURL")
		return
	}

	sfxResponse, err := sfx.Do(sfxRequest)
	if err != nil {
		handleError(err, w, "Invalid response from SFX")
		return
	}

	responseJSON := makeJSONResponseFromSFXResponse(sfxResponse)

	fmt.Fprintln(w, string(responseJSON))
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}

func handleError(err error, w http.ResponseWriter, message string) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "message": message})
}

// healthCheck returns a successful response, that's it
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}

func makeJSONResponseFromSFXResponse(sfxResponse *sfx.SFXResponse) []byte {
	// Remove the Ask a Librarian target -- for details, see:
	// https://nyu-lib.monday.com/boards/765008773/pulses/3548498827
	sfxResponse.RemoveTarget("http://library.nyu.edu/ask/")

	links := []Link{}
	targets := (*(*sfxResponse.XMLResponseBody.ContextObject)[0].SFXContextObjectTargets)[0].Targets
	for _, target := range *targets {
		coverageText := ""
		if target.Coverage != nil {
			firstCoverage := (*(target.Coverage))[0]
			if firstCoverage.CoverageText != nil {
				firstCoverageText := (*(firstCoverage.CoverageText))[0]
				if firstCoverageText.ThresholdText != nil {
					firstThresholdText := (*(firstCoverageText.ThresholdText))[0]
					coverageText = strings.Join(firstThresholdText.CoverageStatement, ". ")
				}
			}
		}
		links = append(links, Link{
			target.TargetPublicName,
			target.TargetUrl,
			coverageText,
		})
	}

	// For now only return one record, but anticipate needing to be able to deliver
	// multiple records later.
	records := []Record{
		{
			CitationSupplemental{},
			links,
		},
	}

	ariadneResponse := Response{
		Errors: []string{},
		// Hardcoding `false` for now.  This is just a placeholder for the value
		// that will be calculated according to https://nyu-lib.monday.com/boards/765008773/pulses/4116986234?userId=27226106
		// acceptance criteria.
		Found:   false,
		Records: records,
	}

	responseJSON, err := json.MarshalIndent(ariadneResponse, "", "    ")
	// Very unlikely that this will error out.  At the moment, can't even think
	// of a way to force an error so that can write a test.  Tested this error handling
	// code during development by setting `err = errors.New("error!")` right after marshalling.
	// Result:
	// {
	// 	   "errors": [
	//		   "Could not marshal ariadne response to JSON: error!"
	//	    ],
	//	    "found": false,
	//	    "records": []
	// }
	if err != nil {
		ariadneResponse = Response{
			Errors:  []string{fmt.Sprintf("Could not marshal ariadne response to JSON: %v", err)},
			Records: []Record{},
		}

		// Even more unlikely that this marshalling will error out, but if it does
		// just let the chips fall.  The frontend will report the error to receive
		// an intelligible response from the backend API.
		responseJSON, _ = json.MarshalIndent(ariadneResponse, "", "    ")
	}

	return responseJSON
}
