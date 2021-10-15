package main

import (
	"encoding/json"
	"fmt"
	"getit/lib/sfx"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Run on port 8080
const appPort = "8080"

func main() {
	router := NewRouter()

	log.Println("Listening on port", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, router))
}

// Setup a new mux router with the appropriate routes for this app
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthcheck", Healthcheck).Methods("GET")
	router.HandleFunc("/", Index).Methods("GET")

	return router
}

// Take an incoming Querystring, convert to context object XML, send a post to SFX
// and write the response JSON
// TODO: Throttle requests
func Index(w http.ResponseWriter, r *http.Request) {
	// s := "http://yourserver:3000/resolve?sid=FirstSearch%3AWorldCat&genre=book&title=Fairy+tales&date=1898&aulast=Andersen&aufirst=H&auinitm=C&rfr_id=info%3Asid%2Ffirstsearch.oclc.org%3AWorldCat&rft.genre=book&rft_id=info%3Aoclcnum%2F7675437&rft.aulast=Andersen&rft.aufirst=H&rft.auinitm=C&rft.btitle=Fairy+tales&rft.date=1898&rft.place=Philadelphia&rft.pub=H.+Altemus+Co.&rft.genre=book"
	w.Header().Add("Content-Type", "application/json")

	ctx, err := sfx.ToCtxObjReq(r.URL.Query())
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := sfx.Post(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(w, resp)
}

// Healthcheck returns a successful response, that's it
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}
