package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (d *Doc) setup() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		// bleve search thing
	})
	r.HandleFunc("/s/{query}", func(w http.ResponseWriter, r *http.Request) {
		// vars := mux.Vars(r)
		// bleve search thing
	})
	r.Path("/r/{group}/{project}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r/{group}/{subgroup}/{project}/files...
		// rest may include subgroup, do full search first then
		vars := mux.Vars(r)
		group := vars["group"]
		proj := vars["project"]
		// subgroup may exist
		println(group, proj)
		// r URL to
		// render index.md
	})

	return r
}
