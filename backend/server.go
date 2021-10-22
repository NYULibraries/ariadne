package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"resolve/lib/sfx"

	"github.com/gorilla/mux"
)

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthcheck", Healthcheck).Methods("GET")
	router.HandleFunc("/", Resolve).Methods("GET")

	return router
}

// Take an incoming Querystring, convert to context object XML, send a post to SFX
// and write the response JSON
// TODO: Throttle requests? Or throttle in frontend?
func Resolve(w http.ResponseWriter, r *http.Request) {
	// s := "http://yourserver:3000/resolve?sid=FirstSearch%3AWorldCat&genre=book&title=Fairy+tales&date=1898&aulast=Andersen&aufirst=H&auinitm=C&rfr_id=info%3Asid%2Ffirstsearch.oclc.org%3AWorldCat&rft.genre=book&rft_id=info%3Aoclcnum%2F7675437&rft.aulast=Andersen&rft.aufirst=H&rft.auinitm=C&rft.btitle=Fairy+tales&rft.date=1898&rft.place=Philadelphia&rft.pub=H.+Altemus+Co.&rft.genre=book"
	w.Header().Add("Content-Type", "application/json")

	ctx, err := sfx.Init(r.URL.Query())
	if err != nil {
		handleError(err, w, "Invalid OpenURL")
		return
	}

	resp, err := ctx.Request()
	if err != nil {
		handleError(err, w, "Invalid response from SFX")
		return
	}

	fmt.Fprintln(w, resp)
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
