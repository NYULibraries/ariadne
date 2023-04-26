package api

import (
	"ariadne/log"
	"ariadne/primo"
	"ariadne/sfx"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const invalidPrimoRequestErrorMessage = "Invalid Primo request"
const invalidSFXRequestErrorMessage = "Invalid SFX request"

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("/healthcheck", http.HandlerFunc(healthCheck))
	router.Handle("/v0/", recoverWrap(http.HandlerFunc(ResolverHandler)))

	return router
}

// Handler for the endpoint used by the frontend
func ResolverHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	sfxResponse, err := getSFXResponse(r.URL.RawQuery)
	if err != nil {
		handleBadRequestError(err, r, w, err.Error())
		return
	}

	sfxAPIResponseLogEntry := makeNewSFXAPIResponseLogEntry(r.URL.RawQuery, sfxResponse.DumpedHTTPResponse)
	log.Debug(MessageKey, "SFX API Response", AriadneKey, sfxAPIResponseLogEntry)

	var ariadneResponse Response

	if sfxResponse.IsFound() {
		ariadneResponse = makeAriadneResponseFromSFXResponse(sfxResponse)
	} else {
		primoResponse, err := getPrimoResponse(r.URL.RawQuery)
		if err != nil {
			// If we got this far, we already know that Ariadne was able to
			// successfully query SFX request, so we do not want this Primo request
			// error to be fatal, since this we still technically have a valid
			// Ariadne request.  We return the SFX results, which at least will
			// have "helper" links.
			ariadneResponse = makeAriadneResponseFromSFXResponse(sfxResponse)
		} else {
			for i, dumpedFRBRMemberHTTPRequest := range primoResponse.DumpedFRBRMemberHTTPRequests {
				primoAPIFRBRMemberRequestLogEntry :=
					makePrimoAPIFRBRMemberRequestLogEntry(r.URL.RawQuery, dumpedFRBRMemberHTTPRequest)
				log.Info(MessageKey, fmt.Sprintf("Primo API FRBR member request #%d", i+1),
					AriadneKey, primoAPIFRBRMemberRequestLogEntry)
			}

			primoAPIISBNSearchResponseLogEntry :=
				makePrimoAPIISBNSearchResponseLogEntry(r.URL.RawQuery, primoResponse.DumpedHTTPResponses[0])
			log.Debug(MessageKey, "Primo API ISBN Search Response",
				AriadneKey, primoAPIISBNSearchResponseLogEntry)

			for i := 1; i < len(primoResponse.DumpedHTTPResponses); i++ {
				primoAPIFRBRMemberResponseLogEntry :=
					makePrimoAPIFRBRMemberResponseLogEntry(r.URL.RawQuery, primoResponse.DumpedHTTPResponses[1])
				log.Debug(MessageKey, fmt.Sprintf("Primo API FRBR Member Response #%d", i),
					AriadneKey, primoAPIFRBRMemberResponseLogEntry)
			}

			if primoResponse.IsFound() {
				ariadneResponse = makeAriadneResponseFromPrimoResponse(primoResponse)
			} else {
				// Back to SFX again, which at least has some "helper" link
				ariadneResponse = makeAriadneResponseFromSFXResponse(sfxResponse)
			}
		}
	}

	ariadneAPIResponseLogEntry :=
		makeAriadneAPIResponseLogEntry(r.URL.RawQuery, ariadneResponse)
	log.Info(MessageKey, "Ariadne API response", AriadneKey, ariadneAPIResponseLogEntry)

	responseJSON := makeAriadneResponseJSON(ariadneResponse)

	fmt.Fprintln(w, responseJSON)
}

func getPrimoResponse(queryString string) (*primo.PrimoResponse, error) {
	primoRequest, err := primo.NewPrimoRequest(queryString)
	if err != nil {
		return &primo.PrimoResponse{}, errors.New(invalidPrimoRequestErrorMessage)
	}

	primoAPIISBNSearchRequestLogEntry :=
		makePrimoAPIISBNSearchRequestLogEntry(queryString, primoRequest.DumpedISBNSearchHTTPRequest)
	log.Info(MessageKey, "Primo API ISBN Search Request", AriadneKey, primoAPIISBNSearchRequestLogEntry)

	return primo.Do(primoRequest)
}

func getSFXResponse(queryString string) (*sfx.SFXResponse, error) {
	sfxResponse := sfx.SFXResponse{}

	sfxRequest, err := sfx.NewSFXRequest(queryString)
	if err != nil {
		return &sfxResponse, errors.New(invalidSFXRequestErrorMessage)
	}

	sfxAPIRequestLogEntry := makeNewSFXAPIRequestLogEntry(queryString, sfxRequest.DumpedHTTPRequest)
	log.Info(MessageKey, "SFX API Request", AriadneKey, sfxAPIRequestLogEntry)

	return sfx.Do(sfxRequest)
}

func handleBadRequestError(err error, r *http.Request, w http.ResponseWriter, message string) {
	response := Response{
		Errors:  []string{message},
		Found:   false,
		Records: []Record{},
	}
	responseJSON, _ := json.MarshalIndent(response, "", "    ")

	ariadneAPIErrorResponseLogEntry :=
		makeAriadneAPIErrorResponseLogEntry(r.URL.RawQuery, err, http.StatusBadRequest, response)
	log.Error(MessageKey, err.Error(), AriadneKey, ariadneAPIErrorResponseLogEntry)

	http.Error(w, string(responseJSON), http.StatusBadRequest)
}

// healthCheck returns a successful response, that's it
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func makeAriadneResponseFromPrimoResponse(primoResponse *primo.PrimoResponse) Response {
	links := []Link{}
	for _, primoLink := range primoResponse.Links {
		displayName := primoLink.HyperlinkText
		if primoLink.HyperlinkText == "" {
			displayName = "Link to Online Resource"
		}
		links = append(links, Link{
			displayName,
			primoLink.LinkURL,
			"",
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

	return Response{
		Errors:  []string{},
		Found:   primoResponse.IsFound(),
		Records: records,
	}
}

func makeAriadneResponseFromSFXResponse(sfxResponse *sfx.SFXResponse) Response {
	// Remove the Ask a Librarian target -- for details, see:
	// https://nyu-lib.monday.com/boards/765008773/pulses/3548498827
	sfxResponse.RemoveTarget(sfx.AskALibrarianLink)
	emptyTarget := sfxResponse.GetTarget("")
	if emptyTarget != nil {
		log.Warn(MessageKey, "Removing target with empty TargetURL", emptyTarget)
		sfxResponse.RemoveTarget("")
	}

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

	return Response{
		Errors:  []string{},
		Found:   sfxResponse.IsFound(),
		Records: records,
	}
}

func makeAriadneResponseJSON(ariadneResponse Response) string {
	responseJSONBytes, err := json.MarshalIndent(ariadneResponse, "", "    ")
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
		responseJSONBytes, _ = json.MarshalIndent(ariadneResponse, "", "    ")
	}

	return string(responseJSONBytes)
}

// Based on accepted answer for:
//
//	https://stackoverflow.com/questions/28745648/global-recover-handler-for-golang-http-panic
func recoverWrap(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recoverValue := recover()
			if recoverValue != nil {
				var err error
				switch recoverValueType := recoverValue.(type) {
				case string:
					err = errors.New(recoverValueType)
				case error:
					err = recoverValueType
				default:
					err = errors.New("Unknown error")
				}

				response := Response{
					Errors:  []string{err.Error()},
					Found:   false,
					Records: []Record{},
				}
				responseJSON, _ := json.MarshalIndent(response, "", "    ")

				http.Error(w, string(responseJSON), http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}
