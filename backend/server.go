package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"resolve/sfx"

	"github.com/gorilla/mux"
)

//go:embed static
var staticFiles embed.FS

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *mux.Router {
	var staticFS = fs.FS(staticFiles)
	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticContent))

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/healthcheck", Healthcheck).Methods("GET")
	router.HandleFunc("/resolve", ResolveHTML).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	router.HandleFunc("/v0/", ResolveJSON).Methods("GET")

	return router
}

// Take an incoming Querystring, convert to context object XML, send a post to SFX
// and write the response JSON
func ResolveJSON(w http.ResponseWriter, r *http.Request) {
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

func ResolveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content Type", "text/html")

	http.ServeFile(w, r, "./templates/index.html")
}

// Healthcheck returns a successful response, that's it
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}

func handleError(err error, w http.ResponseWriter, message string) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "message": message})
}
