package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/mux"
)

func (d *Doc) setup() http.Handler {
	// s/ -> search related
	// r/ -> render related
	// a/ -> asset related, css etc.

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		io.WriteString(w, htmlForm)
	})
	r.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		io.WriteString(w, htmlForm)
	})
	r.HandleFunc("/s/{query}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := bleve.NewMatchQuery(vars["query"])
		search := bleve.NewSearchRequest(query)
		results, err := d.Index().Search(search)
		if err != nil {
			log.Print(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		io.WriteString(w, d.HTML(results))
	})
	r.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		query := bleve.NewMatchQuery(r.FormValue("q"))
		search := bleve.NewSearchRequest(query)
		results, err := d.Index().Search(search)
		if err != nil {
			log.Print(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		io.WriteString(w, d.HTML(results))
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

		buf := d.FetchFile(gl.Project, "index.md")
		if buf == nil {
			http.Error(w, fmt.Sprintf("file %q: %s", "index.md", http.StatusText(http.StatusNotFound)), http.StatusNotFound)
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
			buf := d.FetchFile(gl.Project, rest)
			if buf == nil {
				http.Error(w, fmt.Sprintf("project %q, file %q: %s", p, rest, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
				return
			}
			render(w, r, buf, "index.md")
			return
		}

		// assume first element of rest is the project and proj is actually a subgroup
		subgroup := proj
		el := strings.Split(p, "/")
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
				http.Error(w, fmt.Sprintf("project %q, file %q: %s", p, file, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
				return
			}
			render(w, r, buf, rest)
			return
		}
		http.Error(w, fmt.Sprintf("project %q: %s", p, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
	})

	r.Path("/a/{asset}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		a := vars["asset"]
		asset, ok := Assets[a]
		if !ok {
			http.Error(w, fmt.Sprintf("asset %q: %s", a, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", asset.contenttype)
		io.WriteString(w, asset.content)
	})

	return r
}
