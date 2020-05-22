package main

import (
	"net/url"
	"path"
	"sync"

	"github.com/xanzy/go-gitlab"
)

// Flavor is the markdown flavor used for parsing the markdown files.
type Flavor int

const (
	None Flavor = iota
	Mmark
)

type Doc struct {
	Projects map[string]*GitLab // Basename(URL) -> Project + potential metadata
	rw       sync.RWMutex       // protects Projects
}

type GitLab struct {
	*gitlab.Project
	Commit string // not used yet
	Lang   string // not used yet
	Flavor        // not used yet
	Files  map[string][]byte
}

// New creates a new, initialized pointer to a GitLab.
func New() *GitLab {
	return &GitLab{Files: make(map[string][]byte)}
}

// Insert inserts a new project into d.
func (d *Doc) Insert(p *gitlab.Project) {
	d.rw.Lock()
	defer d.rw.Unlock()
	urlp := ProjectToPath(p)
	d.Projects[urlp] = &GitLab{Project: p, Commit: ""}
}

// Fetch will return the project belonging to path. Will return nil if not found.
func (d *Doc) Fetch(path string) *GitLab {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.Projects[path]
}

func (d *Doc) InsertFile(p *gitlab.Project, pathname string, data []byte) {
	urlp := ProjectToPath(p)
	gl := d.Fetch(urlp)
	if gl == nil {
		return
	}

	d.rw.Lock()
	defer d.rw.Unlock()
	full := path.Join(urlp, pathname)
	gl.Files[full] = data
}

func (d *Doc) FetchFile(p *gitlab.Project, pathname string) []byte {
	urlp := ProjectToPath(p)
	gl := d.Fetch(urlp)
	if gl == nil {
		return nil
	}

	d.rw.Lock()
	defer d.rw.Unlock()
	full := path.Join(urlp, pathname)
	return gl.Files[full]
}

// ProjectToPath converts a gitlab project to a path that can be used in Fetch.
func ProjectToPath(p *gitlab.Project) string {
	url, _ := url.Parse(p.WebURL)
	return url.Path
}
