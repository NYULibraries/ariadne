package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"resolve/sfx"
)

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", healthCheck)
	router.HandleFunc("/v0/", multipleRecordsHandler)

	return router
}

func enableCors(w *http.ResponseWriter) {
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

// Take an incoming Querystring, convert to context object XML, send a post to SFX
// and write the response JSON
func multipleRecordsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	w.Header().Add("Content-Type", "application/json")

	sfxRequest, err := sfx.NewMultipleObjectsRequest(r.URL.Query())
	if err != nil {
		handleError(err, w, "Invalid OpenURL")
		return
	}

	sfxResponse, err := sfx.Do(sfxRequest)
	if err != nil {
		handleError(err, w, "Invalid response from SFX")
		return
	}

	ariadneResponse := response{
		Errors:  []string{},
		Records: sfxResponse.MultiObjXMLResponseBody,
	}

	responseJSON, err := json.MarshalIndent(ariadneResponse, "", "    ")
	if err != nil {
		ariadneResponse = response{
			Errors:  []string{fmt.Sprintf("could not marshal ariadne response to JSON: %v", err)},
			Records: sfx.MultiObjXMLResponseBody{},
		}
	}

	fmt.Fprintln(w, string(responseJSON))
}
