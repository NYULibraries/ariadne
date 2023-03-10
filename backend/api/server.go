package api

import (
	"ariadne/sfx"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Handler for the endpoint used by the frontend
func ResolverHandler(w http.ResponseWriter, r *http.Request) {
	setCORS(&w)

	w.Header().Add("Content-Type", "application/json")

	sfxRequest, err := sfx.NewSFXRequest(r.URL.Query())
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

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", healthCheck)
	router.HandleFunc("/v0/", ResolverHandler)

	return router
}

func setCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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

func makeJSONResponseFromSFXResponse(sfxResponse *sfx.MultipleObjectsResponse) []byte {
	// Remove the Ask a Librarian target -- for details, see:
	// https://nyu-lib.monday.com/boards/765008773/pulses/3548498827
	sfxResponse.RemoveTarget("http://library.nyu.edu/ask/")

	ariadneResponse := Response{
		Errors:  []string{},
		Records: sfxResponse.MultiObjXMLResponseBody,
	}

	responseJSON, err := json.MarshalIndent(ariadneResponse, "", "    ")
	// Very unlikely that this will error out.  At the moment, can't even think
	// of a way to force an error so that can write a test.  Tested this error handling
	// code during development by setting `err = errors.New("error!")` right after marshalling.
	// Result:
	// {
	//     "errors": [
	//         "Could not marshal ariadne response to JSON: error!"
	//     ],
	//     "records": {
	//         "ctx_obj": null
	//     }
	// }
	if err != nil {
		ariadneResponse = Response{
			Errors:  []string{fmt.Sprintf("Could not marshal ariadne response to JSON: %v", err)},
			Records: sfx.MultiObjXMLResponseBody{},
		}

		// Even more unlikely that this marshalling will error out, but if it does
		// just let the chips fall.  The frontend will report the error to receive
		// an intelligible response from the backend API.
		responseJSON, _ = json.MarshalIndent(ariadneResponse, "", "    ")
	}

	return responseJSON
}
