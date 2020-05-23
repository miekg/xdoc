package main

import (
	"log"
	"net/url"
	"path"
	"path/filepath"
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
	log.Printf("InsertFile %q with %d bytes", full, len(buf))
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

// RemoveFirstPathElement removes the first element from the path p. This is need to remove the Docs dir from
// the files downloaded from GitLab.
func RemoveFirstPathElement(p string) string {
	el := filepath.SplitList(p) // this works on Linux, but will break on Windows
	if len(el) == 0 {
		return ""
	}
	return path.Join(el[1:]...)
}
