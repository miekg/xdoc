package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (d *Doc) setup() {
	r := mux.NewRouter()

	//
	r.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		// bleve search thing
	})
	r.HandleFunc("/s/{query}", func(w http.ResponseWriter, r *http.Request) {
		// vars := mux.Vars(r)
		// bleve search thing
	})
	r.HandleFunc("/r/{group}/{subgroup}/{project}", func(w http.ResponseWriter, r *http.Request) {
		// vars := mux.Vars(r)
		// render index.md
	})
	r.HandleFunc("/r/{group}/{subgroup}/{project}/{file}", func(w http.ResponseWriter, r *http.Request) {
		// vars := mux.Vars(r)
		// render file from the docs dir
	})
	r.HandleFunc("/r/{group}/{project}", func(w http.ResponseWriter, r *http.Request) {
		// vars := mux.Vars(r)
		// render index.md
	})
	r.HandleFunc("/r/{group}/{project}/{file}", func(w http.ResponseWriter, r *http.Request) {
		// vars := mux.Vars(r)
		// render file from the docs dir
	})
}
