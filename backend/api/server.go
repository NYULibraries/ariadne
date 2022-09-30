package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"resolve/sfx"
)

// Healthcheck returns a successful response, that's it
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", Healthcheck)
	router.HandleFunc("/v0/", ResolveJSON)

	return router
}

// Take an incoming Querystring, convert to context object XML, send a post to SFX
// and write the response JSON
func ResolveJSON(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	w.Header().Add("Content-Type", "application/json")

	sfxContext, err := sfx.NewSFXContextObjectRequest(r.URL.Query())
	if err != nil {
		handleError(err, w, "Invalid OpenURL")
		return
	}

	response, err := sfxContext.Request()
	if err != nil {
		handleError(err, w, "Invalid response from SFX")
		return
	}

	fmt.Fprintln(w, response)

}

func handleError(err error, w http.ResponseWriter, message string) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "message": message})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
