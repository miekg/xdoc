package main

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

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

		p := PathToProject(group, proj)
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
		vars := mux.Vars(r)
		group := vars["group"]
		proj := vars["project"]
		rest := vars["rest"]
		// rest is either a path, or the first element is the project and the project is actually a subgroup

		// First check if group/project exists.
		p := PathToProject(group, proj)
		gl := d.Fetch(p)
		if gl != nil {
			file := path.Join(p, rest)
			buf := d.FetchFile(gl.Project, file)
			if buf == nil {
				http.Error(w, fmt.Sprintf("file %q: %s", file, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
				return
			}
			render(w, r, buf, "index.md")
			return
		}

		// assume first element of rest is the project and proj is actually a subgroup
		subgroup := proj
		el := filepath.SplitList(p) // this works on Linux, but will break on Windows
		proj = el[0]
		rest = RemoveFirstPathElement(p)
		if rest == "" {
			rest = "index.md"
		}
		p = PathToProject(group, subgroup, proj)
		gl = d.Fetch(p)
		if gl != nil {
			file := path.Join(p, rest)
			buf := d.FetchFile(gl.Project, file)
			if buf == nil {
				http.Error(w, fmt.Sprintf("file %q: %s", file, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
				return
			}
			render(w, r, buf, rest)
			return
		}
		// TODO(miek): this error is too generic for both group and subgroup cases
		http.Error(w, fmt.Sprintf("project %q: %s", p, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
	})

	return r
}
