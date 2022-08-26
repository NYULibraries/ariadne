package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"resolve/lib/sfx"

	"github.com/gorilla/mux"
)

//go:embed static
var staticFiles embed.FS

// The function template.Must is a convenience wrapper that panics when passed a non-nil error value,
// and otherwise returns the *Template unaltered
var templates = template.Must(template.ParseFiles("templates/title.html"))

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := os.ReadFile("templates/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

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

func ResolveHTML(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "./templates/index.html")
	p, err := loadPage("title")
	if err != nil {
		handleHtmlError(err, w, "Page not found")
		return
	}
	renderHTML(w, "title", p)
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

func renderHTML(w http.ResponseWriter, tmpl string, p *Page) {

	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleHtmlError(err error, w http.ResponseWriter, message string) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "message": message})
}
