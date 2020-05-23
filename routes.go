package main

import (
	"fmt"
	"net/http"
	"path"

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
	r.HandleFunc("/r/{group}/{project}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		group := vars["group"]
		proj := vars["project"]

		p := path.Join(group, proj)
		gl := d.Fetch(p)
		if gl == nil {
			http.Error(w, fmt.Sprintf("project %q: %s", p, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
			return
		}

		file := path.Join(p, "index.md")
		buf := d.FetchFile(gl.Project, file)
		if buf == nil {
			http.Error(w, fmt.Sprintf("file %q: %s", file, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
			return
		}
		render(w, r, buf, "index.md")
	})
	r.Path("/r/{group}/{project}/{rest:.*}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
