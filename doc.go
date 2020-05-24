package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/xanzy/go-gitlab"
)

// Flavor is the markdown flavor used for parsing the markdown files.
type Flavor int

const (
	None Flavor = iota
	Mmark
)

type Doc struct {
	projects map[string]*GitLab // Basename(URL) -> Project + potential metadata
	i        bleve.Index
	rw       sync.RWMutex // protects Projects and Index
}

// New returns a new and initialized pointer to a Doc. Note the Bleve index is not set.
func New() *Doc {
	return &Doc{projects: make(map[string]*GitLab), rw: sync.RWMutex{}}
}

type GitLab struct {
	*gitlab.Project
	Options
	Files map[string][]byte
}

func (g *GitLab) String() string {
	// strings.Builder ?
	s := fmt.Sprintf("** %s: %d files\n", g.WebURL, len(g.Files))
	for k, v := range g.Files {
		s += "\t" + k + fmt.Sprintf(", %d bytes\n", len(v))
	}
	return s
}

func (d *Doc) String() string {
	d.rw.RLock()
	defer d.rw.RUnlock()
	s := fmt.Sprintf("%d Projects\n", len(d.projects))
	for _, v := range d.projects {
		s += v.String()
	}
	return s
}

// Insert inserts a new project into d.
func (d *Doc) Insert(p *gitlab.Project, opts Options) {
	d.rw.Lock()
	defer d.rw.Unlock()
	urlp := ProjectToPath(p)
	d.projects[urlp] = &GitLab{Project: p, Options: opts, Files: make(map[string][]byte)}
}

// Fetch will return the project belonging to path. Will return nil if not found.
func (d *Doc) Fetch(path string) *GitLab {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.projects[path]
}

func (d *Doc) InsertFile(p *gitlab.Project, pathname string, buf []byte) {
	urlp := ProjectToPath(p)
	gl := d.Fetch(urlp)
	if gl == nil {
		return
	}

	d.rw.Lock()
	defer d.rw.Unlock()
	stripped := RemoveFirstPathElement(pathname)
	full := path.Join(urlp, stripped)
	gl.Files[full] = buf
}

// FetchFile return the file for p associated with pathname. Pathname must be without the doc dir.
func (d *Doc) FetchFile(p *gitlab.Project, pathname string) []byte {
	urlp := ProjectToPath(p)
	gl := d.Fetch(urlp)
	if gl == nil {
		return nil
	}

	d.rw.RLock()
	defer d.rw.RUnlock()
	full := path.Join(urlp, pathname)
	return gl.Files[full]
}

// SetIndex sets the Bleve index in doc.
func (d *Doc) SetIndex(i bleve.Index) {
	d.rw.Lock()
	defer d.rw.Unlock()
	d.i = i
}

// Index returns the Bleve index from doc.
func (d *Doc) Index() bleve.Index {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.i
}

// ProjectToPath converts a gitlab project to a path that can be used in Fetch.
func ProjectToPath(p *gitlab.Project) string {
	url, _ := url.Parse(p.WebURL)
	return url.Path
}

// PathToProject joins the elements and create a project string
func PathToProject(elem ...string) string {
	elem = append([]string{"/"}, elem...)
	return path.Join(elem...)
}

// RemoveFirstPathElement removes the first element from the path p. This is need to remove the Docs dir from
// the files downloaded from GitLab.
func RemoveFirstPathElement(p string) string {
	el := strings.Split(p, "/")
	if len(el) == 0 {
		return ""
	}
	return path.Join(el[1:]...)
}
